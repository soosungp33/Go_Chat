package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // 경로에 요청이 오는지 수신 대기(요청이 오면 HTML 보내기)
		w.Write([]byte(`
		<html>
		  <head>
		    <title>Chat</title>
		  </head>
		  <body>
		    Let's chat!
		  </body>
		</html>
		`))
	})

	// 	웹 서버 시작
	err := http.ListenAndServe(":8080", nil) // 8080 포트에서 웹 서버 시작
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
