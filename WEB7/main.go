package main

import (
	"goweb/web7/decoHandler"
	"goweb/web7/myapp"
	"log"
	"net/http"
	"time"
)

// handler 의 decorator
func logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Print("[LOGGER1] Started")
	h.ServeHTTP(w, r)
	log.Println("[LOGGER1] Completed time", time.Since(start).Milliseconds())

}

func logger2(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Print("[LOGGER2] Started")
	// param 전달 역할 w, r 에 대한 컨트롤이 가능
	h.ServeHTTP(w, r)
	log.Println("[LOGGER2] Completed time", time.Since(start).Milliseconds())

}

func NewHandler() http.Handler {
	mux := myapp.NewHandler()
	h := decoHandler.NewDecoHandler(mux, logger)
	h = decoHandler.NewDecoHandler(h, logger2)
	return h
}

func main() {
	mux := NewHandler()

	http.ListenAndServe(":3000", mux)
}
