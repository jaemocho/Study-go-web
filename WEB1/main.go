package main

import (
	"fmt"
	"net/http"
)

type fooHandler struct {
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello foo")
}
func barHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world bar")
}

func main() {

	// function 형태 , 정적 등록
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	// function 형태
	http.HandleFunc("/bar", barHandler)

	// intance 성택
	http.Handle("/foo", &fooHandler{})

	http.ListenAndServe(":3000", nil)

}
