package main

import (
	"goweb/web5/myapp"

	"net/http"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHandler())

}
