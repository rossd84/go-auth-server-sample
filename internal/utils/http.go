package utils

import (
	"net"
	"net/http"
	"strings"
)

type Metadata struct {
	UserAgent string
	IPAddress string
	DeviceID  string
	Location  string
	Platform  string
	Browser   string
}

func ExtractMetadata(r *http.Request) Metadata {
	return Metadata{
		UserAgent: r.UserAgent(),
		IPAddress: GetIPAddress(r),
		DeviceID:  r.Header.Get("X-Device-ID"),
		Location:  r.Header.Get("X-Location"),
		Platform:  r.Header.Get("X-Platform"),
		Browser:   r.Header.Get("X-Browser"),
	}
}

func GetIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// May contain multiple IPs
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	// Fallback to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
