package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"srcp-rs/srcp"
)

func CreateGL(w http.ResponseWriter, r *http.Request) {
	_, bus, _ := extract(r)

	var wrapper Wrapper
	var gl srcp.GeneralLoco
	unmarshal(&wrapper, &gl, r, w)

	srcpReply := srcpConnection.SendAndReceive(fmt.Sprintf("INIT %d GL %d %s %d %d %d", bus, gl.Address, gl.Protocol, gl.ProtocolVersion, gl.DecoderSpeedSteps, gl.NumberOfDecoderFunctions))

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		reply(Wrapper{Data{strconv.Itoa(gl.Address), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func GetGL(w http.ResponseWriter, r *http.Request) {
	_, bus, address := extract(r)

	srcpReply1 := srcpConnection.SendAndReceive(fmt.Sprintf("GET %d DESCRIPTION GL %d", bus, address))
	message1 := srcp.Parse(srcpReply1)

	srcpReply2 := srcpConnection.SendAndReceive(fmt.Sprintf("GET %d GL %d", bus, address))
	message2 := srcp.Parse(srcpReply2)

	if message1.Code == 101 && message2.Code == 100 {
		w.WriteHeader(http.StatusOK)
		var gl srcp.GeneralLoco
		srcp.UpdateGeneralLoco(message1.Code, message1.Message, &gl)
		srcp.UpdateGeneralLoco(message2.Code, message2.Message, &gl)
		reply(Wrapper{Data{strconv.Itoa(address), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		if message1.Code != 101 {
			reply(SrcpError{message1.Code, message1.Status, message1.Message}, w)
		} else {
			reply(SrcpError{message2.Code, message2.Status, message2.Message}, w)
		}
	}
}

func UpdateGL(w http.ResponseWriter, r *http.Request) {
	_, bus, address := extract(r)

	var wrapper Wrapper
	var gl srcp.GeneralLoco
	unmarshal(&wrapper, &gl, r, w)

	request := fmt.Sprintf("SET %d GL %d %d %d %d", bus, address, gl.Drivemode, gl.V, gl.Vmax)
	for _, function := range gl.Function {
		request += fmt.Sprintf(" %d", function)
	}
	srcpReply := srcpConnection.SendAndReceive(request)

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		gl.Bus = bus
		reply(Wrapper{Data{strconv.Itoa(address), "gl", gl}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}

func DeleteGL(w http.ResponseWriter, r *http.Request) {
	_, bus, address := extract(r)

	srcpReply := srcpConnection.SendAndReceive(fmt.Sprintf("TERM %d GL %d", bus, address))

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}
