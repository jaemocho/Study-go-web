package app

import (
	"fmt"
	"goWeb/WEB15/model"
	"os"
	"strings"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var rd *render.Render = render.New()

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

// test를 위해 function pointer를 갖는 variable로 변경
// 아래와 같은 형태로 사용
// getSessionID = func(r *http.Request) string {
//		return "testsessionId"
//	}

var getSessionID = func(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["id"]
	if val == nil {
		return ""
	}

	return val.(string)
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	list := a.db.GetTodos(sessionId)
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	name := r.FormValue("name")
	todo := a.db.AddTodo(name, sessionId)

	rd.JSON(w, http.StatusCreated, todo)
}

type Success struct {
	Success bool `json:"success"`
}

func (a *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := a.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusInternalServerError, Success{false})
	}

}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	fmt.Println(id, vars)
	complete := r.FormValue("complete") == "true"
	ok := a.db.CompleteTodo(id, complete)
	fmt.Println(ok)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusInternalServerError, Success{false})
	}

}

// func addTestTodos() {
// 	todoMap[1] = &model.Todo{1, "Buy a milk", false, time.Now()}
// 	todoMap[2] = &model.Todo{2, "Buy a coffee", true, time.Now()}
// 	todoMap[3] = &model.Todo{3, "go gym", false, time.Now()}

// }

func (a *AppHandler) Close() {
	a.db.Close()
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// if request URL is /signin.html, then next()
	if strings.Contains(r.URL.Path, "/signin") ||
		strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	// if user already signed in
	sessionID := getSessionID(r)
	if sessionID != "" {
		next(w, r)
		return
	}
	// if not user sign in
	// redirect signin.html
	// signin.html이 계속 호출되는 무한 루프를 돌 수 있기 때문에 먼저 처리
	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)

}

func MakeHandler(filepath string) *AppHandler {

	r := mux.NewRouter()
	// negroni.Classic은 3가지 decorator 가지고 있다. NewRecovery, NewLogger, NewStatic(http.Dir("public"))
	// custom 하게 사용하려면 아래와 같이 decorator를 직접 선언 chaining
	// negroni.HandlerFunc(CheckSignin) 와 같은 형태로 선언해서 사용
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.HandlerFunc(CheckSignin), negroni.NewStatic(http.Dir("public")))
	// n := negroni.Classic()

	n.UseHandler(r)

	a := &AppHandler{
		Handler: n,
		db:      model.NewDBHandler(filepath),
	}

	r.HandleFunc("/todos", a.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")
	// login page call back (로그인 후 )
	r.HandleFunc("/auth/google/callback", googleAuthCallback)
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/", a.indexHandler)

	return a
}
