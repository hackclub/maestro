package baton

import (
	"encoding/json"
	"fmt"
	"log"
)

type command struct {
	Module string      `json:"module"`
	Call   string      `json:"call"`
	ID     string      `json:"id"`
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
			var cmd command
			if err := json.Unmarshal(rawMsg.data, &cmd); err != nil {
				fmt.Println("nu")
				log.Println(err)
				break
			}
			fmt.Println(cmd)
			msg := msg{rawMsg.conn, cmd}
			processMessage(msg)
		}
	}
}

func processMessage(message msg) {
	resp := make(chan interface{})
	module, ok := modules[message.cmd.Module]
	if !ok {
		fmt.Println(message.cmd.Module, "not in", modules)
		return
	}

	go func() {
		for {
			r, ok := <-resp
			if !ok {
				close(resp)
				break
			}
			bytes, err := json.Marshal(command{message.cmd.Module, message.cmd.Call, message.cmd.ID, r}) //add in Module and Call info for client
			if err != nil {
				log.Println("Error Marshaling")
				log.Println(err)
				break
			}
			message.conn.send <- bytes
		}
	}()
	go func() {
		if err := module.RunCommand(message.cmd.Call, message.cmd.Body, resp); err != nil {
			fmt.Println(err)
		}
	}()
}

var h = hub{
	conns:      make(map[conn]struct{}),
	register:   make(chan conn),
	unregister: make(chan conn),
	receive:    make(chan rawMsg),
}

func Run() {
	h.run()
}
