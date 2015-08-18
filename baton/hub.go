package baton

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
)

var log = logrus.WithField("module", "Hub")

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
			log.Debug("Registering conn", c)
			h.conns[c] = make([]CommandID, 0)
		case c := <-h.unregister:
			log.Debug("Unregistering conn", c)
			if ids, ok := h.conns[c]; ok {
				log.Debug("Associated ids", ids)
				delete(h.conns, c)
				for _, id := range ids {
					delete(h.ids, id)
				}
				close(c.send)
			}
		case rawMsg := <-h.receive:
			var cmd Command
			if err := json.Unmarshal(rawMsg.data, &cmd); err != nil {
				log.WithFields(logrus.Fields{
					"error":  err,
					"rasMsg": rawMsg,
				}).Error("Error unmarshaling message into a command")
				break
			}
			log.WithField("command", cmd).Info("Command recieved")
			module, ok := h.moduleChannels[cmd.Module]
			if !ok {
				log.WithField("command", cmd).Error("Module not found")
				break
			}
			h.conns[rawMsg.conn] = append(h.conns[rawMsg.conn], cmd.ID)
			h.ids[cmd.ID] = rawMsg.conn
			module <- cmd
		case outMsg := <-h.send:
			c, ok := h.ids[outMsg.ID]
			log.WithField("command", outMsg).Debug("Message recieved")
			if !ok {
				log.WithField("command", outMsg).Error("Unassociated ID")
				break
			}
			bytes, err := json.Marshal(outMsg)
			if err != nil {
				log.WithFields(logrus.Fields{
					"error":   err,
					"command": outMsg,
				}).Error("Error unmarshaling message into a command")
				break
			}
			log.WithField("command", outMsg).Info("Message sent")
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
