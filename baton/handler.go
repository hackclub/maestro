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

func WebhookHandler() *mux.Router {
	m := mux.NewRouter()
	for name, module := range modules {
		//PathPrefix needed to make it behave like http.Handle
		m.PathPrefix(fmt.Sprintf("/%s", name)).Handler(http.StripPrefix(fmt.Sprintf("/%s", name), module.Handler()))
	}
	return m
}
