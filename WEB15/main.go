package main

import (
	"goWeb/WEB15/app"
	"log"
	"net/http"
)

func main() {
	// sqlight 사용
	m := app.MakeHandler("./test.db")

	// postgres 사용
	//m := app.MakeHandler(os.Getenv("DATABASE_URL"))

	defer m.Close()
	// public 밑에 폴더가 기본 위치

	err := http.ListenAndServe(":3000", m)

	log.Println("Started App")
	if err != nil {
		panic(err)
	}

}
