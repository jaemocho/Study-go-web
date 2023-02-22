WEB1 : 기본 WEB
WEB1_1 : 기본 WEB/JSON/MUX/TEST 
WEB4: 기본 WEB/JSON/MUX/TEST/JSONTEST 
WEB4-1 : FILE 전송/TEST
WEB5 : RestAPI CRUD
WEB6 : Decorator pattern
Decorator : Decorator pattern example 압호화 -> 압축 -> 전송
WEB7 : Decorator pattern 적용한 handler example
WEB8, template : template, 
WEB9 : pat, render, negroni 사용 예제 
WEB10 : EventSource(SSE)를 이용한 채팅 프로그램/ 
        gopkg.in/antage/eventsource.v1 사용
        channel 사용 
WEB11 : Oauth(google), pat
WEB12 : todo web, bootstrap, gorilla, render 
WEB13 : todo web, bootstrap, gorilla, render/TEST 코드 작성, model 분리 (refactoring)
WEB14 : sqlite3(tdm-gcc 설치 필요), 
       tdm-gcc 설치 후 MinGW command prompt 실행한 곳에서 go get github.com/mattn/go-sqlite3
        C:\Users\조재모\go\src\goWeb\WEB14>go get github.com/mattn/go-sqlite3
        go: downloading github.com/mattn/go-sqlite3 v1.14.16
        go: added github.com/mattn/go-sqlite3 v1.14.16
WEB15: WEB14에 google login 추가 및 세션관리(gorilla/sessions 사용 https://github.com/gorilla/sessions)
        session 별 데이터 저장 (여기서 session 은 구글 id )
WEB16 : websocket (github.com/gorilla/websocket) echo server

※ oauth
https://console.cloud.google.com/ 에서 project 생성 및 인증정보 획득

1. 프로젝트 생성 
2. OAuth 동의 화면
3. 사용자 인증정보 -> 사용자 인증 정보 만들기 
4. 승인된 리디렉션 : http://localhost:3000/auth/google/callback
5. 획득한 client_id, client_secret os 환경변수로 등록(사용자 변수)
  GOOGLE_CLIENT_ID
  GOOGLE_SECRET_KEY


{"web":{"client_id":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxx","project_id":"golangstudy-jjm","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"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"}}