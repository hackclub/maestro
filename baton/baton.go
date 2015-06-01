package baton

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	receive:    make(chan rawMsg),
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
		h.receive <- rawMsg{c, data}
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

type command struct {
	Module string      `json:"module"`
	Call   string      `json:"call"`
	Body   interface{} `json:"body"`
}

type rawMsg struct {
	conn conn
	data []byte
}

type msg struct {
	conn conn
	cmd  command
}

type hub struct {
	conns      map[conn]struct{}
	register   chan conn
	unregister chan conn
	receive    chan rawMsg
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
		case rawMsg := <-h.receive:
			fmt.Println(rawMsg)
			var cmd command
			if err := json.Unmarshal(rawMsg.data, &cmd); err != nil {
				fmt.Println("nu")
				log.Println(err)
				break
			}
			
			msg := msg{rawMsg.conn, cmd}

			fmt.Println(cmd)
			resp := make(chan interface{})
			module, ok := modules[msg.cmd.Module]
			if !ok {
				fmt.Println("the fuck is this module?")
				break
			}

			go func() {
				for {
					r, ok := <-resp
					if !ok {
						close(resp)
					}
					bytes, err := json.Marshal(command{cmd.Module,cmd.Call,r}) //add in Module and Call info for client
					if err != nil {
					  log.Println("Error Marshaling")
						log.Println(err)
						break
					}
        
					msg.conn.send <- bytes
				}
			}()

			if err := module.RunCommand(cmd.Call, cmd.Body, resp); err != nil {
			  fmt.Println(err)
				fmt.Println("the fuck happened to this command?")
				break
			}
		}
	}
}
