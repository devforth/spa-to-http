package app_test

import (
	"bytes"
	"compress/gzip"
	"errors"
	"go-http-server/app"
	"go-http-server/param"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/andybalholm/brotli"
)

func TestNewApp(t *testing.T) {
	params := param.Params{
		Address:                 "0.0.0.0",
		Port:                    8080,
		Gzip:                    true,
		Brotli:                  true,
		Threshold:               1024,
		Directory:               "../../test/frontend/dist",
		CacheControlMaxAge:      604800,
		SpaMode:                 true,
		IgnoreCacheControlPaths: nil,
		CacheEnabled:            true,
		CacheBuffer:             50 * 1024,
	}

	app1 := app.NewApp(&params)
	if reflect.TypeOf(app1) == reflect.TypeOf(nil) {
		t.Errorf("app1 is nil")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	params.CacheBuffer = 0
	app2 := app.NewApp(&params)
	if reflect.TypeOf(app2) != nil {
		t.Errorf("app2 is not nil")
	}
}

func TestCompressFiles(t *testing.T) {
	params := param.Params{
		Address:                 "0.0.0.0",
		Port:                    8080,
		Gzip:                    false,
		Brotli:                  false,
		Threshold:               1024,
		Directory:               "../../test/frontend/dist/vite.svg",
		CacheControlMaxAge:      604800,
		SpaMode:                 true,
		IgnoreCacheControlPaths: nil,
		CacheEnabled:            true,
		CacheBuffer:             50 * 1024,
	}
	app1 := app.NewApp(&params)
	app1.CompressFiles()

	params.Gzip = true
	app2 := app.NewApp(&params)
	app2.CompressFiles()

	info, _ := os.Stat(params.Directory)
	info_gz, _ := os.Stat(params.Directory + ".gz")
	if info.Size() < info_gz.Size() {
		t.Error("Original file is smaller than compressed .gz!")
	}

	params.Gzip = false
	params.Brotli = true
	app3 := app.NewApp(&params)
	app3.CompressFiles()

	info_br, _ := os.Stat(params.Directory + ".br")
	if info.Size() < info_br.Size() {
		t.Error("Original file is smaller than compressed .br!")
	}

	params.Directory = "fasdfa.go"
	app4 := app.NewApp(&params)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}

		params.Directory = "../app"
		app5 := app.NewApp(&params)
		app5.CompressFiles()
		os.Remove("app_test.go.br")
		os.Remove("app.go.br")

		params.Directory = "../../test/frontend/dist/vite.svg.br"
		app6 := app.NewApp(&params)
		app6.CompressFiles()
	}()
	app4.CompressFiles()
}

func TestGetOrCreateResponseItem(t *testing.T) {
	params := param.Params{
		Address:                 "0.0.0.0",
		Port:                    8000,
		Gzip:                    true,
		Brotli:                  false,
		Threshold:               1024,
		Directory:               "../../test/frontend/dist/",
		CacheControlMaxAge:      604800,
		SpaMode:                 true,
		IgnoreCacheControlPaths: nil,
		CacheEnabled:            true,
		CacheBuffer:             50 * 1024,
	}
	apl := app.NewApp(&params)
	text := "text/html; charset=utf-8"

	resp1, _ := apl.GetOrCreateResponseItem("/", 0, &text)
	if resp1.Name != "index.html" {
		t.Errorf("Expected index.html to return, got %s", resp1.Name)
	}

	params.SpaMode = false
	apl1 := app.NewApp(&params)
	resp11, _ := apl1.GetOrCreateResponseItem("../../test/frontend/dist/", 0, &text)
	if resp11.Name != "index.html" {
		t.Errorf("Expected index.html to return, got %s", resp11.Name)
	}

	resp2, _ := apl.GetOrCreateResponseItem("../../test/frontend/dist/vite.svg", 0, &text)
	if resp2.Name != "vite.svg" {
		t.Errorf("Expected vite.svg to return, got %s", resp2.Name)
	}

	resp3, _ := apl.GetOrCreateResponseItem("../../src/app/app_test.go", 0, &text)
	if resp3.Name != "app_test.go" {
		t.Errorf("Expected app_test.go to return, got %s", resp3.Name)
	}

	resp4, _ := apl.GetOrCreateResponseItem("../../test/frontend/dist/vite.svg", 1, &text)
	if resp4.Name != "vite.svg.gz" {
		t.Errorf("Expected vite.svg.gz to return, got %s", resp4.Name)
	}

	resp5, _ := apl.GetOrCreateResponseItem("../../test/frontend/dist/vite.svg", 2, &text)
	if resp5.Name != "vite.svg.br" {
		t.Errorf("Expected vite.svg.br to return, got %s", resp5.Name)
	}

	resp6, _ := apl.GetOrCreateResponseItem("/fdsfds.go", 0, &text)
	if resp6 != nil {
		t.Errorf("Expected nil to return, got %s", resp5)
	}
}

