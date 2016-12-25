package main

import (
	"log"
	"net/http"
	"os"

	"srcp-rs/handlers"
)

func main() {
	if len(os.Args) > 1 {
		handlers.GetStore().SetSrcpEndpoint(os.Args[1])
	} else {
		handlers.GetStore().SetSrcpEndpoint("localhost:4303")
	}

	http.HandleFunc("/info", handlers.Info)
	go http.ListenAndServe(":4201", nil)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":4202", router))
}
