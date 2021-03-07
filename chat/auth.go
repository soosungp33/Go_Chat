package main

import "net/http"

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