func TestHandlerFuncNew(t *testing.T) {
	params := param.Params{
		Address:                 "0.0.0.0",
		Port:                    8080,
		Gzip:                    false,
		Brotli:                  true,
		Threshold:               1024,
		Directory:               "../../test/frontend/dist",
		CacheControlMaxAge:      99999999999,
		SpaMode:                 true,
		IgnoreCacheControlPaths: []string{"../../test/frontend/dist/example.html"},
		CacheEnabled:            true,
		CacheBuffer:             50 * 1024,
	}
	app1 := app.NewApp(&params)
	app1.CompressFiles()
	index_content, _ := ioutil.ReadFile("../../test/frontend/dist/index.html")
	example_content, _ := ioutil.ReadFile("../../test/frontend/dist/example.html")
	vite_content, _ := ioutil.ReadFile("../../test/frontend/dist/vite.svg")

	req1, _ := http.NewRequest("GET", "/", nil)
	recorder1 := httptest.NewRecorder()
	app1.HandlerFuncNew(recorder1, req1)
	if recorder1.HeaderMap["Cache-Control"][0] != "no-store" {
		t.Errorf("Expected no-store to return, got %s", recorder1.HeaderMap["Cache-Control"])
	}
	if recorder1.Body.String() != string(index_content) {
		t.Errorf("Expected index.html body to return, got %s", recorder1.Body)
	}

	req2, _ := http.NewRequest("GET", "/example.html", nil)
	recorder2 := httptest.NewRecorder()
	app1.HandlerFuncNew(recorder2, req2)
	if recorder2.HeaderMap["Cache-Control"][0] != "no-store" {
		t.Errorf("Expected no-store to return, got %s", recorder2.HeaderMap["Cache-Control"])
	}
	if recorder2.Body.String() != string(example_content) {
		t.Errorf("Expected exapmle.html body to return, got %s", recorder2.Body)
	}

	req3, _ := http.NewRequest("GET", "/random", nil)
	recorder3 := httptest.NewRecorder()
	app1.HandlerFuncNew(recorder3, req3)
	if recorder3.HeaderMap["Cache-Control"][0] != "no-store" {
		t.Errorf("Expected no-store to return, got %s", recorder3.HeaderMap["Cache-Control"])
	}
	if recorder3.Body.String() != string(index_content) {
		t.Errorf("Expected index.html body to return, got %s", recorder3.Body)
	}

	req4, _ := http.NewRequest("GET", "/vite.svg", nil)
	req4.Header.Set("Accept-Encoding", "*")
	recorder4 := httptest.NewRecorder()
	app1.HandlerFuncNew(recorder4, req4)
	resp_body_decoded3, _ := ioutil.ReadAll(brotli.NewReader(recorder4.Body))
	if recorder4.HeaderMap["Cache-Control"][0] != "max-age=99999999999" {
		t.Errorf("Expected Cache-Control = no-store to return, got %s", recorder4.HeaderMap["Cache-Control"])
	}
	if recorder4.HeaderMap["Content-Encoding"][0] != "br" {
		t.Errorf("Expected Content-Encoding = br to return, got %s", recorder4.HeaderMap["Content-Encoding"])
	}
	if string(resp_body_decoded3) != string(vite_content) {
		t.Errorf("Expected vite.svg to return, got %s", recorder4.Body)
	}

	req5, _ := http.NewRequest("GET", "/vite.svg", nil)
	req5.Header.Set("Accept-Encoding", "br, gzip")
	recorder5 := httptest.NewRecorder()
	app1.HandlerFuncNew(recorder5, req5)
	resp_body_decoded5, _ := ioutil.ReadAll(brotli.NewReader(recorder5.Body))
	if recorder5.HeaderMap["Cache-Control"][0] != "max-age=99999999999" {
		t.Errorf("Expected no-store to return, got %s", recorder5.HeaderMap["Cache-Control"])
	}
	if string(resp_body_decoded5) != string(vite_content) {
		t.Errorf("Expected vite.svg body to return, got %s", recorder5.Body)
	}

	req6, _ := http.NewRequest("GET", "/vite.svg", nil)
	req6.Header.Set("Accept-Encoding", "gzip")
	recorder6 := httptest.NewRecorder()
	params.Brotli = false
	params.Gzip = true
	app2 := app.NewApp(&params)
	app2.HandlerFuncNew(recorder6, req6)
	reader6 := bytes.NewReader(recorder6.Body.Bytes())
	gzreader6, _ := gzip.NewReader(reader6)
	resp_body_decoded6, _ := ioutil.ReadAll(gzreader6)
	if recorder6.HeaderMap["Cache-Control"][0] != "max-age=99999999999" {
		t.Errorf("Expected no-store to return, got %s", recorder6.HeaderMap["Cache-Control"])
	}
	if recorder6.HeaderMap["Content-Encoding"][0] != "gzip" {
		t.Errorf("Expected Content-Encoding = gzip to return, got %s", recorder6.HeaderMap["Cache-Control"])
	}
	if string(resp_body_decoded6) != string(vite_content) {
		t.Errorf("Expected vite.svg to return, got %s", recorder6.Body)
	}

	req7, _ := http.NewRequest("GET", "/vite.svg", nil)
	req7.Header.Set("Accept-Encoding", "gzip")
	recorder7 := httptest.NewRecorder()
	params.Brotli = true
	params.Gzip = false
	app3 := app.NewApp(&params)
	app3.HandlerFuncNew(recorder7, req7)
	_, err := gzip.NewReader(bytes.NewReader(recorder7.Body.Bytes()))
	if err.Error() != errors.New("gzip: invalid header").Error() {
		t.Errorf("Expected: \"gzip: invalid header\", got %s", err)
	}
	if recorder7.HeaderMap["Cache-Control"][0] != "max-age=99999999999" {
		t.Errorf("Expected Cache-Control = max-age=99999999999 to return, got %s", recorder7.HeaderMap["Cache-Control"])
	}
	if string(resp_body_decoded6) != string(vite_content) {
		t.Errorf("Expected vite.svg body to return, got %s", recorder7.Body)
	}
	os.Remove("../../test/frontend/dist/vite.svg.gz")
	os.Remove("../../test/frontend/dist/vite.svg.br")
}

