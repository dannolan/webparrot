package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/dannolan/webparrot/handlers"
)

type key int

const (
	requestIDKey key = 0
)

var defaultAddress = ":5000"

func getHostingAddress() string {
	address := os.Getenv("HOST_ADDRESS")
	if len(address) != 0 {
		return address
	}
	return defaultAddress
}

func runningInProduction() bool {
	inProd := os.Getenv("PRODUCTION_ENV")
	if len(inProd) > 0 {
		return true
	}
	return false
}

func hostDomain() string {
	hostDomain := os.Getenv("PRODUCTION_DOMAIN")
	if len(hostDomain) > 0 {
		return hostDomain
	}
	return "example.com"
}

func main() {
	logger := log.New(os.Stdout, "webparrot:", log.LstdFlags)
	logger.Println("Loading...")
	listenAddr := getHostingAddress()

	router := http.NewServeMux()
	router.Handle("/healthz", handlers.HealthHandler{
		Logger: logger,
	})
	router.Handle("/api/v1/parrot", handlers.ParrotHandler{
		Logger: logger,
	})

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	var server *http.Server

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostDomain()), //Your domain here
		Cache:      autocert.DirCache("certs"),           //Folder for storing certificates
	}

	if runningInProduction() {
		server = &http.Server{
			Addr:         ":https",
			Handler:      tracing(nextRequestID)(logging(logger)(router)),
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		log.Fatal(server.ListenAndServeTLS("", ""))

	} else {
		server = &http.Server{
			Addr:         listenAddr,
			Handler:      tracing(nextRequestID)(logging(logger)(router)),
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		}
		log.Fatal(server.ListenAndServe())

	}

}

//SHAMELESSLY STOLEN FROM https://gist.github.com/enricofoltran/10b4a980cd07cb02836f70a4ab3e72d7

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
