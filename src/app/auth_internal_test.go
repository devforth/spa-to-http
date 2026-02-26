package app

import (
	"go-http-server/param"
	"net/http"
	"os"
	"testing"
)

func TestAuthLoggerDisabled(t *testing.T) {
	params := param.Params{
		Logger: false,
	}
	app := NewApp(&params)
	if logger := app.authLogger(); logger != nil {
		t.Fatalf("expected nil logger when disabled")
	}
}

func TestAuthLoggerLevelsAndFormats(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		pretty   bool
	}{
		{"debug json", "debug", false},
		{"warn pretty", "warn", true},
		{"error json", "error", false},
		{"default info pretty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := os.Getenv("LOG_LEVEL")
			if tt.logLevel == "" {
				_ = os.Unsetenv("LOG_LEVEL")
			} else {
				_ = os.Setenv("LOG_LEVEL", tt.logLevel)
			}
			t.Cleanup(func() {
				if orig == "" {
					_ = os.Unsetenv("LOG_LEVEL")
				} else {
					_ = os.Setenv("LOG_LEVEL", orig)
				}
			})

			params := param.Params{
				Logger:    true,
				LogPretty: tt.pretty,
			}
			app := NewApp(&params)
			logger := app.authLogger()
			if logger == nil {
				t.Fatalf("expected logger, got nil")
			}
		})
	}
}

func TestListenWithAuthAndLoggerEnabled(t *testing.T) {
	params := param.Params{
		BasicAuthEnabled: true,
		BasicAuthUser:    "user",
		BasicAuthPass:    "pass",
		Logger:           true,
		LogPretty:        true,
		Address:          "127.0.0.1",
		Port:             8085,
	}

	called := false
	app := NewAppWithListenAndServe(&params, func(server *http.Server) error {
		if server == nil || server.Handler == nil {
			t.Fatalf("expected server and handler to be set")
		}
		called = true
		return nil
	})

	app.Listen()
	if !called {
		t.Fatalf("expected listenAndServe to be called")
	}
}
