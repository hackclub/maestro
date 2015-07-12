package baton

import (
	"testing"

	"github.com/hackedu/maestro/baton/commands"
)

func TestHub(t *testing.T) {
	h := hub{
		conns:      make(map[conn][]string),
		ids:        make(map[string]conn),
		register:   make(chan conn),
		unregister: make(chan conn),
		receive:    make(chan rawMsg),
		send:       make(chan commands.Command),
		modules:    make(map[string]chan<- commands.Command),
	}

	module := make(chan commands.Command)
	h.modules["module"] = module

	go h.run()

	user := conn{nil, make(chan []byte)}
	h.register <- user

	data := []byte(`{"module":"module","call":"testCall","id":"0-0","body":{"data":"data"}}`)
	go func() {
		h.receive <- rawMsg{user, data}
	}()
	command := <-module
	if command.ID != "0-0" || command.Module != "module" || command.Call != "testCall" {
		t.Error("Command metadata is incorrect")
	}

	go func() {
		h.send <- command
	}()
	response := <-user.send
	if string(response) != string(data) {
		t.Error("Echoed data is not equal")
	}

	h.unregister <- user
	user = conn{nil, make(chan []byte)}
	h.register <- user

	data = []byte(`{"module":"module","call":"testCall","id":"1-0","body":{"data":"data"}}`)
	go func() {
		h.receive <- rawMsg{user, data}
	}()
	command2 := <-module
	go func() {
		h.send <- command
	}()
	go func() {
		h.send <- command2
	}()
	response = <-user.send
	if string(response) != string(data) {
		t.Error("Echoed data is not equal")
	}
}
