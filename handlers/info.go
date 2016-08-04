package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	var srcpConnection srcp.SrcpConnection
	var session Session

	srcpConnection.Connect("localhost:4303")

	srcpReply := srcpConnection.Receive()

	session.Infos = srcp.ExtractSessionInfos(srcpReply)

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

func listenAndSend(srcpConnection *srcp.SrcpConnection, websocket *websocket.Conn) {
	defer srcpConnection.Close()
	defer websocket.Close()

	var store map[int]map[int]*srcp.GeneralLoco = make(map[int]map[int]*srcp.GeneralLoco)
	for {
		message := srcp.Parse(srcpConnection.Receive())
		timestamp, error := strconv.ParseFloat(message.Time, 64)
		if error != nil {
			log.Printf("error converting timestamp", error)
			return
		}
		if srcp.ExtractDeviceGroup(message.Message) == "GL" {
			bus, address := srcp.ExtractBusAndAddress(message.Message)
			if bus > -1 && address > -1 {
				if store[bus] == nil {
					store[bus] = make(map[int]*srcp.GeneralLoco)
				}
				var gl *srcp.GeneralLoco = store[bus][address]
				if gl == nil {
					gl = new(srcp.GeneralLoco)
					store[bus][address] = gl
				}
				if timestamp > gl.LastTimestamp {
					srcp.UpdateGeneralLoco(message.Code, message.Message, gl)
					switch message.Code {
					case 100:
						websocket.WriteJSON(InfoMessage{"GL updated", "update", Wrapper{Data{strconv.Itoa(address), "gl", gl}}})
					case 101:
						websocket.WriteJSON(InfoMessage{"GL created", "create", Wrapper{Data{strconv.Itoa(address), "gl", gl}}})
					case 102:
						websocket.WriteJSON(InfoMessage{"GL deleted", "delete", Wrapper{Data{strconv.Itoa(address), "gl", gl}}})
						store[bus][address] = nil
					}
				}
				gl.LastTimestamp = timestamp
			} else {
				log.Printf("bus: %d address: %s", bus, address)
			}
		}
	}
}
