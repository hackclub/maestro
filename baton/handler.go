package baton

import(
	"net/http"
	"log"

	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/router"
)

func Handler() *mux.Router {
	m := router.Baton()
	m.Get(router.BatonConnect).HandlerFunc(serveWs)
	return m
}

func Run() {
	h.run()
}

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
	c.readPump()
}
