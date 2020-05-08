package isso

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func jsonBind(r io.ReadCloser, obj interface{}) error {
	defer r.Close()

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return fmt.Errorf("invalid JSON payload: %v", err)
	}

	return nil
}

func findClientIP(r *http.Request) string {
	headers := []string{"X-Forwarded-For", "X-Real-Ip"}
	for _, header := range headers {
		value := r.Header.Get(header)

		if value != "" {
			addresses := strings.Split(value, ",")
			address := strings.TrimSpace(addresses[0])

			if net.ParseIP(address) != nil {
				return address
			}
		}
	}

	// Fallback to TCP/IP source IP address.
	var remoteIP string
	if strings.ContainsRune(r.RemoteAddr, ':') {
		remoteIP, _, _ = net.SplitHostPort(r.RemoteAddr)
	} else {
		remoteIP = r.RemoteAddr
	}

	// When listening on a Unix socket, RemoteAddr is empty.
	if remoteIP == "" {
		remoteIP = "127.0.0.1"
	}

	return remoteIP
}

func findOrigin(r *http.Request) string {
	origin := r.Header.Get("origin")
	if origin == "" {
		u, err := url.Parse(r.Referer())
		if err != nil {
			return ""
		}
		return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}
	return origin
}
