package baton

import (
	"encoding/json"
	"fmt"
	"log"
)

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
				fmt.Println(msg.cmd.Module, "not in", modules)
				break
			}

			go func() {
				for {
					r, ok := <-resp
					if !ok {
						close(resp)
					}
					bytes, err := json.Marshal(command{cmd.Module, cmd.Call, r}) //add in Module and Call info for client
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
