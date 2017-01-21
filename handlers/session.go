package handlers

import (
	"fmt"
	"net/http"
	"srcp-rs/srcp"
	"strconv"
	"strings"
)

func CreateSession(w http.ResponseWriter, r *http.Request) {
	var wrapper Wrapper
	var session Session

	unmarshal(&wrapper, &session, r, w)

	srcpConnection := srcp.NewSrcpConnection()
	srcpConnection.Connect(store.GetSrcpEndpoint())

	reply := srcpConnection.Receive()
	message := srcp.Parse(reply)

	session.Infos = message.ExtractSessionInfos()

	reply = srcpConnection.SendAndReceive(fmt.Sprintf("SET CONNECTIONMODE SRCP %s", strings.ToUpper(session.Mode)))

	if message := srcp.Parse(reply); message.Code != 202 {
		writeReply(SrcpError{message.Code, message.Status, message.Message}, w)
		return
	}

	reply = srcpConnection.SendAndReceive("GO")

	if message := srcp.Parse(reply); message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		session.SessionId = message.ExtractSessionId()
		store.SaveConnection(session.SessionId, srcpConnection)
		writeReply(Wrapper{Data{strconv.Itoa(session.SessionId), "session", session}}, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		writeReply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}
