package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"srcp-rs/srcp"
)

var srcpConnection srcp.SrcpConnection

func CreateSession(w http.ResponseWriter, r *http.Request) {
	var wrapper Wrapper
	var session Session

	unmarshal(&wrapper, &session, r, w)

	srcpConnection.Connect("localhost:4303")

	srcpReply := srcpConnection.Receive()

	session.Infos = srcp.ExtractSessionInfos(srcpReply)

	srcpReply = srcpConnection.SendAndReceive(fmt.Sprintf("SET CONNECTIONMODE SRCP %s", strings.ToUpper(session.Mode)))

	if message := srcp.Parse(srcpReply); message.Code != 202 {
		reply(SrcpError{message.Code, message.Status, message.Message}, w)
		return
	}

	srcpReply = srcpConnection.SendAndReceive("GO")

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
