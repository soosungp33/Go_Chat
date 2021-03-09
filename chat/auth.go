package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth") // auth라는 특수 쿠키를 찾는다.
	if err == http.ErrNoCookie {
		// 인증 불가
		w.Header().Set("Location", "/login") // 쿠키가 없는 경우 로그인 페이지로 리다이렉션한다.
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// 다른 에러
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 성공 - 다음 핸들러 호출
	h.next.ServeHTTP(w, r)
}

// 단순히 다른 http.Handler를 저장(래핑)하는 authHandler이다.
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler는 서드파티 로그인 프로세스를 처리한다.
// 형식 : /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) { // 단순한 함수이며, handler 인터페이스를 구현하는 객체가 아니므로 http.HandleFunc를 사용
	segs := strings.Split(r.URL.Path, "/") // 경로를 "/"기준으로 나눠서 segs에 넣는다. 0에는 공백, 1에는 auth가 들어가있음
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login": // 동작 값을 알고 있으면 실행
		log.Println("TODO handle login for", provider)
	default: // 아니면 오류 메시지 출력
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
