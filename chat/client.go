package main

import (
	"github.com/gorilla/websocket"
)

type client struct { // client는 한 명의 채팅 사용자를 나타낸다.
	socket *websocket.Conn // socket은 이 클라이언트의 웹 소켓이다(클라이언트와 통신할 수 있는 웹 소켓에 대한 참조)
	send   chan []byte     // send는 메시지가 전송되는 채널
	room   *room           // room은 클라이언트가 채팅하는 방
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		c.room.forward <- msg
	}
}
