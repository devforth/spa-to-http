package util

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/felixge/httpsnoop"
)

func TestLogHTTPReqInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	ri := &HTTPReqInfo{
		method:    "GET",
		path:      "/test/path",
		code:      200,
		size:      1234,
		duration:  150 * time.Millisecond,
		ipAddress: net.ParseIP("127.0.0.1"),
		userAgent: "Go-http-client/1.1",
		referer:   "http://example.com",
	}

	logHTTPReqInfo(logger, ri)

	logged := buf.String()

	// Parse the JSON log entry
	var logData map[string]interface{}
	if err := json.Unmarshal([]byte(logged), &logData); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v\nLog output: %s", err, logged)
	}

	tests := []struct {
		field string
		want  interface{}
	}{
		{"method", "GET"},
		{"path", "/test/path"},
		{"code", float64(200)}, // JSON numbers are float64
		{"size", float64(1234)},
		{"ipAddress", "127.0.0.1"},
		{"userAgent", "Go-http-client/1.1"},
		{"referer", "http://example.com"},
	}

	for _, tt := range tests {
		if value, ok := logData[tt.field]; !ok {
			t.Errorf("Expected log to contain field %q, got: %s", tt.field, logged)
		} else if value != tt.want {
			t.Errorf("Expected field %q to be %v, got %v", tt.field, tt.want, value)
		}
	}

	// Check that duration field exists and is a number (slog duration format in nanoseconds)
	if duration, ok := logData["duration"]; !ok {
		t.Errorf("Expected log to contain 'duration' field, got: %s", logged)
	} else if _, ok := duration.(float64); !ok {
		t.Errorf("Expected duration to be a number, got %T: %v", duration, duration)
	}

	// Check that msg field exists
	if msg, ok := logData["msg"]; !ok {
		t.Errorf("Expected log to contain 'msg' field, got: %s", logged)
	} else if msg != "HTTP Request" {
		t.Errorf("Expected msg to be 'HTTP Request', got %v", msg)
	}
}

func TestLogRequestHandler(t *testing.T) {
	tests := []struct {
		name        string
		pretty      bool
		method      string
		path        string
		userAgent   string
		referer     string
		remoteAddr  string
		wantMethod  string
		wantPath    string
		wantAgent   string
		wantReferer string
	}{
		{
			name:        "GET request with all headers",
			pretty:      false,
			method:      "GET",
			path:        "/api/test",
			userAgent:   "Mozilla/5.0",
			referer:     "https://example.com",
			remoteAddr:  "192.168.1.1:12345",
			wantMethod:  "GET",
			wantPath:    "/api/test",
			wantAgent:   "Mozilla/5.0",
			wantReferer: "https://example.com",
		},
		{
			name:        "POST request without headers",
			pretty:      false,
			method:      "POST",
			path:        "/api/submit",
			userAgent:   "",
			referer:     "",
			remoteAddr:  "127.0.0.1:8080",
			wantMethod:  "POST",
			wantPath:    "/api/submit",
			wantAgent:   "",
			wantReferer: "",
		},
		{
			name:        "PUT request with query parameters",
			pretty:      false,
			method:      "PUT",
			path:        "/api/update?id=123&param=value",
			userAgent:   "Go-http-client/1.1",
			referer:     "http://localhost:3000",
			remoteAddr:  "10.0.0.1:54321",
			wantMethod:  "PUT",
			wantPath:    "/api/update?id=123&param=value",
			wantAgent:   "Go-http-client/1.1",
			wantReferer: "http://localhost:3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output to a buffer instead of stdout
			var buf bytes.Buffer

			// Create the logging handler with a custom logger that writes to our buffer
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test response"))
			})

			// Create a custom version of LogRequestHandler that uses our buffer for logging
			var logger *slog.Logger
			if tt.pretty {
				logger = slog.New(slog.NewTextHandler(&buf, nil))
			} else {
				logger = slog.New(slog.NewJSONHandler(&buf, nil))
			}

			// Create a modified version of LogRequestHandler that uses our logger
			fn := func(w http.ResponseWriter, r *http.Request) {
				// runs handler and captures information about HTTP request
				mtr := httpsnoop.CaptureMetrics(dummyHandler, w, r)

				logHTTPReqInfo(logger, &HTTPReqInfo{
					method:    r.Method,
					path:      r.URL.String(),
					code:      mtr.Code,
					size:      mtr.Written,
					duration:  mtr.Duration,
					ipAddress: requestGetRemoteAddress(r),
					userAgent: r.Header.Get("User-Agent"),
					referer:   r.Header.Get("Referer"),
				})
			}
			handler := http.HandlerFunc(fn)

			// Create test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.userAgent != "" {
				req.Header.Set("User-Agent", tt.userAgent)
			}
			if tt.referer != "" {
				req.Header.Set("Referer", tt.referer)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute the handler
			handler.ServeHTTP(w, req)

			// Verify the response
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			if w.Body.String() != "test response" {
				t.Errorf("Expected body 'test response', got '%s'", w.Body.String())
			}

			// Verify the log output
			logged := buf.String()

			// Basic check that we got some log output
			if len(logged) == 0 {
				t.Errorf("Expected log output, got empty string")
				return
			}

			if !tt.pretty {
				// For JSON format, parse and validate the structure
				var logData map[string]interface{}
				if err := json.Unmarshal([]byte(logged), &logData); err != nil {
					t.Errorf("Failed to parse log output as JSON: %v\nLog output: %s", err, logged)
					return
				}

				// Check required fields
				if method, ok := logData["method"]; !ok || method != tt.wantMethod {
					t.Errorf("Expected method %q, got %v", tt.wantMethod, method)
				}
				if path, ok := logData["path"]; !ok || path != tt.wantPath {
					t.Errorf("Expected path %q, got %v", tt.wantPath, path)
				}
				if code, ok := logData["code"]; !ok || code != float64(200) {
					t.Errorf("Expected code 200, got %v", code)
				}
				if size, ok := logData["size"]; !ok || size != float64(13) { // "test response" is 13 bytes
					t.Errorf("Expected size 13, got %v", size)
				}

				// Check optional header fields only if they should be present
				if tt.wantAgent != "" {
					if userAgent, ok := logData["userAgent"]; !ok || userAgent != tt.wantAgent {
						t.Errorf("Expected userAgent %q, got %v", tt.wantAgent, userAgent)
					}
				}
				if tt.wantReferer != "" {
					if referer, ok := logData["referer"]; !ok || referer != tt.wantReferer {
						t.Errorf("Expected referer %q, got %v", tt.wantReferer, referer)
					}
				}

				// Verify duration and ipAddress fields exist
				if _, ok := logData["duration"]; !ok {
					t.Errorf("Expected log to contain 'duration' field, got: %s", logged)
				}
				if _, ok := logData["ipAddress"]; !ok {
					t.Errorf("Expected log to contain 'ipAddress' field, got: %s", logged)
				}
			} else {
				// For pretty format, just check that key information is present
				if !strings.Contains(logged, tt.wantMethod) {
					t.Errorf("Expected log to contain method %q, got: %s", tt.wantMethod, logged)
				}
				if !strings.Contains(logged, tt.wantPath) {
					t.Errorf("Expected log to contain path %q, got: %s", tt.wantPath, logged)
				}
			}
		})
	}
}

func TestLogRequestHandlerWithDifferentStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{"Not Found", http.StatusNotFound, "not found"},
		{"Internal Server Error", http.StatusInternalServerError, "error occurred"},
		{"Created", http.StatusCreated, "resource created"},
		{"No Content", http.StatusNoContent, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			// Create a dummy handler that returns the specified status code
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			})

			// Create a custom handler that uses our logger
			fn := func(w http.ResponseWriter, r *http.Request) {
				mtr := httpsnoop.CaptureMetrics(dummyHandler, w, r)

				logHTTPReqInfo(logger, &HTTPReqInfo{
					method:    r.Method,
					path:      r.URL.String(),
					code:      mtr.Code,
					size:      mtr.Written,
					duration:  mtr.Duration,
					ipAddress: requestGetRemoteAddress(r),
					userAgent: r.Header.Get("User-Agent"),
					referer:   r.Header.Get("Referer"),
				})
			}
			handler := http.HandlerFunc(fn)

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "127.0.0.1:12345"

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute the handler
			handler.ServeHTTP(w, req)

			// Verify the response
			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}

			// Verify the log contains the correct status code
			logged := buf.String()

			var logData map[string]interface{}
			if err := json.Unmarshal([]byte(logged), &logData); err != nil {
				t.Errorf("Failed to parse log output as JSON: %v", err)
				return
			}

			if code, ok := logData["code"]; !ok || code != float64(tt.statusCode) {
				t.Errorf("Expected code %d, got %v", tt.statusCode, code)
			}

			// Verify expected response size
			expectedSize := int64(len(tt.response))
			if size, ok := logData["size"]; !ok || size != float64(expectedSize) {
				t.Errorf("Expected size %d, got %v", expectedSize, size)
			}
		})
	}
}

func TestLogRequestHandlerPrettyLogging(t *testing.T) {
	// Test that the pretty option works without errors
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Create a custom handler with pretty logging
	fn := func(w http.ResponseWriter, r *http.Request) {
		mtr := httpsnoop.CaptureMetrics(dummyHandler, w, r)

		logHTTPReqInfo(logger, &HTTPReqInfo{
			method:    r.Method,
			path:      r.URL.String(),
			code:      mtr.Code,
			size:      mtr.Written,
			duration:  mtr.Duration,
			ipAddress: requestGetRemoteAddress(r),
			userAgent: r.Header.Get("User-Agent"),
			referer:   r.Header.Get("Referer"),
		})
	}
	handler := http.HandlerFunc(fn)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify the handler still works correctly
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("Expected body 'ok', got '%s'", w.Body.String())
	}

	// Verify that we got some log output (text format)
	logged := buf.String()
	if len(logged) == 0 {
		t.Errorf("Expected log output, got empty string")
	}

	// For text format, just verify key information is present
	if !strings.Contains(logged, "GET") {
		t.Errorf("Expected log to contain method 'GET', got: %s", logged)
	}
	if !strings.Contains(logged, "HTTP Request") {
		t.Errorf("Expected log to contain message 'HTTP Request', got: %s", logged)
	}
}
