package main

import (
	"goweb/web14/app"
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler("./test.db")
	defer m.Close()
	// public 밑에 폴더가 기본 위치
	n := negroni.Classic()
	n.UseHandler(m)

	err := http.ListenAndServe(":3000", n)

	log.Println("Started App")
	if err != nil {
		panic(err)
	}

}
