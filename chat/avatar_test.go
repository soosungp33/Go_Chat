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
