package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"srcp-rs/srcp"
	"strconv"
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

	srcpConnection.Connect(store.GetSrcpEndpoint())

	srcpReply := srcpConnection.Receive()

	session.Infos = srcp.ExtractSessionInfos(srcpReply)

	srcpReply = srcpConnection.SendAndReceive("SET CONNECTIONMODE SRCP INFORMATION")

	if message := srcp.Parse(srcpReply); message.Code != 202 {
		websocket.WriteJSON(SrcpError{message.Code, message.Status, message.Message})
		websocket.Close()
		return
	}

	srcpReply = srcpConnection.SendAndReceive("GO")

	message := srcp.Parse(srcpReply)
	if message.Code == 200 {
		session.SessionId = srcp.ExtractSessionId(message.Message)
		websocket.WriteJSON(InfoMessage{"Info session created", "created", Data{fmt.Sprintf("%d", session.SessionId), "session", session}})
		listenAndSend(&srcpConnection, websocket)
	} else {
		websocket.WriteJSON(SrcpError{message.Code, message.Status, message.Message})
		websocket.Close()
	}
}

func listenAndSend(srcpConnection *srcp.SrcpConnection, websocket *websocket.Conn) {
	defer srcpConnection.Close()
	defer websocket.Close()

	for {
		message := srcp.Parse(srcpConnection.Receive())
		timestamp, error := strconv.ParseFloat(message.Time, 64)
		if error != nil {
			log.Printf("error converting timestamp", error)
			return
		}
		deviceGroup := srcp.ExtractDeviceGroup(message.Message)
		switch deviceGroup {
		case "GL":
			bus, address := srcp.ExtractBusAndAddress(message.Message)
			if bus > -1 && address > -1 {
				var gl = store.GetGL(bus, address)
				if gl == nil && message.Code == 101 {
					gl = store.CreateGL(bus, address)
				}
				if gl != nil {
					if timestamp >= gl.LastTimestamp {
						srcp.UpdateGeneralLoco(message.Code, message.Message, gl)
						switch message.Code {
						case 100:
							if err := websocket.WriteJSON(InfoMessage{"GL updated", "update", Data{fmt.Sprintf("%d-%d", bus, address), "gl", gl}}); err != nil {
								panic(err)
							}
						case 101:
							if err := websocket.WriteJSON(InfoMessage{"GL created", "create", Data{fmt.Sprintf("%d-%d", bus, address), "gl", gl}}); err != nil {
								panic(err)
							}
						case 102:
							if err := websocket.WriteJSON(InfoMessage{"GL deleted", "delete", Data{fmt.Sprintf("%d-%d", bus, address), "gl", gl}}); err != nil {
								panic(err)
							}
						}
					}
					gl.LastTimestamp = timestamp
				}
			}
		case "GM":
			gm, err := srcp.ExtractGM(message.Message)
			if err != nil {
				panic(err)
			}
			var data Data
			if err := json.Unmarshal([]byte(gm.Message), &data); err != nil {
				log.Printf("json: %s", gm.Message)
				panic(err)
			}
			switch data.Type {
			case "gl":
				if err := websocket.WriteJSON(InfoMessage{"GL updated", "update", data}); err != nil {
					panic(err)
				}
			}
		}
	}
}
