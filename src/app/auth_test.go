package app

import (
	"go-http-server/param"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthMiddlewareDisabled(t *testing.T) {
	params := param.Params{
		BasicAuthEnabled: false,
	}
	app := NewApp(&params)

	handler := app.BasicAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("expected body ok, got %s", rec.Body.String())
	}
}

func TestBasicAuthMiddlewareMissingHeader(t *testing.T) {
	params := param.Params{
		BasicAuthEnabled: true,
		BasicAuthUser:    "user",
		BasicAuthPass:    "pass",
	}
	app := NewApp(&params)

	handler := app.BasicAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if got := rec.Header().Get("WWW-Authenticate"); got != "Basic realm=\"Restricted\"" {
		t.Fatalf("expected realm header, got %s", got)
	}
}

func TestBasicAuthMiddlewareWrongCredentials(t *testing.T) {
	params := param.Params{
		BasicAuthEnabled: true,
		BasicAuthUser:    "user",
		BasicAuthPass:    "pass",
		BasicAuthRealm:   "CustomRealm",
	}
	app := NewApp(&params)

	handler := app.BasicAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("user", "wrong")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if got := rec.Header().Get("WWW-Authenticate"); got != "Basic realm=\"CustomRealm\"" {
		t.Fatalf("expected realm header, got %s", got)
	}
}

func TestBasicAuthMiddlewareSuccess(t *testing.T) {
	params := param.Params{
		BasicAuthEnabled: true,
		BasicAuthUser:    "user",
		BasicAuthPass:    "pass",
	}
	app := NewApp(&params)

	handler := app.BasicAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("user", "pass")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("expected body ok, got %s", rec.Body.String())
	}
}
