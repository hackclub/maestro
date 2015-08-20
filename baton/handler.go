package baton

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/router"
)

func (h Hub) Handler() *mux.Router {
	m := router.Baton()
	m.Get(router.BatonConnect).HandlerFunc(h.serveWs)
	for name, module := range h.modules {
		//PathPrefix needed to make it behave like http.Handle
		// /webhooks must be included because their documentation lies
		if module, ok := module.(ModuleHandler); ok {
			m.PathPrefix(fmt.Sprintf("/webhooks/%s", name)).Handler(http.StripPrefix(fmt.Sprintf("/webhooks/%s", name), module.Handler()))
		}
	}
	return m
}

var i = 0

func (h Hub) serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithField("error", err).Error("Couldn't upgrade connection")
		return
	}
	c := conn{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.send <- []byte(strconv.Itoa(i))
	i++
	c.readPump(h)
}
