package baton

import (
  "github.com/hackedu/maestro/baton/echo"
  "github.com/hackedu/maestro/baton/giphy"
  "github.com/hackedu/maestro/baton/neutrino"
)
var modules = map[string]Module{
	"Echo": echo.Echo{},
	"Giphy": giphy.Giphy{"dc6zaTOxFJmzC"},
	"Neutrino": neutrino.Neutrino{"user-id","api-key"},
}

type Module interface {
	RunCommand(cmd string, body interface{}, resp chan<- interface{}) error
}
