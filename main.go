package main

import (
	"log"
	"net/http"

	"srcp-rs/handlers"
)

func main() {
	http.HandleFunc("/info", handlers.Info)
	go http.ListenAndServe(":4201", nil)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":4202", router))
}
