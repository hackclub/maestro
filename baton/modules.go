package baton

import "github.com/hackedu/maestro/baton/echo"

var modules = map[string]Module{
	"echo": echo.Echo{},
}

type Module interface {
	RunCommand(cmd string, body interface{}, resp chan<- interface{}) error
}
