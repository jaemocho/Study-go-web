package myapp

// go get -u github.com/gorilla/mux
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var userMap map[int]*User
var lastId int

// User struct
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

/*
 빈 값과 default 값 구분이 어려워서
 아래와 같이 bool type으로 만들어서 해결하는 방법도 있다.

type UpdateUser struct {
	ID        int       `json:"id"`
	UpdatedFirstName bool `json:"updated_first_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
*/

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Get UserInfo by /users/{id}")
	if len(userMap) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No Users")
		return
	}

	// User struct pointer slice 형태로 선언
	users := []*User{}
	// map에서 꺼내와서 위에 선언한 attr에 append
	for _, u := range userMap {
		users = append(users, u)
	}

	// json 형태로 marsahling
	data, _ := json.Marshal(users)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(data))
}

func getUserInfo89Handler(w http.ResponseWriter, r *http.Request) {
	// gorilla mux 기능
	// vars := mux.Vars(r)
	// fmt.Fprint(w, "User Id:", vars["id"])
	vars := mux.Vars(r) // vars는 map[string]string

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	user, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func createUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	// Created User
	lastId++
	user.ID = lastId
	user.CreatedAt = time.Now()
	userMap[user.ID] = user
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))

}

func deleteUserInfo89Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	_, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID:", id)
		return
	}

	delete(userMap, id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User ID:", id)

}

func updateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	updateUser := new(User)
	err := json.NewDecoder(r.Body).Decode(updateUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[updateUser.ID]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID:", updateUser.ID)
		return
	}

	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}

	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}

	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	// userMap[updateUser.ID] = user // pointer type이라 굳이 이렇게 안해줘도 된다.

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(updateUser)
	fmt.Fprint(w, string(data))
}

// Newhandler make a new myapp handler
func NewHandler() http.Handler {

	userMap = make(map[int]*User)
	lastId = 0
	//mux := http.NewServeMux()
	// gorilla mux 사용
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	//Methods("") - gorilla 에서 제공 요청 마다 다르게 동작
	mux.HandleFunc("/users", usersHandler).Methods("GET")
	mux.HandleFunc("/users", createUserInfoHandler).Methods("POST")
	mux.HandleFunc("/users", updateUserInfoHandler).Methods("PUT")
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfo89Handler).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserInfo89Handler).Methods("DELETE")
	return mux
}
