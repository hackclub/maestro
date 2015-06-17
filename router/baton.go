package router

import "github.com/gorilla/mux"

func Baton() *mux.Router {
	m := mux.NewRouter()
	m.Path("/connect").Methods("GET").Name(BatonConnect)
	return m
}

func App() *mux.Router {
	m := mux.NewRouter()
	m.PathPrefix("/static/").Name(AppStatic)
	return m
}
