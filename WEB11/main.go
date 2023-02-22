package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/pat"

	"github.com/urfave/negroni"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	//cloud.google.com/go 도 설치
)

// google login 을 위한 config struct 선언
// os.Getenv 로 os에 설정한 환경변수 값을 이용
var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	// user email 에 접근하겠다 permission 요청
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
	// google login 시 정해진 endpoint
	Endpoint: google.Endpoint,
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// cookie 생성 및 uniq 로 생성된 id 반환
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	// 만료 시간 지정
	expiration := time.Now().Add(1 * 24 * time.Hour)

	// uniq key 생성 (random)
	b := make([]byte, 16)
	rand.Read(b)

	// string으로 변환
	state := base64.URLEncoding.EncodeToString(b)

	// cookie struct 생성 name, value는 필수, 위에 생성한 uniq key를 value로 사용
	cookie := &http.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: expiration,
	}

	// 생성한 cookie set
	http.SetCookie(w, cookie)

	return state
}

// 로그인 후 /auth/google/callback로 redirect되면서 호출되는 callback function
// user info를 받아서 화면에 표시
func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")

	// cookie 정보가 잘못된 경우 기본 경로로 redirect
	if r.FormValue("state") != oauthstate.Value {
		log.Printf("invalid google oauth state cookie: %s state:%s\n", oauthstate.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(error.Error(err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	fmt.Fprintf(w, string(data))
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	// context thread safe 한 저장공간, 기본 제공하는 Background 이용
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s", err.Error())
	}

	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s", err.Error())
	}

	return ioutil.ReadAll(resp.Body)

}

func main() {
	mux := pat.New()

	// login page
	mux.HandleFunc("/auth/google/login", googleLoginHandler)

	// login page call back (로그인 후 )
	mux.HandleFunc("/auth/google/callback", googleAuthCallback)

	// 파일 서버가 기본적으로 지원 되기 때문에 public/index.html 접근 가능
	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}

/*

로그인 후 반환해주는 id를 user key 로 사용하면 된다
email/프로필사진/친구목록 이런것들은 google에 요청해서 사용

{
  "id": "104288306882129194173",
  "email": "mot882000@gmail.com",
  "verified_email": true,
  "picture": "https://lh3.googleusercontent.com/a/default-user=s96-c"
}

// console

PS C:\Users\조재모\go\src\goWeb\web11> .\web11.exe
[negroni] 2023-02-14T11:52:37+09:00 | 304 |      108.0715ms | localhost:3000 | GET /
[negroni] 2023-02-14T11:52:39+09:00 | 307 |      1.0278ms | localhost:3000 | GET /auth/google/login
[negroni] 2023-02-14T11:52:43+09:00 | 200 |      677.337ms | localhost:3000 | GET /auth/google/callback

*/
