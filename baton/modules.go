package baton

import (
  "github.com/hackedu/maestro/baton/echo"
  "github.com/hackedu/maestro/baton/giphy"
)
var modules = map[string]Module{
	"Echo": echo.Echo{},
	"Giphy": giphy.Giphy{},
}

type Module interface {
	RunCommand(cmd string, body interface{}, resp chan<- interface{}) error
}
