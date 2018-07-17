package handlers

import (
	"encoding/json"
	"errors"
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

func forwardedIP(r *http.Request) (string, error) {
	forwarded := r.Header.Get("X-Forwarded-For")
	if len(forwarded) == 0 {
		return "", errors.New("No forwarded IP provided")
	}
	items := strings.Split(forwarded, ",")
	// Technically this validates that it's actually an IP
	ip := net.ParseIP(items[0])
	return ip.String(), nil
}

func (p ParrotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	remoteAddress := r.RemoteAddr
	host := r.Host
	ip, err := forwardedIP(r)
	if err != nil {
		ip = remoteAddress
	}
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
