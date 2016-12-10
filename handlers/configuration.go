package handlers

import (
	"srcp-rs/srcp"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"gopkg.in/yaml.v2"
)

func GetConfiguration(w http.ResponseWriter, r *http.Request) {
	configurationData, error := yaml.Marshal(store.GetGLS());
	if error != nil {
		panic(error)
	}
	io.WriteString(w, string(configurationData))
}

func UpdateConfiguration(w http.ResponseWriter, r *http.Request) {
	configurationData, error := ioutil.ReadAll(io.LimitReader(r.Body, 1048576));
	if error != nil {
		panic(error)
	}
	configuration := make(map[int]map[int]*srcp.GeneralLoco)
	if error = yaml.Unmarshal(configurationData, &configuration); error != nil {
		panic(error)
	}
	session, _, _ := extract(r)
	srcpConnection := store.GetConnection(session)
	for bus, gls := range(configuration) {
			for address, gl := range gls {
				srcpConnection.SendAndReceive(fmt.Sprintf("INIT %d GL %d %s %d %d %d", bus, address, gl.Protocol, gl.ProtocolVersion, gl.DecoderSpeedSteps, gl.NumberOfDecoderFunctions))
			}
	}
}
