package echo

import (
	"errors"
	"io"
	"log"
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
	m.PathPrefix("/").HandlerFunc(echo)
	return m
}

func echo(w http.ResponseWriter, r *http.Request) {
	_, err := io.Copy(w, r.Body)
	if err != nil {
		log.Println("Echo:", "Error copying request body")
		log.Println("Echo:", err)
	}
}
