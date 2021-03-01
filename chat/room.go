package main

type room struct {
	forward chan []byte // forward는 수신 메시지를 보관하는 채널이며 수신한 메시지는 다른 클라이언트로 전달돼야 한다
	// join과 leave는 clients 맵에서 클라이언트를 안전하게 추가 및 제거하기 위해 존재
	join    chan *client     // 방에 들어오려는 클라이언트를 위한 채널
	leave   chan *client     // 방을 나가길 원하는 클라이언트를 위한 채널
	clients map[*client]bool // 현재 채팅방에 있는 모든 클라이언트를 보유
}

func (r *room) run() {
	for {
		select { // 한 번에 한 케이스 코드만 실행되므로 맵이 동시에 여러 개 수정되는 가능성을 방지하며 동기화한다.
		case client := <-r.join: // join 채널에서 메시지를 받으면
			// 입장
			r.clients[client] = true
		case client := <-r.leave: // leave 채널에서 메시지를 받으면
			// 퇴장
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward: // forward 채널에서 메시지를 받으면
			// 모든 클라이언트에게 메시지 전달
			for client := range r.clients {
				client.send <- msg // 각 클라이언트의 send 채널에 메시지를 추가하고 클라이언트 타입의 write 메소드가 이를 받아들여 소켓에서 브라우저로 보낸다.
			}
		}
	}
}
