package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func extract(r *http.Request) (int, int, int) {
	vars := mux.Vars(r)
	sessionId, _ := strconv.Atoi(vars["sessionId"])
	bus, _ := strconv.Atoi(vars["bus"])
	address, _ := strconv.Atoi(vars["address"])
	return sessionId, bus, address
}

func unmarshal(wrapper *Wrapper, payload interface{}, reader *http.Request, writer http.ResponseWriter) {
	var body []byte
	var error error
	if body, error = ioutil.ReadAll(io.LimitReader(reader.Body, 1048576)); error != nil {
		panic(error)
	}
	if error = reader.Body.Close(); error != nil {
		panic(error)
	}
	if error = json.Unmarshal(body, wrapper); error != nil {
		if error = json.NewEncoder(writer).Encode(error); error != nil {
			panic(error)
		}
	}
	j, error := json.Marshal(wrapper.Attributes)
	if error != nil {
		panic(error)
	}
	error = json.Unmarshal(j, &payload)
	if error != nil {
		panic(error)
	}
}

func writeReply(response interface{}, writer http.ResponseWriter) {
	if error := json.NewEncoder(writer).Encode(response); error != nil {
		panic(error)
	}
}
