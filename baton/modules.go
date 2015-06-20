package baton

import (
	"github.com/gorilla/mux"
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
	RunCommand(cmd string, body interface{}, resp chan<- interface{}) error
	Handler() *mux.Router
}
