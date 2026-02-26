package app

import (
	"go-http-server/param"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewAppWithListenAndServeDefault(t *testing.T) {
	params := param.Params{
		Directory: ".",
	}
	app := NewAppWithListenAndServe(&params, nil)
	if app.listenAndServe == nil {
		t.Fatal("expected listenAndServe to be set")
	}
}

func TestGetOrCreateResponseItemCacheHit(t *testing.T) {
	params := param.Params{
		Directory:    ".",
		CacheEnabled: true,
		CacheBuffer:  10,
	}
	app := NewApp(&params)

	cached := ResponseItem{
		Name:        "cached.txt",
		Path:        "/cached.txt",
		ModTime:     time.Now(),
		Content:     []byte("cached"),
		ContentType: "text/plain",
	}

	app.cache.Add("/cached.txt", cached)

	got, code := app.GetOrCreateResponseItem("/cached.txt", None, nil)
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if got == nil || got.Name != "cached.txt" {
		t.Fatalf("expected cached response item, got %#v", got)
	}
}

func TestGetOrCreateResponseItemCacheRedirect(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.html")
	if err := os.WriteFile(indexPath, []byte("index"), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	params := param.Params{
		Directory:    dir,
		CacheEnabled: true,
		CacheBuffer:  10,
		SpaMode:      true,
	}
	app := NewApp(&params)

	app.cache.Add("/alias", indexPath)

	got, code := app.GetOrCreateResponseItem("/alias", None, nil)
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if got == nil || got.Name != "index.html" {
		t.Fatalf("expected index.html, got %#v", got)
	}
}

type errStatFile struct{}

func (f *errStatFile) Close() error                       { return nil }
func (f *errStatFile) Read(_ []byte) (int, error)         { return 0, io.EOF }
func (f *errStatFile) Seek(_ int64, _ int) (int64, error) { return 0, nil }
func (f *errStatFile) Readdir(_ int) ([]fs.FileInfo, error) {
	return nil, nil
}
func (f *errStatFile) Stat() (fs.FileInfo, error) {
	return nil, fs.ErrInvalid
}

type fixedFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

func (fi fixedFileInfo) Name() string       { return fi.name }
func (fi fixedFileInfo) Size() int64        { return fi.size }
func (fi fixedFileInfo) Mode() fs.FileMode  { return fi.mode }
func (fi fixedFileInfo) ModTime() time.Time { return fi.modTime }
func (fi fixedFileInfo) IsDir() bool        { return fi.isDir }
func (fi fixedFileInfo) Sys() interface{}   { return nil }

type shortReadFile struct {
	info fs.FileInfo
}

func (f *shortReadFile) Close() error                       { return nil }
func (f *shortReadFile) Read(_ []byte) (int, error)         { return 0, io.EOF }
func (f *shortReadFile) Seek(_ int64, _ int) (int64, error) { return 0, nil }
func (f *shortReadFile) Readdir(_ int) ([]fs.FileInfo, error) {
	return nil, nil
}
func (f *shortReadFile) Stat() (fs.FileInfo, error) { return f.info, nil }

func TestGetOrCreateResponseItemOpenErrorSpaRedirect(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.html")
	if err := os.WriteFile(indexPath, []byte("index"), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	params := param.Params{
		Directory:    dir,
		SpaMode:      true,
		CacheEnabled: true,
		CacheBuffer:  10,
	}
	app := NewApp(&params)

	opener := func(dirPath, fileName string) (http.File, error) {
		if fileName == "index.html" {
			return os.Open(filepath.Join(dirPath, fileName))
		}
		return nil, fs.ErrNotExist
	}

	got, code := app.getOrCreateResponseItemWithOpener(filepath.Join(dir, "missing.txt"), None, nil, opener)
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if got == nil || got.Name != "index.html" {
		t.Fatalf("expected index.html, got %#v", got)
	}
}

func TestGetOrCreateResponseItemStatError(t *testing.T) {
	params := param.Params{
		Directory: ".",
		SpaMode:   false,
	}
	app := NewApp(&params)

	opener := func(_, _ string) (http.File, error) {
		return &errStatFile{}, nil
	}

	got, code := app.getOrCreateResponseItemWithOpener("file.txt", None, nil, opener)
	if got != nil {
		t.Fatalf("expected nil response, got %#v", got)
	}
	if code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, code)
	}
}

func TestGetOrCreateResponseItemStatErrorSpaRedirect(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.html")
	if err := os.WriteFile(indexPath, []byte("index"), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	params := param.Params{
		Directory:    dir,
		SpaMode:      true,
		CacheEnabled: true,
		CacheBuffer:  10,
	}
	app := NewApp(&params)

	opener := func(dirPath, fileName string) (http.File, error) {
		if fileName == "index.html" {
			return os.Open(filepath.Join(dirPath, fileName))
		}
		return &errStatFile{}, nil
	}

	got, code := app.getOrCreateResponseItemWithOpener(filepath.Join(dir, "missing.txt"), None, nil, opener)
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if got == nil || got.Name != "index.html" {
		t.Fatalf("expected index.html, got %#v", got)
	}
}

func TestGetOrCreateResponseItemReadError(t *testing.T) {
	params := param.Params{
		Directory: ".",
		SpaMode:   false,
	}
	app := NewApp(&params)

	info := fixedFileInfo{
		name:    "file.txt",
		size:    10,
		mode:    0600,
		modTime: time.Now(),
		isDir:   false,
	}
	opener := func(_, _ string) (http.File, error) {
		return &shortReadFile{info: info}, nil
	}

	got, code := app.getOrCreateResponseItemWithOpener("file.txt", None, nil, opener)
	if got != nil {
		t.Fatalf("expected nil response, got %#v", got)
	}
	if code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, code)
	}
}

func TestGetOrCreateResponseItemDirWithCompressionNotFound(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "dir.gz"), 0700); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	params := param.Params{
		Directory: dir,
		SpaMode:   true,
	}
	app := NewApp(&params)

	got, code := app.getOrCreateResponseItemWithOpener(filepath.Join(dir, "dir"), Gzip, nil, nil)
	if got != nil {
		t.Fatalf("expected nil response, got %#v", got)
	}
	if code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, code)
	}
}
