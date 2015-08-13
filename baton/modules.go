package baton

import (
	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/baton/commands"
)

type Module interface {
	Init(cmd <-chan commands.Command, resp chan<- commands.Command)
}

type ModuleHandler interface {
	Handler() *mux.Router
}

func (h *Hub) InitModules(modules map[string]Module) {
	h.moduleChannels = make(map[string]chan<- commands.Command)
	for name, module := range modules {
		cmd := make(chan commands.Command, 0)
		h.moduleChannels[name] = cmd
		module.Init(cmd, h.send)
	}
	h.modules = modules
}
