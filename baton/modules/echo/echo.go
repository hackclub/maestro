package echo

import "errors"

type Echo struct {
}

func (e Echo) RunCommand(cmd string, body interface{}, resp chan<- interface{}) error {
	if cmd != "echo" {
		return errors.New("unknown command: " + cmd)
	}
	//send resp <- msg
	resp <- body
	return nil
}