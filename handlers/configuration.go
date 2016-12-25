package handlers

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"srcp-rs/model"
)

func GetConfiguration(w http.ResponseWriter, r *http.Request) {
	configurationData, error := yaml.Marshal(store.GetGLS())
	if error != nil {
		panic(error)
	}
	io.WriteString(w, string(configurationData))
}

func UpdateConfiguration(w http.ResponseWriter, r *http.Request) {
	configurationData, error := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if error != nil {
		panic(error)
	}
	configuration := make(map[int]map[int]*model.GeneralLoco)
	if error = yaml.Unmarshal(configurationData, &configuration); error != nil {
		panic(error)
	}
	session, _, _ := extract(r)
	srcpConnection := store.GetConnection(session)
	for bus, gls := range configuration {
		for address, gl := range gls {
			store.SaveGL(bus, address, gl)
			srcpConnection.SendAndReceive(fmt.Sprintf("INIT %d GL %d %s %d %d %d", bus, address, gl.Protocol, gl.ProtocolVersion, gl.DecoderSpeedSteps, gl.NumberOfDecoderFunctions))
			data, error := json.Marshal(Data{fmt.Sprintf("%d-%d", bus, address), "gl", gl})
			if error != nil {
				panic(error)
			}
			srcpConnection.SendAndReceive(fmt.Sprintf("SET 0 GM 0 0 TEXT %s", string(data)))
		}
	}
}
