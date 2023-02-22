package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	//cloud.google.com/go 도 설치
)

type GoogleUserId struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

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
		errMsg := fmt.Sprintf("invalid google oauth state cookie: %s state:%s\n", oauthstate.Value, r.FormValue("state"))
		log.Printf(errMsg)
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(error.Error(err))
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store Id info into Session cookie
	// google에서 return 해주는 json 형태의 struct 생성 및 선언
	var userInfo GoogleUserId
	// unmarshal 을 통해 data get
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, "session")
	// Set some session values.
	session.Values["id"] = userInfo.ID

	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// fmt.Fprintf(w, string(data))
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
