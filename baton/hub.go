package baton

import (
	"encoding/json"
	"log"
)

type rawMsg struct {
	conn conn
	data []byte
}

type CommandID string

type Command struct {
	Module string      `json:"module"`
	Call   string      `json:"call"`
	ID     CommandID   `json:"id"`
	Body   interface{} `json:"body"`
}

type Hub struct {
	conns          map[conn][]CommandID
	ids            map[CommandID]conn
	register       chan conn
	unregister     chan conn
	receive        chan rawMsg
	send           chan Command
	modules        map[string]Module
	moduleChannels map[string]chan<- Command
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			log.Println("Hub: Registering conn", c)
			h.conns[c] = make([]CommandID, 0)
		case c := <-h.unregister:
			log.Println("Hub: Unregistering conn", c)
			if ids, ok := h.conns[c]; ok {
				log.Println("Hub: Associated ids", ids)
				delete(h.conns, c)
				for _, id := range ids {
					delete(h.ids, id)
				}
				close(c.send)
			}
		case rawMsg := <-h.receive:
			var cmd Command
			if err := json.Unmarshal(rawMsg.data, &cmd); err != nil {
				log.Println("Hub: Error unmarshaling message into a command")
				log.Println("Hub:", err)
				break
			}
			log.Println("Hub: Recieved command", cmd.ID)
			log.Println("Hub: Content", cmd)
			module, ok := h.moduleChannels[cmd.Module]
			if !ok {
				log.Println("Hub:", cmd.Module, "not in modules")
				break
			}
			h.conns[rawMsg.conn] = append(h.conns[rawMsg.conn], cmd.ID)
			h.ids[cmd.ID] = rawMsg.conn
			module <- cmd
		case outMsg := <-h.send:
			c, ok := h.ids[outMsg.ID]
			log.Println("Hub: Message for ", outMsg.ID)
			if !ok {
				log.Println("Hub:", outMsg.ID, "is associated with a disconnected client.")
				break
			}
			log.Println(outMsg)
			bytes, err := json.Marshal(outMsg)
			if err != nil {
				log.Println("Hub: Error marshaling Command into JSON")
				log.Println("Hub:", err)
				break
			}
			c.send <- bytes
		}
	}
}

func NewHub() Hub {
	return Hub{
		conns:          make(map[conn][]CommandID),
		ids:            make(map[CommandID]conn),
		register:       make(chan conn),
		unregister:     make(chan conn),
		receive:        make(chan rawMsg),
		send:           make(chan Command),
		modules:        nil,
		moduleChannels: nil,
	}
}
