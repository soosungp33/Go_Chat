package main

type room struct {
	// forward는 수신 메시지를 보관하는 채널이며 수신한 메시지는 다른 클라이언트로 전달돼야 한다
	forward chan []byte
}
