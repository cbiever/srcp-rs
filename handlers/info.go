package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"srcp-rs/srcp"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}

func Info(w http.ResponseWriter, r *http.Request) {
	websocket, error := upgrader.Upgrade(w, r, nil)
	if error != nil {
		panic(error)
	}

	var srcpConnection SrcpConnection
	var session Session

	srcpConnection.Connect("localhost:4303")

	srcpReply := srcpConnection.Receive()

	session.Infos = make(map[string]string)
	for _, info := range strings.Split(srcpReply, ";") {
		keyValue := strings.Split(strings.Trim(info, " "), " ")
		session.Infos[keyValue[0]] = keyValue[1]
	}

	srcpReply = srcpConnection.SendAndReceive("SET CONNECTIONMODE SRCP INFORMATION")

	if message := srcp.Parse(srcpReply); message.Code != 202 {
		websocket.WriteJSON("")
		websocket.Close()
		return
	}

	srcpReply = srcpConnection.SendAndReceive("GO")

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		session.SessionId = srcp.ExtractSessionId(message.Message)
		websocket.WriteJSON(InfoMessage{"Info session created", "created", Wrapper{Data{fmt.Sprintf("%d", session.SessionId), "session", session}}})
		listenAndSend(&srcpConnection, websocket)
	} else {
		websocket.WriteJSON("")
		websocket.Close()
	}
}

func listenAndSend(srcpConnection *SrcpConnection, websocket *websocket.Conn) {
	defer srcpConnection.Close()
	defer websocket.Close()

	var lastTimestamp float64 = 0
	var store map[int]map[int]*srcp.GeneralLoco = make(map[int]map[int]*srcp.GeneralLoco)
	for {
		message := srcp.Parse(srcpConnection.Receive())
		timestamp, error := strconv.ParseFloat(message.Time, 64)
		if error != nil {
			log.Printf("error converting timestamp", error)
			return
		}
		if timestamp > lastTimestamp {
			if srcp.ExtractDeviceGroup(message.Message) == "GL" {
				var gl *srcp.GeneralLoco = new(srcp.GeneralLoco)
				srcp.UpdateGeneralLoco(message.Message, gl)
				if store[gl.Bus] == nil {
					store[gl.Bus] = make(map[int]*srcp.GeneralLoco)
				}
				if store[gl.Bus][gl.Address] == nil {
					store[gl.Bus][gl.Address] = gl
				} else {
					gl = store[gl.Bus][gl.Address]
					srcp.UpdateGeneralLoco(message.Message, gl)
				}
				id := fmt.Sprintf("%d-%d", gl.Bus, gl.Address)
				switch message.Code {
				case 100:
					websocket.WriteJSON(InfoMessage{"GL updated", "update", Wrapper{Data{id, "gl", gl}}})
				case 101:
					websocket.WriteJSON(InfoMessage{"GL created", "created", Wrapper{Data{id, "gl", gl}}})
				case 102:
					websocket.WriteJSON(InfoMessage{"GL deleted", "deleted", Wrapper{Data{id, "gl", gl}}})
				}
			}
		}
		lastTimestamp = timestamp
	}
}
