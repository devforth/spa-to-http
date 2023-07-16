package util

import (
	"net/http"
	"testing"
)

func TestIPAddrFromRemoteAddr(t *testing.T) {
	tests := []struct {
		remoteAddr string
		expected   string
	}{
		{"[::1]:58292", "[::1]"},
		{"127.0.0.1:12345", "127.0.0.1"},
		{"[::1]", "[::1]"},
		{"127.0.0.1", "127.0.0.1"},
	}

	for _, tt := range tests {
		actual := ipAddrFromRemoteAddr(tt.remoteAddr)
		if actual != tt.expected {
			t.Errorf("ipAddrFromRemoteAddr(%s): expected %s, got %s", tt.remoteAddr, tt.expected, actual)
		}
	}
}

func TestRequestGetRemoteAddress(t *testing.T) {
	tests := []struct {
		headerRealIP       string
		headerForwardedFor string
		remoteAddr         string
		expected           string
	}{
		{"", "", "127.0.0.1:12345", "127.0.0.1"},
		{"", "192.168.0.1, 127.0.0.1", "127.0.0.1:12345", "192.168.0.1"},
		{"192.168.0.1", "", "127.0.0.1:12345", "192.168.0.1"},
		{"192.168.0.1", "192.168.0.2, 127.0.0.1", "127.0.0.1:12345", "192.168.0.2"},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("X-Real-Ip", tt.headerRealIP)
		req.Header.Set("X-Forwarded-For", tt.headerForwardedFor)
		req.RemoteAddr = tt.remoteAddr

		actual := requestGetRemoteAddress(req)
		if actual != tt.expected {
			t.Errorf("requestGetRemoteAddress(%s, %s, %s): expected %s, got %s", tt.headerRealIP, tt.headerForwardedFor, tt.remoteAddr, tt.expected, actual)
		}
	}
}
