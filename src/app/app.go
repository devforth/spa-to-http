package app

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	lru "github.com/hashicorp/golang-lru"
	"go-http-server/param"
	"go-http-server/util"
	"golang.org/x/exp/slices"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type App struct {
	params *param.Params
	server *http.Server
	cache  *lru.TwoQueueCache
}

type ResponseItem struct {
	Name        string
	Path        string
	ModTime     time.Time
	Content     []byte
	ContentType string
}

type Compression int

const (
	None Compression = iota
	Gzip
	Brotli
)

func NewApp(params *param.Params) App {
	var cache *lru.TwoQueueCache = nil
	var err error

	if params.CacheEnabled {
		cache, err = lru.New2Q(params.CacheBuffer)
		if err != nil {
			panic(err)
		}
	}

	return App{params: params, server: nil, cache: cache}
}

func (app *App) shouldSkipCompression(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, blocked := range app.params.NoCompress {
		if strings.ToLower(blocked) == ext {
			return true
		}
	}
	return false
}

func (app *App) CompressFiles() {
	if !app.params.Gzip && !app.params.Brotli {
		return
	}
	err := filepath.Walk(app.params.Directory, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := path.Ext(filePath)
		if ext == ".br" || ext == ".gz" {
			return nil
		}

		if app.shouldSkipCompression(filePath) {
			return nil
		}

		if info.Size() > app.params.Threshold {
			data, _ := os.ReadFile(filePath)

			if app.params.Gzip {
				newName := filePath + ".gz"
				file, _ := os.Create(newName)

				writer := gzip.NewWriter(file)

				_, _ = writer.Write(data)
				_ = writer.Close()
			}

			if app.params.Brotli {
				newName := filePath + ".br"
				file, _ := os.Create(newName)

				writer := brotli.NewWriter(file)

				_, _ = writer.Write(data)
				_ = writer.Close()
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (app *App) GetOrCreateResponseItem(requestedPath string, compression Compression, actualContentType *string) (*ResponseItem, int) {
	rootIndexPath := path.Join(app.params.Directory, "index.html")

	switch compression {
	case Gzip:
		requestedPath = requestedPath + ".gz"
	case Brotli:
		requestedPath = requestedPath + ".br"
	}

	if app.cache != nil {
		cacheValue, ok := app.cache.Get(requestedPath)
		if ok {
			responseItem, ok := cacheValue.(ResponseItem)
			if ok {
				return &responseItem, 0
			} else {
				localRedirectItem := cacheValue.(string)
				return app.GetOrCreateResponseItem(localRedirectItem, compression, actualContentType)
			}
		}
	}

	dirPath, fileName := filepath.Split(requestedPath)
	dir := http.Dir(dirPath)

	file, err := dir.Open(fileName)
	if err != nil {
		if app.params.SpaMode && compression == None && requestedPath != rootIndexPath {
			newPath := path.Join(app.params.Directory, "index.html")
			if app.cache != nil {
				app.cache.Add(requestedPath, newPath)
			}
			return app.GetOrCreateResponseItem(newPath, compression, actualContentType)
		}
		return nil, http.StatusNotFound
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		if app.params.SpaMode && compression == None && requestedPath != rootIndexPath {
			newPath := path.Join(app.params.Directory, "index.html")
			if app.cache != nil {
				app.cache.Add(requestedPath, newPath)
			}
			return app.GetOrCreateResponseItem(newPath, compression, actualContentType)
		}
		return nil, http.StatusNotFound
	}

	if stat.IsDir() && requestedPath != rootIndexPath {
		if app.params.SpaMode && compression == None {
			newPath := path.Join(app.params.Directory, "index.html")
			if app.cache != nil {
				app.cache.Add(requestedPath, newPath)
			}
			return app.GetOrCreateResponseItem(newPath, compression, actualContentType)
		} else if !app.params.SpaMode && compression == None {
			newPath := path.Join(requestedPath, "index.html")
			if app.cache != nil {
				app.cache.Add(requestedPath, newPath)
			}
			return app.GetOrCreateResponseItem(newPath, compression, actualContentType)
		}

		return nil, http.StatusNotFound
	}

	content := make([]byte, stat.Size())
	_, err = io.ReadFull(file, content)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	name := stat.Name()
	var contentType string
	if compression == None {
		contentType = mime.TypeByExtension(filepath.Ext(name))
	} else {
		contentType = *actualContentType
	}

	responseItem := ResponseItem{
		Path:        requestedPath,
		Name:        name,
		ModTime:     stat.ModTime(),
		Content:     content,
		ContentType: contentType,
	}

	if app.cache != nil {
		app.cache.Add(requestedPath, responseItem)
	}

	return &responseItem, 0
}

func (app *App) GetFilePath(urlPath string) (string, bool) {
	requestedPath := path.Join(app.params.Directory, urlPath)

	if _, err := os.Stat(requestedPath); !os.IsNotExist(err) {
		requestedPath, err = filepath.EvalSymlinks(requestedPath)
	}

	if !strings.HasPrefix(requestedPath, app.params.Directory) {
		return "", false
	}

	return requestedPath, true
}

func (app *App) HandlerFuncNew(w http.ResponseWriter, r *http.Request) {
	requestedPath, valid := app.GetFilePath(r.URL.Path)

	if !valid {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responseItem, errorCode := app.GetOrCreateResponseItem(requestedPath, None, nil)
	if errorCode != 0 {
		w.WriteHeader(errorCode)
		return
	}

	if r.Header.Get("Range") != "" || app.shouldSkipCompression(requestedPath) {
		if responseItem.ContentType != "" {
			w.Header().Set("Content-Type", responseItem.ContentType)
		}
		http.ServeContent(w, r, responseItem.Name, responseItem.ModTime, bytes.NewReader(responseItem.Content))
		return
	}

	if slices.Contains(app.params.IgnoreCacheControlPaths, r.URL.Path) || path.Ext(responseItem.Name) == ".html" {
		w.Header().Set("Cache-Control", "no-store")
	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", app.params.CacheControlMaxAge))
	}

	var brotliApplicable bool
	var gzipApplicable bool
	var overThreshold = int64(len(responseItem.Content)) > app.params.Threshold

	acceptEncodingHeader := r.Header.Get("Accept-Encoding")
	if app.params.Brotli || app.params.Gzip {
		if strings.Contains(acceptEncodingHeader, "*") {
			brotliApplicable = app.params.Brotli
			gzipApplicable = app.params.Gzip
		} else {
			brotliApplicable = app.params.Brotli && strings.Contains(acceptEncodingHeader, "br")
			gzipApplicable = app.params.Gzip && strings.Contains(acceptEncodingHeader, "gzip")

		}
	}

	if brotliApplicable && overThreshold {
		brotliResponseItem, _ := app.GetOrCreateResponseItem(responseItem.Path, Brotli, &responseItem.ContentType)

		if brotliResponseItem != nil {
			responseItem = brotliResponseItem
			w.Header().Set("Content-Encoding", "br")
		}
	} else if gzipApplicable && overThreshold {
		gzipResponseItem, _ := app.GetOrCreateResponseItem(responseItem.Path, Gzip, &responseItem.ContentType)

		if gzipResponseItem != nil {
			responseItem = gzipResponseItem
			w.Header().Set("Content-Encoding", "gzip")
		}
	}

	if responseItem.ContentType != "" {
		w.Header().Set("Content-Type", responseItem.ContentType)
	}

	http.ServeContent(w, r, responseItem.Name, responseItem.ModTime, bytes.NewReader(responseItem.Content))
}

func (app *App) Listen() {
	var handlerFunc http.Handler = http.HandlerFunc(app.HandlerFuncNew)
	if app.params.Logger {
		handlerFunc = util.LogRequestHandler(handlerFunc, &util.LogRequestHandlerOptions{
			Pretty: app.params.LogPretty,
		})
	}

	app.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", app.params.Address, app.params.Port),
		Handler: handlerFunc,
	}

	fmt.Printf("Server listening on http://%s\n", app.server.Addr)
	err := app.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
