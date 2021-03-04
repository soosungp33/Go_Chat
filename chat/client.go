package main

import (
	"github.com/gorilla/websocket"
)

type client struct { // client는 한 명의 채팅 사용자를 나타낸다.
	socket *websocket.Conn // socket은 이 클라이언트의 웹 소켓이다(클라이언트와 통신할 수 있는 웹 소켓에 대한 참조)
	send   chan []byte     // send는 메시지가 전송되는 채널
	room   *room           // room은 클라이언트가 채팅하는 방
}

// 글을 쓰면 소켓에 글이 들어감.
// read 메소드에서 소켓에 있는 글을 읽고 forward 채널로 메시지를 전송한다.
// forward 채널에 메시지가 전송되면 그 메시지를 모든 클라이언트의 send 채널에 메시지를 추가한다.
// write 메소드에서 각 클라이언트는 send 채널에 의해 메시지를 기다리고 있다가 send 채널에 온 메시지를 수신한다.
func (c *client) read() {
	defer c.socket.Close()
	for { // 무한루프
		_, msg, err := c.socket.ReadMessage() // 소켓에서 읽고
		if err != nil {
			return
		}

		c.room.forward <- msg // room의 forward 채널로 계속 전송
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg) // 소켓에서 메시지를 계속 수신
		if err != nil {
			return
		}
	}
}
