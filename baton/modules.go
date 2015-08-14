package baton

import (
	"github.com/gorilla/mux"
)

type Module interface {
	Init(cmd <-chan Command, resp chan<- Command)
}

type ModuleHandler interface {
	Handler() *mux.Router
}

func (h *Hub) InitModules(modules map[string]Module) {
	h.moduleChannels = make(map[string]chan<- Command)
	for name, module := range modules {
		cmd := make(chan Command, 0)
		h.moduleChannels[name] = cmd
		module.Init(cmd, h.send)
	}
	h.modules = modules
}
