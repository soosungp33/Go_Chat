package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar

	client := new(client) // 사용자 데이터가 없는 클라이언트를 사용해 ErrorNoAvatarURL이 리턴되는지 확인
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	// 적절한 값을 설정해서
	testUrl := "http://url-to-gravatar/"
	client.userData = map[string]interface{}{
		"avatar_url": testUrl,
	}

	url, err = authAvatar.GetAvatarURL(client) // 올바른 값을 리턴하는지 확인
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar

	client := new(client)

	client.userData = map[string]interface{}{ // gravatar는 이메일 주소의 해시를 사용해 각 프로필 이미지의 고유 ID를 생성
		"uesrid": "0bc83cb571cd1c50ba6f3e8a78ef1346",
	}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" { // gravatar계정이 없으면 gravatar에 있는 임의의 이미지가 나온다.
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
