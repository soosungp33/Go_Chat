package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

// avatars 폴더에 업로드한 이미지를 저장
func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userID := req.FormValue("userid")               // HTML 폼에 숨겨진 입력에 배치한 사용자 ID를 가져온다. <input type="hidden" name="userid" value="{{.UserData.userid}}" /> 이부분
	file, header, err := req.FormFile("avatarFile") // 파일 자체(io.Reader타입), 메타데이터를 포함하는 파일 헤더, 오류 -> 파일 업로드칸에 들어오는 파일
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file) // 모든 바이트가 수신될 때까지 계속 읽는다.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("avatars", userID+path.Ext(header.Filename)) // userID로 새 파일명을 만들고 headr에서 가져올 수 있는 원래 파일명의 확장자를 복사한다.
	err = ioutil.WriteFile(filename, data, 0777)                       // avatars 폴더에 새 파일을 만드는데 userID를 사용해 gravatar와 같은 방식으로 사용자에게 이미지를 연결시킨다.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 결론적으로 고유 ID.확장자 로 저장된다.
	io.WriteString(w, "Successful")
}
