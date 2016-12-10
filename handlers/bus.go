package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"srcp-rs/srcp"
)

func GetBuses(w http.ResponseWriter, r *http.Request) {
	session, _, _ := extract(r)
	srcpConnection := store.GetConnection(session)
	var buses []Data
	var bus = 1
	for {
		srcpReply := srcpConnection.SendAndReceive(fmt.Sprintf("GET %d DESCRIPTION", bus))
		message := srcp.Parse(srcpReply)
		if message.Code == 100 {
			buses = append(buses, Data{strconv.Itoa(bus), "bus", srcp.Bus{srcp.ExtractDeviceGroups(message.Message)}})
		} else {
			w.WriteHeader(http.StatusOK)
			reply(ArrayWrapper{buses}, w)
			return
		}
		bus++
	}
}

func DeleteBus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	reply(nil, w)
}
