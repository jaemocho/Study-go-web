package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at`
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "jjm", Email: "jjm@naver.com"}

	// 아래 네줄을 rd를 사용해서 한줄로 치환
	rd.JSON(w, http.StatusOK, user)

	// w.Header().Add("content-type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// data, _ := json.Marshal(user)
	// fmt.Fprint(w, string(data))

}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {

		// rd를 사용해 아래 두줄을 한줄로 치환
		rd.Text(w, http.StatusBadRequest, err.Error())

		// w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprint(w, err)
		return
	}
	user.CreatedAt = time.Now()

	// 아래 네줄을 rd를 사용해서 한줄로 치환
	rd.JSON(w, http.StatusOK, user)

	// w.Header().Add("content-type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// data, _ := json.Marshal(user)
	// fmt.Fprint(w, string(data))

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "jjm", Email: "jjm@naver.com"}

	// tmpl, err := template.New("Hello").ParseFiles("templates/hello.tmpl")
	// if err != nil {
	// 	// rd를 사용해 아래 두줄을 한줄로 치환
	// 	rd.Text(w, http.StatusBadRequest, err.Error())

	// 	// w.WriteHeader(http.StatusBadRequest)
	// 	// fmt.Fprint(w, err)
	// 	return
	// }
	// tmpl.ExecuteTemplate(w, "hello.tmpl", "jjm")

	// 위 주석처리된 전체를 아래 한줄로 치환 가능 render package가 확장자 빼고 등록해서 hello 라고만 입력
	rd.HTML(w, http.StatusOK, "body", user)

	// template 폴더아래 tmpl 로 넣으면 인식
	// render 선언 시 지정할 수 있음 main func 참고

	// hello.html에 {{ yield }} 선언하고 render 선언 시 Layout으로 지정
	// 위에 rendering 할 template은 body로 지정하고 실행하면
	// body 가 layout의 yield 부분으로 들어간다

	// hello.html 에 {{ partial "title" }} 과 같이 선언한 후
	// title-body.html 을 만들어주면 안의 내용이 parsing 되어서 들어간다
}

func main() {
	rd = render.New(render.Options{
		Directory:  "templates",
		Extensions: []string{".html", ".tmpl"},
		Layout:     "hello",
	})
	// mux := http.NewServeMux()
	mux := pat.New()

	mux.Get("/users", getUserInfoHandler)
	mux.Post("/users", addUserHandler)
	mux.Get("/hello", helloHandler)

	// negroni file server/log 기능 기본 제공
	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", mux)

}
