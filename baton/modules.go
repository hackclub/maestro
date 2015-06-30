package baton

import (
	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/baton/commands"
	"github.com/hackedu/maestro/baton/modules/echo"
	"github.com/hackedu/maestro/baton/modules/giphy"
	"github.com/hackedu/maestro/baton/modules/neutrino"
	"github.com/hackedu/maestro/baton/modules/twilio"
)

var modules = map[string]Module{
	"Echo":  echo.Echo{},
	"Giphy": giphy.Giphy{"dc6zaTOxFJmzC"}, //testing key from Giphy
	"Neutrino": neutrino.Neutrino{"user-id",
		"api-key"},
	"Twilio": twilio.Twilio{"user-id",
		"api-key"},
}

type Module interface {
	Init(cmd <-chan commands.Command, resp chan<- commands.Command)
	Handler() *mux.Router
}

var moduleChannels map[string]chan<- commands.Command

func InitModules() {
	moduleChannels = make(map[string]chan<- commands.Command)
	for name, module := range modules {
		cmd := make(chan commands.Command, 0)
		moduleChannels[name] = cmd
		module.Init(cmd, h.send)
	}
}
