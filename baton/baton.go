package baton

import (
	"fmt"
	"time"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var h = hub{
	conns:      make(map[conn]struct{}),
	register:   make(chan conn),
	unregister: make(chan conn),
	receive:    make(chan msg),
}

// conn is the middleman between the websocket connection and the hub
type conn struct {
	ws   *websocket.Conn
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub
func (c conn) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		h.receive <- msg{c, data}
	}
}

// write writes a message with the given message type and payload
func (c conn) write(msgType int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(msgType, payload)
}

// writePump pumps messages from the hub to the websocket connection
func (c conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

type msg struct {
	conn conn
	data []byte
}

type hub struct {
	conns      map[conn]struct{}
	register   chan conn
	unregister chan conn
	receive    chan msg
}

func (h hub) run() {
	for {
		select {
		case c := <-h.register:
			h.conns[c] = struct{}{}
		case c := <-h.unregister:
			if _, ok := h.conns[c]; ok {
				delete(h.conns, c)
				close(c.send)
			}
		case msg := <-h.receive:
			// TODO: Process message data
			select {
			case msg.conn.send <- msg.data: // echo
			default:
				close(msg.conn.send)
				delete(h.conns, msg.conn)
			}
		}
	}
}
