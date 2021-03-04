package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
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

	t.templ.Execute(w, nil) // 응답으로 템플릿을 보낸다.
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"}) // 경로에 요청이 오는지 수신 대기(요청이 오면 HTML 보내기)
	http.Handle("/room", r)

	// 방을 가져옴
	go r.run() // 고루틴을 통해 채팅 작업을 백그라운드에서 실행(메인하고 같이 돌고 run이 무한루프for문이므로 계속 돈다.)

	// 	웹 서버 시작
	err := http.ListenAndServe(":8080", nil) // 8080 포트에서 웹 서버 시작
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
