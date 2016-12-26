package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"srcp-rs/model"
	"srcp-rs/srcp"
	"strings"
)

func CreateGL(w http.ResponseWriter, r *http.Request) {
	session, bus, _ := extract(r)
	srcpConnection := store.GetConnection(session)

	var wrapper Wrapper
	var gl model.GeneralLoco
	unmarshal(&wrapper, &gl, r, w)

	reply := srcpConnection.SendAndReceive(fmt.Sprintf("INIT %d GL %d %s %d %d %d", bus, gl.Address, gl.Protocol, gl.ProtocolVersion, gl.DecoderSpeedSteps, gl.NumberOfDecoderFunctions))

	message := srcp.Parse(reply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		writeReply(Wrapper{Data{fmt.Sprintf("%d-%d", bus, gl.Address), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		writeReply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func GetGL(w http.ResponseWriter, r *http.Request) {
	var err error
	session, bus, address := extract(r)
	srcpConnection := store.GetConnection(session)

	reply1 := srcpConnection.SendAndReceive(fmt.Sprintf("GET %d DESCRIPTION GL %d", bus, address))
	message1 := srcp.Parse(reply1)

	reply2 := srcpConnection.SendAndReceive(fmt.Sprintf("GET %d GL %d", bus, address))
	message2 := srcp.Parse(reply2)

	if message1.Code == 101 && message2.Code == 100 {
		w.WriteHeader(http.StatusOK)
		var gl model.GeneralLoco
		if gl.Bus, gl.Address, gl.Drivemode, gl.V, gl.Vmax, gl.Function, err = message1.ExtractGLDescriptionValues(); err != nil {
			panic(err)
		}
		if gl.Bus, gl.Address, gl.Protocol, gl.ProtocolVersion, gl.DecoderSpeedSteps, gl.NumberOfDecoderFunctions, err = message2.ExtractGLInitValues(); err != nil {
			panic(err)
		}
		writeReply(Wrapper{Data{fmt.Sprintf("%d-%d", bus, gl.Address), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		if message1.Code != 101 {
			writeReply(SrcpError{message1.Code, message1.Status, message1.Message}, w)
		} else {
			writeReply(SrcpError{message2.Code, message2.Status, message2.Message}, w)
		}
	}
}

func UpdateGL(w http.ResponseWriter, r *http.Request) {
	session, bus, address := extract(r)
	srcpConnection := store.GetConnection(session)

	var wrapper Wrapper
	var gl model.GeneralLoco
	unmarshal(&wrapper, &gl, r, w)

	request := fmt.Sprintf("SET %d GL %d %d %d %d", bus, address, gl.Drivemode, gl.V, gl.Vmax)
	for _, function := range gl.Function {
		request += fmt.Sprintf(" %d", function)
	}
	reply := srcpConnection.SendAndReceive(request)

	message := srcp.Parse(reply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		gl.Bus = bus
		writeReply(Wrapper{Data{fmt.Sprintf("%d-%d", bus, gl.Address), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		writeReply(SrcpError{message.Code, message.Status, message.Message}, w)
	}

	if existingGL := store.GetGL(bus, address); existingGL != nil && strings.Compare(existingGL.Name, gl.Name) != 0 {
		existingGL.Name = gl.Name
		data, error := json.Marshal(Data{fmt.Sprintf("%d-%d", bus, gl.Address), "gl", gl})
		if error != nil {
			panic(error)
		}
		srcpConnection.SendAndReceive(fmt.Sprintf("SET 0 GM 0 0 TEXT %s", string(data)))
	}
}

func DeleteGL(w http.ResponseWriter, r *http.Request) {
	session, bus, address := extract(r)
	srcpConnection := store.GetConnection(session)

	reply := srcpConnection.SendAndReceive(fmt.Sprintf("TERM %d GL %d", bus, address))

	message := srcp.Parse(reply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		writeReply(nil, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		writeReply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}
