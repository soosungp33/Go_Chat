package main

import (
	// 해시 패키지
	"errors"
	"io/ioutil"
	"path"
)

// ErrNoAvatar는 Avatar 인스턴스가 아바타 URL을 제공할 수 없을 때 리턴되는 에러다
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

// Avatar는 사용자 프로필 사진을 표현할 수 있는 타입을 나타낸다.
type Avatar interface {
	// GetAvatarURL은 지정된 클라이언트에 대한 아바타 URL을 가져오고, 문제가 발생하면 에러를 리턴한다.
	// 객체가 지정된 클라이언트의 URL을 가져올 수 없는 경우 ErrNoAvatarURL이 리턴된다.
	GetAvatarURL(c *client) (string, error) // URL을 리턴할 사용자를 알 수 있도록 클라이언트를 인수로 사용
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

// 객체가 nil값을 가질 수 있으므로 리시버의 변수명을 생략해 Go에 참조를 버리라고 전달
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok { // 사용자 데이터가 있으면
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar와 같은 패턴
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

// 객체가 nil값을 가질 수 있으므로 리시버의 변수명을 생략해 Go에 참조를 버리라고 전달
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok { // 사용자 데이터가 있으면
			return "//www.gravatar.com/avatar/" + useridStr, nil // 기본 URL에 추가
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			files, err := ioutil.ReadDir("avatars") // avatars 폴데 있는 파일 목록을 전부 가져온다.(avatars 디렉터리도 포함됨) -> 뒤에 확장자까지 다 저장됨
			if err != nil {
				return "", ErrNoAvatarURL
			}
			for _, file := range files {
				if file.IsDir() { // 목록에는 디렉터리도 포함되므로 디렉터리면 건너뛴다.
					continue
				}
				if match, _ := path.Match(useridStr+"*", file.Name()); match { // 각 파일이 userid와 일치하는지 확인
					return "/avatars/" + file.Name(), nil // 일치하면 해당 파일을 찾은 것이므로 경로를 리턴
				}
			}
			//return "/avatars/" + useridStr + ".jpg", nil // 하드코딩 -> jpg만 지원
		}
	}

	return "", ErrNoAvatarURL
}
