package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
