package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type templateHandler struct { // 템플릿을 로드하고 컴파일하며 전달하는 구조체
	once     sync.Once // 함수를 한 번만 실행하기 위해 사용
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { // templateHandler 타입에는 ServeHTTP라는 단일 메소드 존재
	// ServeHTTP 메소드는 소스 파일을 로드하고, 템플릿을 컴파일한 후 실행하고 지정된 http.ResponseWriter 메소드에 출력을 작성한다.
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	// 전체 http.Request 객체를 전달하는 것 대신에 Host 및 UserData가 있는 데이터를 만들어 전달한다.
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value) // authCookie.Value는 user name이 저장되어 있다.
	}

	t.templ.Execute(w, data) // 응답으로 템플릿을 보낸다.
	// r을 data인수로 전달하므로써 http.Request에서 추출할 수 있는 데이터를 사용해 템플릿을 표시하도록 지시(호스트 주소가 포함됨)
	// 따라서 chat.html 파일의 소켓 생성하는 라인에서 {{.Host}}를 사용할 수 있다.
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.") // *string 타입을 반환(주소)
	flag.Parse()                                                            // 플래그 파싱
	// gomniauth 설정
	gomniauth.SetSecurityKey("PUT YOUR AUTH KEY HERE")
	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret", "http://localhost:8080/auth/callback/github"),
		google.New("1084570662586-jq97d3vj6vmtf2919g567p4at6c7qqrf.apps.googleusercontent.com", "HA9ze1I0TYSN7lFRkoZsl9u_", "http://localhost:8080/auth/callback/google"),
	)

	// r:= newRoom() // 프로필 사진 x
	//r := newRoom(UseAuthAvatar) // 프로필 사진 o
	//r := newRoom(UseGravatar) // 프로필 사진 gravatar 이미지로 변경
	r := newRoom(UseFileSystemAvatar) // 프로필 사진 업로드 가능

	//r.tracer = trace.New(os.Stdout)                           // 추적 결과를 터미널로 출력하고 싶을 때 사용(Trace의 t에 쓰인 내용이 터미널에 나옴)
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"})) // 경로에 요청이 오는지 수신 대기(요청이 오면 HTML 보내기), 채팅
	// MustAuth는 authHandler를 통한 권한 수행이 먼저 실행되고 인증되면 templateHandler가 실행된다.
	http.Handle("/login", &templateHandler{filename: "login.html"}) // 로그인
	http.HandleFunc("/auth/", loginHandler)                         // 권한 요청
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) { // 로그아웃
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "", // 빈 문자열을 넣어 이전에 저장돼 있던 사용자 데이터를 제거한다.
			Path:   "/",
			MaxAge: -1, // 브라우저에서 즉시 삭제돼야 함을 나타낸다.
		})
		w.Header().Set("Location", "/chat") // 로그아웃 후 채팅 페이지로 리다이렉션하면 채팅 페이지에서 로그인 페이지로 리다이렉션 된다.
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler) // 업로드 핸들러 매핑

	http.Handle("/avatars/",
		http.StripPrefix("/avatars/", // 지정된 접두사를 제거해 경로를 수정한 후 핸들러로 전달(제거하지 않으면 /avatars/avatars/filename과 같은 경로가 된다.)
			http.FileServer(http.Dir("./avatars")))) // 공개할 폴더를 지정

	// 방을 가져옴
	go r.run() // 고루틴을 통해 채팅 작업을 백그라운드에서 실행(메인하고 같이 동시에 돌고 run이 무한루프for문이므로 계속 돈다.)

	// 	웹 서버 시작
	log.Println("Starting web server on", *addr)
	err := http.ListenAndServe(*addr, nil) // 8080 포트에서 웹 서버 시작
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