func TestListen(t *testing.T) {
	params := param.Params{
		Address:                 "0.0.0.0",
		Port:                    8000,
		Gzip:                    false,
		Brotli:                  false,
		Threshold:               1024,
		Directory:               "../../test/frontend/dist",
		CacheControlMaxAge:      604800,
		SpaMode:                 true,
		IgnoreCacheControlPaths: nil,
		CacheEnabled:            true,
		CacheBuffer:             50 * 1024,
	}
	a := app.NewApp(&params)

	var l *http.Server
	monkey.PatchInstanceMethod(reflect.TypeOf(l), "ListenAndServe", func(*http.Server) error {
		return http.ErrServerClosed
	})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	a.Listen()
}

func TestGetFilePath(t *testing.T) {
	params := param.Params{
		Address:                 "0.0.0.0",
		Port:                    8080,
		Gzip:                    false,
		Brotli:                  true,
		Threshold:               1024,
		Directory:               "../../test/frontend/dist",
		CacheControlMaxAge:      99999999999,
		SpaMode:                 true,
		IgnoreCacheControlPaths: []string{"../../test/frontend/dist/example.html"},
		CacheEnabled:            true,
		CacheBuffer:             50 * 1024,
	}
	app := app.NewApp(&params)

	_, valid := app.GetFilePath("../test/index.html")
	if valid {
		t.Errorf("Expected false, got %t", valid)
	}

	_, valid = app.GetFilePath("../../test/index.html")
	if valid {
		t.Errorf("Expected false, got %t", valid)
	}

	_, valid = app.GetFilePath("test/../index.html")
	if !valid {
		t.Errorf("Expected false, got %t", valid)
	}

	_, valid = app.GetFilePath("test/../test/../index.html")
	if !valid {
		t.Errorf("Expected false, got %t", valid)
	}
}
