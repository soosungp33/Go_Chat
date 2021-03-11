package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
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
	case "login": // 사용자에게 권한 부여
		provider, err := gomniauth.Provider(provider) // URL에 지정된 객체(google or github 등)와 일치하는 프로바이더 객체를 가져온다.
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil) // 인증 프로세스를 시작하기 위해 사용자를 보내야 하는 위치를 가져온다.
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", loginUrl) // GetBeginAuthURL 호출시 오류가 없으면 사용자의 브라우저를 반환된 URL로 리디렉션한다.
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback": // 사용자에게 권한을 부여한 후 리다이렉션하면 이 case로 온다.
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery)) // URL을 파싱해서 OAuth2 핸드셰이크를 완료한다.(자격증명을 받음)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err), http.StatusInternalServerError)
			return
		}
		user, err := provider.GetUser(creds) // 제공자에 대해 자격증명 정보를 사용해 사용자에 대한 몇 가지 기본 정보에 액세스한다.
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err), http.StatusInternalServerError)
			return
		}
		authCookieValue := objx.New(map[string]interface{}{ // 사용자가 있으면 JSON 객체의 Name 필드를 Base64로 인코딩한다.(Base64는 데이터를 URL이나 쿠키에 저장하는 경우 유용하다.)
			"name":       user.Name(),      // 사용자명
			"avatar_url": user.AvatarURL(), // 사용자 사진
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{ // 나중에 사용할 수 있도록 auth 쿠키 값으로 저장한다.(func (h *authHandler) ServeHTTP 메소드에서 사용)
			Name:  "auth",
			Value: authCookieValue, // auth의 value 값에 user name이 저장되어 있다.
			Path:  "/"})

		w.Header().Set("Location", "/chat") // 원래 목적지인 chat으로 리다이렉션
		w.WriteHeader(http.StatusTemporaryRedirect)

	default: // 아니면 오류 메시지 출력
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
