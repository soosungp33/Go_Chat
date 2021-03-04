package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/soosungp33/Go_MSA/trace"
)

type room struct {
	forward chan []byte // forward는 수신 메시지를 보관하는 채널이며 수신한 메시지는 다른 클라이언트로 전달돼야 한다
	// join과 leave는 clients 맵에서 클라이언트를 안전하게 추가 및 제거하기 위해 존재
	join    chan *client     // 방에 들어오려는 클라이언트를 위한 채널
	leave   chan *client     // 방을 나가길 원하는 클라이언트를 위한 채널
	clients map[*client]bool // 현재 채팅방에 있는 모든 클라이언트를 보유
	tracer  trace.Tracer     // tracer는 방 안에서 활동의 추적 정보를 수신한다.
}

func newRoom() *room { // 채팅방 만드는 함수
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select { // 한 번에 한 케이스 코드만 실행되므로 맵이 동시에 여러 개 수정되는 가능성을 방지하며 동기화한다.
		case client := <-r.join: // join 채널에서 메시지를 받으면
			// 입장
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave: // leave 채널에서 메시지를 받으면
			// 퇴장
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward: // forward 채널에서 메시지를 받으면
			// 모든 클라이언트에게 메시지 전달
			for client := range r.clients {
				client.send <- msg // 각 클라이언트의 send 채널에 메시지를 추가하고 클라이언트 타입의 write 메소드가 이를 받아들여 소켓에서 브라우저로 보낸다.
				r.tracer.Trace(" -- set to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// 웹 소켓을 사용하려면 websocket.Upgrader 타입을 사용해 HTTP 연결을 업그레이드 해야 한다.(재사용 가능)
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil) // 소켓 가져오기
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	client := &client{ //  문제가 없다면 클라이언트 생성
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client // 생성한 클라이언트를 join채널에 전달
	defer func() { r.leave <- client }()
	go client.write() // 고루틴으로 클라이언트의 write 메소드를 호출
	client.read()     // 메인 스레드에서 read 메소드를 호출해 닫을 때까지 작업을 차단(연결을 활성 상태로 유지)
}
