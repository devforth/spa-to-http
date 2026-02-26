package util

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/felixge/httpsnoop"
)

type LogRequestHandlerOptions struct {
	Pretty bool
}

// LogReqInfo describes info about HTTP request
type HTTPReqInfo struct {
	// GET etc.
	method string
	// requested path
	path string
	// response code, like 200, 404
	code int
	// number of bytes of the response sent
	size int64
	// how long did it take to
	duration time.Duration
	// client IP Address
	ipAddress net.IP
	// client UserAgent
	userAgent string
	// referer header
	referer string
}

func logHTTPReqInfo(l *slog.Logger, ri *HTTPReqInfo) {
	l.Info("HTTP Request",
		"method", ri.method,
		"path", ri.path,
		slog.Int("code", ri.code),
		slog.Int64("size", ri.size),
		slog.Int64("duration", ri.duration.Milliseconds()), // in milliseconds
		"ipAddress", ri.ipAddress,
		"userAgent", ri.userAgent,
		"referer", ri.referer,
	)
}

func LogRequestHandler(h http.Handler, opt *LogRequestHandlerOptions) http.Handler {
	var logger *slog.Logger
	if opt.Pretty {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		// runs handler h and captures information about HTTP request
		mtr := httpsnoop.CaptureMetrics(h, w, r)

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

	return http.HandlerFunc(fn)
}
