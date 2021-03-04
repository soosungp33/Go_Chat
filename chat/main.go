package main

import (
	"flag"
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

	t.templ.Execute(w, r) // 응답으로 템플릿을 보낸다.
	// r을 data인수로 전달하므로써 http.Request에서 추출할 수 있는 데이터를 사용해 템플릿을 표시하도록 지시(호스트 주소가 포함됨)
	// 따라서 chat.html 파일의 소켓 생성하는 라인에서 {{.Host}}를 사용할 수 있다.
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.") // *string 타입을 반환(주소)
	flag.Parse()                                                            // 플래그 파싱

	r := newRoom()
	//r.tracer = trace.New(os.Stdout)                           // 추적 결과를 터미널로 출력하고 싶을 때 사용(Trace의 t에 쓰인 내용이 터미널에 나옴)
	http.Handle("/", &templateHandler{filename: "chat.html"}) // 경로에 요청이 오는지 수신 대기(요청이 오면 HTML 보내기)
	http.Handle("/room", r)

	// 방을 가져옴
	go r.run() // 고루틴을 통해 채팅 작업을 백그라운드에서 실행(메인하고 같이 동시에 돌고 run이 무한루프for문이므로 계속 돈다.)

	// 	웹 서버 시작
	log.Println("Starting web server on", *addr)
	err := http.ListenAndServe(*addr, nil) // 8080 포트에서 웹 서버 시작
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
