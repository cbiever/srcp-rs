package handlers

import (
	"fmt"
	"net/http"
	"srcp-rs/model"
	"srcp-rs/srcp"
)

func UpdateCV(w http.ResponseWriter, r *http.Request) {
	session, bus, address, cv := extract(r)
	srcpConnection := store.GetConnection(session)

	var cvUpdate model.CV
	unmarshal(nil, &cvUpdate, r, w)

	reply := srcpConnection.SendAndReceive(fmt.Sprintf("INIT %d SM NMRA", bus))
	message := srcp.Parse(reply)

	if message.Code != 200 {
		w.WriteHeader(http.StatusBadRequest)
		writeReply(SrcpError{message.Code, message.Status, message.Message}, w)
		return
	}

	request := fmt.Sprintf("SET %d SM %d %s %d", bus, address, cvUpdate.Type, cv)
	for _, value := range cvUpdate.Value {
		request += fmt.Sprintf(" %d", value)
	}

	reply = srcpConnection.SendAndReceive(request)
	message = srcp.Parse(reply)

	if message.Code != 200 {
		srcpConnection.SendAndReceive(fmt.Sprintf("TERM %d SM", bus))
		w.WriteHeader(http.StatusBadRequest)
		writeReply(SrcpError{message.Code, message.Status, message.Message}, w)
		return
	}

	reply = srcpConnection.SendAndReceive(fmt.Sprintf("TERM %d SM", bus))
	message = srcp.Parse(reply)
	if message.Code == 200 {
		w.WriteHeader(http.StatusOK)
		writeReply(nil, w)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		writeReply(SrcpError{message.Code, message.Status, message.Message}, w)
	}
}
