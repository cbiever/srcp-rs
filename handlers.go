package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"./srcp"

	"github.com/gorilla/mux"
)

var connection net.Conn

func CreateSession(w http.ResponseWriter, r *http.Request) {
	var error error

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, error := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if error != nil {
		panic(error)
	}
	if error = r.Body.Close(); error != nil {
		panic(error)
	}

	var request SessionRequest
	if error = json.Unmarshal(body, &request); error != nil {
		if error := json.NewEncoder(w).Encode(error); error != nil {
			panic(error)
		}
	}

	connection, error = net.Dial("tcp", "localhost:4303")
	if error != nil {
		panic(error)
	}
	reader := bufio.NewReader(connection)
	var reply string
	reply, error = reader.ReadString('\n')

	var response SessionResponse
	response.Infos = make(map[string]string)
	for _, info := range strings.Split(reply, ";") {
		keyValue := strings.Split(strings.Trim(info, " "), " ")
		response.Infos[keyValue[0]] = keyValue[1]
	}

	if _, error = connection.Write([]byte("SET CONNECTIONMODE SRCP " + strings.ToUpper(request.Mode) + "\n")); error != nil {
		panic(error)
	}
	if reply, error = reader.ReadString('\n'); error != nil {
		panic(error)
	}
	if message := srcp.Parse(reply); message.Code != 202 {
		if error = json.NewEncoder(w).Encode(response); error != nil {
			panic(error)
		}
		return
	}

	if _, error = connection.Write([]byte("GO\n")); error != nil {
		panic(error)
	}
	if reply, error = reader.ReadString('\n'); error != nil {
		panic(error)
	}

	message := srcp.Parse(reply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		response.SessionId = srcp.ExtractSessionId(message.Message)
		if error = json.NewEncoder(w).Encode(response); error != nil {
			panic(error)
		}
	}
}

func GetGL(w http.ResponseWriter, r *http.Request) {
	var error error

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	sessionId := vars["sessionId"]
	bus := vars["bus"]
	address := vars["address"]

	log.Printf("sessionId: %s bus: %s address: %s", sessionId, bus, address)

	if _, error = connection.Write([]byte("GET " + bus + " GL " + address + "\n")); error != nil {
		panic(error)
	}
	var reply string
	if reply, error = bufio.NewReader(connection).ReadString('\n'); error != nil {
		panic(error)
	}

	message := srcp.Parse(reply)
	if message.Code == 100 {
		w.WriteHeader(http.StatusOK)
		var response GeneralLocoGetResponse
		response.Time = message.Time
		if error = json.NewEncoder(w).Encode(response); error != nil {
			panic(error)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		var errorMessage jsonErr
		errorMessage.Code = message.Code
		errorMessage.Text = message.Message
		if error = json.NewEncoder(w).Encode(errorMessage); error != nil {
			panic(error)
		}
	}
}

func CreateGL(w http.ResponseWriter, r *http.Request) {
	var error error

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	sessionId := vars["sessionId"]

	log.Printf("sessionId: %s", sessionId)

	var body []byte
	if body, error = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); error != nil {
		panic(error)
	}
	if error = r.Body.Close(); error != nil {
		panic(error)
	}

	var request GeneralLocoCreateRequest
	if error = json.Unmarshal(body, &request); error != nil {
		if error = json.NewEncoder(w).Encode(error); error != nil {
			panic(error)
		}
	}

	if _, error = connection.Write([]byte("INIT " + strconv.Itoa(request.Bus) + " GL " + strconv.Itoa(request.Address) + " " + strings.ToUpper(request.Protocol) + " " + strconv.Itoa(request.ProtocalVersion) + " " + strconv.Itoa(request.DecoderSpeedSteps) + " " + strconv.Itoa(request.NumberOfDecoderFunctions) + "\n")); error != nil {
		panic(error)
	}
	var reply string
	if reply, error = bufio.NewReader(connection).ReadString('\n'); error != nil {
		panic(error)
	}

	message := srcp.Parse(reply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		var response GeneralLocoCreateResponse
		response.Time = message.Time
		if error = json.NewEncoder(w).Encode(response); error != nil {
			panic(error)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		var errorMessage jsonErr
		errorMessage.Code = message.Code
		errorMessage.Text = message.Message
		if error = json.NewEncoder(w).Encode(errorMessage); error != nil {
			panic(error)
		}
	}
}
