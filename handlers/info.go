package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"srcp-rs/model"
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

	reply := srcpConnection.Receive()
	message := srcp.Parse(reply)

	session.Infos = message.ExtractSessionInfos()

	reply = srcpConnection.SendAndReceive("SET CONNECTIONMODE SRCP INFORMATION")

	if message := srcp.Parse(reply); message.Code != 202 {
		websocket.WriteJSON(SrcpError{message.Code, message.Status, message.Message})
		websocket.Close()
		return
	}

	reply = srcpConnection.SendAndReceive("GO")

	if message := srcp.Parse(reply); message.Code == 200 {
		session.SessionId = message.ExtractSessionId()
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
		deviceGroup := message.ExtractDeviceGroup()
		switch deviceGroup {
		case "GL":
			bus, address := message.ExtractBusAndAddress()
			if bus > -1 && address > -1 {
				var gl = store.GetGL(bus, address)
				if gl == nil && message.Code == 101 {
					gl = store.CreateGL(bus, address)
				}
				if gl != nil {
					if timestamp >= gl.LastTimestamp {
						updateGeneralLoco(&message, gl)
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
			gm, err := srcp.ParseGM(message.Message)
			if err != nil {
				panic(err)
			}
			var data Data
			if err := json.Unmarshal([]byte(gm.Message), &data); err != nil {
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

func updateGeneralLoco(message *srcp.SrcpMessage, gl *model.GeneralLoco) {
	var err error
	switch message.Code {
	case 101:
		if gl.Bus, gl.Address, gl.Protocol, gl.ProtocolVersion, gl.DecoderSpeedSteps, gl.NumberOfDecoderFunctions, err = message.ExtractGLInitValues(); err != nil {
			panic(err)
		}
	case 100:
		if gl.Bus, gl.Address, gl.Drivemode, gl.V, gl.Vmax, gl.Function, err = message.ExtractGLDescriptionValues(); err != nil {
			panic(err)
		}
	}
}
