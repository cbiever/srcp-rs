package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func extract(r *http.Request) (int, int, int, int) {
	vars := mux.Vars(r)
	sessionId, _ := strconv.Atoi(vars["sessionId"])
	bus, _ := strconv.Atoi(vars["bus"])
	address, _ := strconv.Atoi(vars["address"])
	cv, _ := strconv.Atoi(vars["cv"])
	return sessionId, bus, address, cv
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
	var j []byte
	if wrapper != nil {
		if error = json.Unmarshal(body, wrapper); error != nil {
			if error = json.NewEncoder(writer).Encode(error); error != nil {
				panic(error)
			}
		}
		if j, error = json.Marshal(wrapper.Attributes); error != nil {
			panic(error)
		}
	} else {
		j = body
	}
	if error = json.Unmarshal(j, &payload); error != nil {
		panic(error)
	}
}

func writeReply(response interface{}, writer http.ResponseWriter) {
	if error := json.NewEncoder(writer).Encode(response); error != nil {
		panic(error)
	}
}
