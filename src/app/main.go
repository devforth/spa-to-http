package app

import (
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	"go-http-server/param"
	"go-http-server/util"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type App struct {
	params *param.CliParams
	server *http.Server
}

func NewApp(params *param.CliParams) App {
	return App{params: params, server: nil}
}

func (app *App) CompressFiles() {
	if !app.params.Gzip && !app.params.Brotli {
		return
	}
	err := filepath.Walk(app.params.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Size() > app.params.Threshold {
			data, _ := ioutil.ReadFile(path)

			if app.params.Gzip {
				newName := path + ".gz"
				file, _ := os.Create(newName)

				writer := gzip.NewWriter(file)

				_, _ = writer.Write(data)
				_ = writer.Close()
			}

			if app.params.Brotli {
				newName := path + ".br"
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

func (app *App) HandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("New request:\n Url: %s\n", r.URL)

	requestedPath := path.Join(app.params.Directory, r.URL.Path)
	fileType := util.GetFileType(requestedPath)

	if fileType == util.FileTypeDirectory {
		if app.params.SpaMode {
			requestedPath = path.Join(app.params.Directory, "/index.html")
		} else if !app.params.DirectoryListing {
			requestedPath = path.Join(requestedPath, "/index.html")
		}
	}

	if fileType == util.FileTypeNotExists {
		if app.params.SpaMode {
			requestedPath = path.Join(app.params.Directory, "/index.html")
			fileType = util.GetFileType(requestedPath)

		}

		if fileType != util.FileTypeFile {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	if path.Ext(requestedPath) != ".html" {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", app.params.CacheControlMaxAge))
	}

	acceptEncodingHeader := r.Header.Get("Accept-Encoding")
	brotliApplicable := app.params.Brotli && (strings.Contains(acceptEncodingHeader, "*") || strings.Contains(acceptEncodingHeader, "br"))
	gzipApplicable := app.params.Gzip && (strings.Contains(acceptEncodingHeader, "*") || strings.Contains(acceptEncodingHeader, "gzip"))

	if brotliApplicable && util.GetFileType(requestedPath+".br") == util.FileTypeFile {
		requestedPath = requestedPath + ".br"
		w.Header().Set("Content-Encoding", "br")
	} else if gzipApplicable && util.GetFileType(requestedPath+".gz") == util.FileTypeFile {
		requestedPath = requestedPath + ".gz"
		w.Header().Set("Content-Encoding", "gzip")
	}

	http.ServeFile(w, r, requestedPath)
}

func (app *App) Listen() {
	app.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", app.params.Host, app.params.Port),
		Handler: http.HandlerFunc(app.HandleFunc),
	}

	fmt.Printf("Server listening on http://%s\n", app.server.Addr)
	err := app.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
