package app

import (
	"go-http-server/param"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
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

func TestMapRequestPath(t *testing.T) {
	params := param.Params{
		Directory: ".",
		BasePath:  "/app",
	}
	app := NewApp(&params)

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "base path exact", in: "/app", want: "/"},
		{name: "base path with trailing slash", in: "/app/", want: "/"},
		{name: "asset under base path", in: "/app/assets/main.js", want: "/assets/main.js"},
		{name: "outside prefix fallback", in: "/assets/main.js", want: "/assets/main.js"},
		{name: "prefix-like but not matching", in: "/application/main.js", want: "/application/main.js"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.mapRequestPath(tt.in)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestMapRequestPathRootBasePathNoop(t *testing.T) {
	params := param.Params{
		Directory: ".",
		BasePath:  "/",
	}
	app := NewApp(&params)

	got := app.mapRequestPath("/app/assets/main.js")
	if got != "/app/assets/main.js" {
		t.Fatalf("expected path unchanged, got %q", got)
	}
}

func TestMapRequestPathEmptyBasePathNoop(t *testing.T) {
	params := param.Params{
		Directory: ".",
		BasePath:  "",
	}
	app := NewApp(&params)

	got := app.mapRequestPath("/app/assets/main.js")
	if got != "/app/assets/main.js" {
		t.Fatalf("expected path unchanged, got %q", got)
	}
}

func TestHandlerFuncNewBasePathCanonicalIndex(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.html")
	if err := os.WriteFile(indexPath, []byte("index"), 0600); err != nil {
		t.Fatalf("failed to write index: %v", err)
	}

	params := param.Params{
		Directory: dir,
		BasePath:  "/app",
		SpaMode:   true,
	}
	app := NewApp(&params)

	for _, reqPath := range []string{"/app", "/app/"} {
		req := httptest.NewRequest("GET", reqPath, nil)
		rec := httptest.NewRecorder()
		app.HandlerFuncNew(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200 for %s, got %d", reqPath, rec.Code)
		}
		if rec.Body.String() != "index" {
			t.Fatalf("expected index body for %s, got %q", reqPath, rec.Body.String())
		}
	}
}

func TestHandlerFuncNewBasePathAssetMappingAndIgnoreCacheControl(t *testing.T) {
	dir := t.TempDir()
	assetPath := filepath.Join(dir, "asset.bin")
	if err := os.WriteFile(assetPath, []byte("asset"), 0600); err != nil {
		t.Fatalf("failed to write asset: %v", err)
	}

	params := param.Params{
		Directory:               dir,
		BasePath:                "/app",
		SpaMode:                 true,
		IgnoreCacheControlPaths: []string{"/asset.bin"},
		CacheControlMaxAge:      3600,
	}
	app := NewApp(&params)

	req := httptest.NewRequest("GET", "/app/asset.bin", nil)
	rec := httptest.NewRecorder()
	app.HandlerFuncNew(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("expected no-store cache control, got %s", rec.Header().Get("Cache-Control"))
	}
	if rec.Body.String() != "asset" {
		t.Fatalf("expected asset body, got %q", rec.Body.String())
	}
}

func TestHandlerFuncNewBasePathTraversalRejected(t *testing.T) {
	dir := t.TempDir()
	params := param.Params{
		Directory: dir,
		BasePath:  "/app",
		SpaMode:   true,
	}
	app := NewApp(&params)

	req := httptest.NewRequest("GET", "/app/../secret.txt", nil)
	rec := httptest.NewRecorder()
	app.HandlerFuncNew(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandlerFuncNewBasePathOutsidePrefixFallsBackToRoot(t *testing.T) {
	dir := t.TempDir()
	assetPath := filepath.Join(dir, "asset.bin")
	if err := os.WriteFile(assetPath, []byte("asset"), 0600); err != nil {
		t.Fatalf("failed to write asset: %v", err)
	}

	params := param.Params{
		Directory:          dir,
		BasePath:           "/app",
		SpaMode:            true,
		CacheControlMaxAge: 3600,
	}
	app := NewApp(&params)

	req := httptest.NewRequest("GET", "/asset.bin", nil)
	rec := httptest.NewRecorder()
	app.HandlerFuncNew(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "asset" {
		t.Fatalf("expected asset body, got %q", rec.Body.String())
	}
}
