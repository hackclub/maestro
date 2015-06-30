package baton

import (
	"encoding/json"
	"log"

	"github.com/hackedu/maestro/baton/commands"
)

type rawMsg struct {
	conn conn
	data []byte
}

type hub struct {
	conns      map[conn][]string
	ids        map[string]conn
	register   chan conn
	unregister chan conn
	receive    chan rawMsg
	send       chan commands.Command
}

func (h hub) run() {
	for {
		select {
		case c := <-h.register:
			h.conns[c] = make([]string, 0)
		case c := <-h.unregister:
			if ids, ok := h.conns[c]; ok {
				delete(h.conns, c)
				for _, id := range ids {
					delete(h.ids, id)
				}
				close(c.send)
			}
		case rawMsg := <-h.receive:
			var cmd commands.Command
			if err := json.Unmarshal(rawMsg.data, &cmd); err != nil {
				log.Println("Error unmarshaling message into a command")
				log.Println(err)
				break
			}
			log.Println("cmd:", cmd)
			module, ok := moduleChannels[cmd.Module]
			if !ok {
				log.Println(cmd.Module, "not in modules")
				break
			}
			h.conns[rawMsg.conn] = append(h.conns[rawMsg.conn], cmd.ID)
			h.ids[cmd.ID] = rawMsg.conn
			module <- cmd
		case outMsg := <-h.send:
			c, _ := h.ids[outMsg.ID]
			log.Println(outMsg)
			bytes, err := json.Marshal(outMsg)
			if err != nil {
				log.Println("Error marshaling commands.Command into JSON")
				log.Println(err)
				break
			}
			c.send <- bytes
		}
	}
}

var h = hub{
	conns:      make(map[conn][]string),
	ids:        make(map[string]conn),
	register:   make(chan conn),
	unregister: make(chan conn),
	receive:    make(chan rawMsg),
	send:       make(chan commands.Command),
}

func Run() {
	h.run()
}
