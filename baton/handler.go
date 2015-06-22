package baton

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/router"
)

func Handler() *mux.Router {
	m := router.Baton()
	m.Get(router.BatonConnect).HandlerFunc(serveWs)
	for name, module := range modules {
		//PathPrefix needed to make it behave like http.Handle
		// /webhooks must be included because their documentation lies
		m.PathPrefix(fmt.Sprintf("/webhooks/%s", name)).Handler(http.StripPrefix(fmt.Sprintf("/webhooks/%s", name), module.Handler()))
	}
	return m
}

var i = 0

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := conn{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.send <- []byte(strconv.Itoa(i))
	i++
	c.readPump(h)
}
