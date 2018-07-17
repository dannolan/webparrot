package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

//HealthHandler - handles healthz responses
type HealthHandler struct {
	Logger *log.Logger
}

//HealthResponse - respond with a health response
type HealthResponse struct {
	Time string `json:"time"`
}

//ErrorMessageResponse - handles error message responses
type ErrorMessageResponse struct {
	Error string `json:"error"`
}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().UTC().Format(time.RubyDate)
	hr := HealthResponse{
		Time: currentTime,
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	resp, err := json.Marshal(&hr)
	if err != nil {
		h.Logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		body, _ := json.Marshal(ErrorMessageResponse{Error: err.Error()})
		w.Write(body)
		return
	}
	w.Write(resp)
	return
}
