package handlers

import (
	"fmt"
	"net/http"
	"srcp-rs/model"
	"srcp-rs/srcp"
	"strconv"
)

func GetBuses(w http.ResponseWriter, r *http.Request) {
	session, _, _ := extract(r)
	srcpConnection := store.GetConnection(session)
	var buses []Data
	var bus = 1
	for {
		reply := srcpConnection.SendAndReceive(fmt.Sprintf("GET %d DESCRIPTION", bus))
		message := srcp.Parse(reply)
		if message.Code == 100 {
			buses = append(buses, Data{strconv.Itoa(bus), "bus", model.Bus{message.ExtractDeviceGroups()}})
		} else {
			w.WriteHeader(http.StatusOK)
			writeReply(ArrayWrapper{buses}, w)
			return
		}
		bus++
	}
}

func DeleteBus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeReply(nil, w)
}
