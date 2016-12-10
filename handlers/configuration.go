package handlers

import (
	"io"
	"log"
	"net/http"
	"gopkg.in/yaml.v2"
)

func Configuration(w http.ResponseWriter, r *http.Request) {
	y, err := yaml.Marshal(store.GetGLS())
	log.Printf("yaml: %s", string(y))
	if err == nil {
		io.WriteString(w, string(y))
	} else {
		log.Println(err)
	}
}
