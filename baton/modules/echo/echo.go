package echo

import (
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/baton"
)

var log = logrus.WithField("module", "Echo")

type Echo struct {
}

func (e Echo) Init(cmd <-chan baton.Command, resp chan<- baton.Command) {
	resp = resp
	go func() {
		for {
			tmp := <-cmd
			if tmp.Call != "echo" {
				log.WithField("command", tmp).Error("Unknown command")
				continue
			}
			log.WithField("command", tmp).Debug("Processing command")
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
	log.Debug("Recieved Message over HTTP")
	_, err := io.Copy(w, r.Body)
	if err != nil {
		log.WithField("error", err).Error("Error copying request body")
	}
}
