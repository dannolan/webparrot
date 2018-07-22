package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHealthHandler(t *testing.T) {
	Convey("We should be able to hit the health handler", t, func() {
		Convey("Success case should work", func() {
			req, err := http.NewRequest("GET", "/healthz", nil)
			if err != nil {
				t.Fatal(err)
			}
			logger := log.New(os.Stdout, "Test", log.LstdFlags)

			handler := HealthHandler{
				Logger: logger,
			}
			rr := httptest.NewRecorder()

			webHandler := http.HandlerFunc(handler.ServeHTTP)

			webHandler.ServeHTTP(rr, req)

			So(rr.Code, ShouldEqual, 200)
		})
	})
}
