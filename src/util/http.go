package util

import (
	"net"
	"net/http"
	"strings"
)

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	// return full string for IPv6 inputs wihtout port
	if strings.LastIndex(s, "]") == len(s)-1 {
		return s
	}

	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}

	return s[:idx]
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func requestGetRemoteAddress(r *http.Request) net.IP {
	hdr := r.Header

	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return net.ParseIP(ipAddrFromRemoteAddr(r.RemoteAddr))
	}

	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		fwdIPs := make([]net.IP, len(parts))
		for i, p := range parts {
			fwdIPs[i] = net.ParseIP(ipAddrFromRemoteAddr(strings.TrimSpace(p)))
		}
		// return first address
		return fwdIPs[0]
	}

	return net.ParseIP(hdrRealIP)
}
