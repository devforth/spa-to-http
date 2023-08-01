package util

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func logHTTPReqInfo(ri *HTTPReqInfo) {
	log.Info().
		Str("method", ri.method).
		Str("path", ri.path).
		Int("code", ri.code).
		Int64("size", ri.size).
		Dur("duration", ri.duration).
		IPAddr("ipAddress", ri.ipAddress).
		Str("userAgent", ri.userAgent).
		Str("referer", ri.referer).
		Send()
}

func LogRequestHandler(h http.Handler, opt *LogRequestHandlerOptions) http.Handler {
	if opt.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.DurationFieldUnit = time.Millisecond

	fn := func(w http.ResponseWriter, r *http.Request) {
		// runs handler h and captures information about HTTP request
		mtr := httpsnoop.CaptureMetrics(h, w, r)

		logHTTPReqInfo(&HTTPReqInfo{
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
