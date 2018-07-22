package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParrotHandler(t *testing.T) {
	Convey("We should be able to handle all valid testing cases", t, func() {
		Convey("We should be able to handle the remote address case", func() {
			req, err := http.NewRequest("GET", "/api/v1/parrot", nil)
			req.RemoteAddr = "192.168.1.1"
			if err != nil {
				t.Fatal(err)
			}
			logger := log.New(os.Stdout, "Test", log.LstdFlags)
			handler := ParrotHandler{
				Logger: logger,
			}
			rr := httptest.NewRecorder()

			webHandler := http.HandlerFunc(handler.ServeHTTP)

			webHandler.ServeHTTP(rr, req)
			var output ParrotResponse
			bytes, _ := ioutil.ReadAll(rr.Body)
			json.Unmarshal(bytes, &output)
			So(output.IPAddress, ShouldEqual, "192.168.1.1")
			So(rr.Code, ShouldEqual, 200)

		})
		Convey("We should be able to handle the remote address with port", func() {
			req, err := http.NewRequest("GET", "/api/v1/parrot", nil)
			req.RemoteAddr = "192.168.1.1:7876"
			if err != nil {
				t.Fatal(err)
			}
			logger := log.New(os.Stdout, "Test", log.LstdFlags)
			handler := ParrotHandler{
				Logger: logger,
			}
			rr := httptest.NewRecorder()

			webHandler := http.HandlerFunc(handler.ServeHTTP)

			webHandler.ServeHTTP(rr, req)
			var output ParrotResponse
			bytes, _ := ioutil.ReadAll(rr.Body)
			json.Unmarshal(bytes, &output)
			So(output.IPAddress, ShouldEqual, "192.168.1.1")
			So(rr.Code, ShouldEqual, 200)

		})
		Convey("We should be able to handle the real IP address case", func() {
			req, err := http.NewRequest("GET", "/api/v1/parrot", nil)
			req.RemoteAddr = "192.168.1.1"
			req.Header.Set("X-Real-Ip", "192.168.1.2")
			if err != nil {
				t.Fatal(err)
			}
			logger := log.New(os.Stdout, "Test", log.LstdFlags)
			handler := ParrotHandler{
				Logger: logger,
			}
			rr := httptest.NewRecorder()

			webHandler := http.HandlerFunc(handler.ServeHTTP)

			webHandler.ServeHTTP(rr, req)
			var output ParrotResponse
			bytes, _ := ioutil.ReadAll(rr.Body)
			json.Unmarshal(bytes, &output)
			fmt.Println(output)
			So(output.IPAddress, ShouldEqual, "192.168.1.2")
			So(rr.Code, ShouldEqual, 200)

		})

		Convey("We should be able to handle the remote address with forwarded for", func() {
			req, err := http.NewRequest("GET", "/api/v1/parrot", nil)
			req.RemoteAddr = "192.168.1.1:7876"
			req.Header.Set("X-Forwarded-For", "192.168.10.1")

			if err != nil {
				t.Fatal(err)
			}
			logger := log.New(os.Stdout, "Test", log.LstdFlags)
			handler := ParrotHandler{
				Logger: logger,
			}
			rr := httptest.NewRecorder()

			webHandler := http.HandlerFunc(handler.ServeHTTP)

			webHandler.ServeHTTP(rr, req)
			var output ParrotResponse
			bytes, _ := ioutil.ReadAll(rr.Body)
			json.Unmarshal(bytes, &output)

			So(output.IPAddress, ShouldEqual, "192.168.10.1")
			logger.Println(output)
			So(rr.Code, ShouldEqual, 200)

		})

	})

}
