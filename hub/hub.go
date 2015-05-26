package hub

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type conn struct {
	ws      *websocket.conn
	send    chan []byte
}

type msg struct {
	sender *conn
	data   []byte
}

type hub struct {
	conns      map[conn]struct{}
	register   chan *conn
	unregister chan *conn
	receive    chan msg
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.conns[c] = struct{}{}
		case c := <-h.unregister:
			if _, ok := h.conns[c]; ok {
				delete(h.conns, c)
				close(c.send)
				close(c.receive)
			}
		case msg := <-h.receive:
			// TODO: Process message data
			fmt.Println(string(msg.data))
		}
	}
}
