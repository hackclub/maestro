package echo

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Echo struct {
}

func (e Echo) RunCommand(cmd string, body interface{}, resp chan<- interface{}) error {
	if cmd != "echo" {
		return errors.New("unknown command: " + cmd)
	}
	resp <- body
	return nil
}

func (e Echo) Handler() *mux.Router {
	m := mux.NewRouter()
	m.Path("*").HandlerFunc(echo)
	return m
}
func echo(w http.ResponseWriter, r *http.Request) {
	_, err := io.Copy(w, r.Body)
	if err != nil {
		fmt.Println(err)
	}
}
