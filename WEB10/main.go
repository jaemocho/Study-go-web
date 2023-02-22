package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/pat"

	"github.com/urfave/negroni"

	eventsource "gopkg.in/antage/eventsource.v1"
)

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	log.Println("postMessageHandler", msg, name)
	sendMessage(name, msg)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	sendMessage("", fmt.Sprintf("add user: %s", username))
}

func leftUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	sendMessage("", fmt.Sprintf("left user: %s", username))
}

type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

var msgCh chan Message

// message 생성/ 사용자 입장/퇴장 시 channel에 데이터 put
func sendMessage(name, msg string) {
	// send message to every clients
	msgCh <- Message{name, msg}
}

func processMsgCh(es eventsource.EventSource) {
	// channel 에 data가 들어올 때마다 eventsource message 전송(구독자들에게 )
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}

func main() {
	msgCh = make(chan Message)

	es := eventsource.New(nil, nil)
	defer es.Close()

	// event 발생 시 마다 처리하기 위해 go routine 생성
	go processMsgCh(es)

	mux := pat.New()
	mux.Post("/messages", postMessageHandler)
	mux.Post("/users", addUserHandler)
	mux.Delete("/users", leftUserHandler)

	// sse(server sent event)
	mux.Handle("/stream", es)

	n := negroni.Classic()
	// file server때문에 사용 public/index.html, chat.js 사용하려고
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
