package app

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

const defaultBasicAuthRealm = "Restricted"

func (app *App) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.params.BasicAuthEnabled {
			next.ServeHTTP(w, r)
			return
		}

		realm := app.params.BasicAuthRealm
		if realm == "" {
			realm = defaultBasicAuthRealm
		}

		logger := app.authLogger()
		if logger != nil {
			logger.Debug("[AUTH] start",
				"path", r.URL.Path,
				"remoteAddr", r.RemoteAddr,
				"hasAuthHeader", r.Header.Get("Authorization") != "",
			)
		}

		user, pass, ok := r.BasicAuth()
		userMatch := constantTimeStringMatch(user, app.params.BasicAuthUser)
		passMatch := constantTimeStringMatch(pass, app.params.BasicAuthPass)
		authorized := ok && userMatch && passMatch

		if !authorized {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))
			w.WriteHeader(http.StatusUnauthorized)
			if logger != nil {
				logger.Info("[AUTH] failure",
					"path", r.URL.Path,
					"remoteAddr", r.RemoteAddr,
					"userProvided", user != "",
					"matchedUser", userMatch,
					"matchedPass", passMatch,
				)
			}
			return
		}

		if logger != nil {
			logger.Info("[AUTH] success",
				"path", r.URL.Path,
				"remoteAddr", r.RemoteAddr,
				"user", user,
			)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) authLogger() *slog.Logger {
	if !app.params.Logger {
		return nil
	}

	level := slog.LevelInfo
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: level}
	if app.params.LogPretty {
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func constantTimeStringMatch(a string, b string) bool {
	sumA := sha256.Sum256([]byte(a))
	sumB := sha256.Sum256([]byte(b))
	return subtle.ConstantTimeCompare(sumA[:], sumB[:]) == 1
}
