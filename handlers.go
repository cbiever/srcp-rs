package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"./srcp"

	"github.com/gorilla/mux"
)

var connection net.Conn

func SetHeader(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		inner.ServeHTTP(w, r)
	})
}

func CreateSession(w http.ResponseWriter, r *http.Request) {
	var error error

	var request SessionRequest
	unmarshall(&request, r, w)

	connection, error = net.Dial("tcp", "localhost:4303")
	if error != nil {
		panic(error)
	}

	srcpReply := send("")

	var response SessionResponse
	response.Infos = make(map[string]string)
	for _, info := range strings.Split(srcpReply, ";") {
		keyValue := strings.Split(strings.Trim(info, " "), " ")
		response.Infos[keyValue[0]] = keyValue[1]
	}

	srcpReply = send(fmt.Sprintf("SET CONNECTIONMODE SRCP %s", strings.ToUpper(request.Mode)))

	if message := srcp.Parse(srcpReply); message.Code != 202 {
		if error = json.NewEncoder(w).Encode(response); error != nil {
			panic(error)
		}
		return
	}

	srcpReply = send("GO")

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		response.SessionId = srcp.ExtractSessionId(message.Message)
		reply(response, w)
	}
}

func CreateGL(w http.ResponseWriter, r *http.Request) {
	//	vars := mux.Vars(r)

	//	sessionId := vars["sessionId"]

	var request GeneralLocoCreateRequest
	unmarshall(&request, r, w)

	srcpReply := send(fmt.Sprintf("INIT %d GL %d %s %d %d %d", request.Bus, request.Address, request.Protocol, request.ProtocalVersion, request.DecoderSpeedSteps, request.NumberOfDecoderFunctions))

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		reply(GeneralLocoCreateResponse{message.Time}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func GetGL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//sessionId := vars["sessionId"]
	bus := vars["bus"]
	address := vars["address"]

	srcpReply := send(fmt.Sprintf("GET %s GL %s", bus, address))

	message := srcp.Parse(srcpReply)
	if message.Code == 100 {
		w.WriteHeader(http.StatusOK)
		glValues := srcp.ExtractGLValues(message.Message)
		reply(GeneralLocoGetResponse{message.Time, glValues.Drivemode, glValues.V, glValues.Vmax, glValues.Function}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func UpdateGL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//	sessionId := vars["sessionId"]
	bus := vars["bus"]
	address := vars["address"]

	var updateRequest GeneralLocoUpdateRequest
	unmarshall(&updateRequest, r, w)

	request := fmt.Sprintf("SET %s GL %s %d %d %d", bus, address, updateRequest.Drivemode, updateRequest.V, updateRequest.Vmax)
	for _, function := range updateRequest.Function {
		request += fmt.Sprintf(" %d", function)
	}
	srcpReply := send(request)

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		reply(GeneralLocoUpdateResponse{message.Time}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func DeleteGL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//	sessionId := vars["sessionId"]
	bus := vars["bus"]
	address := vars["address"]

	srcpReply := send(fmt.Sprintf("TERM %s GL %s", bus, address))

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		reply(GeneralLocoUpdateResponse{message.Time}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func unmarshall(request interface{}, reader *http.Request, writer http.ResponseWriter) {
	var body []byte
	var error error
	if body, error = ioutil.ReadAll(io.LimitReader(reader.Body, 1048576)); error != nil {
		panic(error)
	}
	if error = reader.Body.Close(); error != nil {
		panic(error)
	}
	if error = json.Unmarshal(body, &request); error != nil {
		if error = json.NewEncoder(writer).Encode(error); error != nil {
			panic(error)
		}
	}
}

func send(request string) string {
	var error error
	if request != "" {
		log.Printf(request)
		if _, error = connection.Write([]byte(request + "\n")); error != nil {
			panic(error)
		}
	}
	var reply string
	if reply, error = bufio.NewReader(connection).ReadString('\n'); error != nil {
		panic(error)
	}
	log.Printf(reply)
	return reply
}

func reply(response interface{}, writer http.ResponseWriter) {
	if error := json.NewEncoder(writer).Encode(response); error != nil {
		panic(error)
	}
}
