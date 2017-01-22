package test

import (
	"github.com/gorilla/mux"
	"net/http"
	"srcp-rs/handlers"
)

var router *mux.Router = nil

func createRouter() {
	if router == nil {
		router = mux.NewRouter()
		router.HandleFunc("/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls", handlers.CreateGL)
		router.HandleFunc("/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}/cv/{cv:[0-9]+}", handlers.UpdateCV)
		http.Handle("/", router)
	}
}
