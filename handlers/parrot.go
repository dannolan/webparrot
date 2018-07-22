package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

//ParrotHandler - handles parrot requests
type ParrotHandler struct {
	Logger *log.Logger
}

//ParrotResponse - parrots a response back
type ParrotResponse struct {
	IPAddress  string                 `json:"ip_address"`
	Headers    map[string]interface{} `json:"headers"`
	Host       string                 `json:"host"`
	RemoteAddr string                 `json:"remote_address"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

// FromRequest return client's real public IP address from http request headers.
// shamelessly nicked from here and changed - https://github.com/tomasen/realip/blob/master/realip.go
func FromRequest(r *http.Request) string {
	// Fetch header value
	xRealIP := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	// If both empty, return IP from remote address
	if len(xRealIP) == 0 && len(xForwardedFor) == 0 {
		var remoteIP string
		// If there are colon in remote address, remove the port number
		// otherwise, return remote address as is
		if strings.ContainsRune(r.RemoteAddr, ':') {
			remoteIP, _, _ = net.SplitHostPort(r.RemoteAddr)
		} else {
			remoteIP = r.RemoteAddr
		}

		return remoteIP
	}

	fmt.Println("checking x forwarded for")
	// Check list of IP in X-Forwarded-For and return the first global address
	for _, address := range strings.Split(xForwardedFor, ",") {
		address = strings.TrimSpace(address)
		if address != "" {
			return address
		}

	}
	// If nothing succeed, return X-Real-IP
	return xRealIP
}

func (p ParrotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	remoteAddress := r.RemoteAddr
	host := r.Host
	ip := FromRequest(r)
	headers := make(map[string]interface{})
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]

		}
	}
	parrot := ParrotResponse{
		IPAddress:  ip,
		Headers:    headers,
		RemoteAddr: remoteAddress,
		Host:       host,
	}

	presp, err := json.Marshal(&parrot)
	if err != nil {
		p.Logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		body, _ := json.Marshal(ErrorMessageResponse{Error: err.Error()})
		w.Write(body)
		return
	}
	w.Write(presp)
	return

}
