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
	"strconv"
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
	var wrapper Wrapper
	var session Session

	unmarshal(&wrapper, &session, r, w)

	connection, error = net.Dial("tcp", "localhost:4303")
	if error != nil {
		panic(error)
	}

	srcpReply := send("")

	session.Infos = make(map[string]string)
	for _, info := range strings.Split(srcpReply, ";") {
		keyValue := strings.Split(strings.Trim(info, " "), " ")
		session.Infos[keyValue[0]] = keyValue[1]
	}

	srcpReply = send(fmt.Sprintf("SET CONNECTIONMODE SRCP %s", strings.ToUpper(session.Mode)))

	if message := srcp.Parse(srcpReply); message.Code != 202 {
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
		return
	}

	srcpReply = send("GO")

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		session.SessionId = srcp.ExtractSessionId(message.Message)
		reply(Wrapper{Data{strconv.Itoa(session.SessionId), "session", session}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func CreateGL(w http.ResponseWriter, r *http.Request) {
	sessionId, bus, _ := extract(r)

	var wrapper Wrapper
	var gl GeneralLoco
	unmarshal(&wrapper, &gl, r, w)

	srcpReply := send(fmt.Sprintf("INIT %d GL %d %s %d %d %d", bus, gl.Address, gl.Protocol, gl.ProtocalVersion, gl.DecoderSpeedSteps, gl.NumberOfDecoderFunctions))

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		reply(Wrapper{Data{fmt.Sprintf("%d-%d", sessionId, bus), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func GetGL(w http.ResponseWriter, r *http.Request) {
	sessionId, bus, address := extract(r)

	srcpReply := send(fmt.Sprintf("GET %d GL %d", bus, address))

	message := srcp.Parse(srcpReply)
	if message.Code == 100 {
		w.WriteHeader(http.StatusOK)
		values := srcp.ExtractGLValues(message.Message)
		gl := GeneralLoco{address, "N", 1, 128, 4, values.Drivemode, values.V, values.Vmax, values.Function}
		reply(Wrapper{Data{fmt.Sprintf("%d-%d", sessionId, bus), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func UpdateGL(w http.ResponseWriter, r *http.Request) {
	sessionId, bus, address := extract(r)

	var wrapper Wrapper
	var gl GeneralLoco
	unmarshal(&wrapper, &gl, r, w)

	request := fmt.Sprintf("SET %d GL %d %d %d %d", bus, address, gl.Drivemode, gl.V, gl.Vmax)
	for _, function := range gl.Function {
		request += fmt.Sprintf(" %d", function)
	}
	srcpReply := send(request)

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		reply(Wrapper{Data{fmt.Sprintf("%d-%d", sessionId, bus), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func DeleteGL(w http.ResponseWriter, r *http.Request) {
	_, bus, address := extract(r)

	srcpReply := send(fmt.Sprintf("TERM %d GL %d", bus, address))

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

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
