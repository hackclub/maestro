package echo

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/baton/commands"
)

type Echo struct {
}

var resp chan<- commands.Command

func (e Echo) Init(cmd <-chan commands.Command, resp chan<- commands.Command) {
	resp = resp
	go func() {
		for {
			tmp := <-cmd
			if tmp.Call != "echo" {
				log.Println("unknown command: " + tmp.Call)
			}
			resp <- tmp
		}
	}()
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
