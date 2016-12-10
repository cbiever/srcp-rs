package srcp

import (
	"regexp"
	"strconv"
	"strings"
)

type SrcpMessage struct {
	Time    string
	Code    int
	Status  string
	Message string
}

type Bus struct {
	DeviceGroups []string `json:"deviceGroups"`
}

type GeneralLoco struct {
	Bus                      int     `json:"bus" yaml:"-"`
	Address                  int     `json:"address" yaml:"-"`
	Protocol                 string  `json:"protocol" yaml:"protocol"`
	ProtocolVersion          int     `json:"protocol-version" yaml:"protocol-version"`
	DecoderSpeedSteps        int     `json:"decoder-speed-steps" yaml:"decoder-speed-steps"`
	NumberOfDecoderFunctions int     `json:"number-of-decoder-functions" yaml:"number-of-decoder-functions"`
	Drivemode                int     `json:"drivemode" yaml:"-"`
	V                        int     `json:"v" yaml:"-"`
	Vmax                     int     `json:"v-max" yaml:"v-max"`
	Function                 []int   `json:"functions" yaml:"-"`
	LastTimestamp            float64 `json:"-" yaml:"-"`
}

var messagePattern = regexp.MustCompile(`(\d{10}\.\d{3}) (\d{3}) ([A-Z]+)[ ]{0,1}([\w ]*)`)
var sessionIdPattern = regexp.MustCompile(`GO (\d+)`)
var deviceGroupPattern = regexp.MustCompile(`\d+ (\w+)`)
var busAndAddressPattern = regexp.MustCompile(`(\d+) \w+ (\d+)`)
var glInitPattern = regexp.MustCompile(`(\d+) GL (\d+) ([\w]) (\d+) (\d+) (\d+)`)
var glDescriptionPattern = regexp.MustCompile(`(\d+) GL (\d+) (\d) (\d+) (\d+)([ \d]*)`)

func Parse(message string) SrcpMessage {
	result := messagePattern.FindStringSubmatch(message)
	var srcpMessage SrcpMessage
	srcpMessage.Time = result[1]
	srcpMessage.Code, _ = strconv.Atoi(result[2])
	srcpMessage.Status = result[3]
	srcpMessage.Message = result[4]
	return srcpMessage
}

func ExtractSessionId(message string) int {
	sessionId, _ := strconv.Atoi(sessionIdPattern.FindStringSubmatch(message)[1])
	return sessionId
}

func ExtractSessionInfos(message string) map[string]string {
	infos := make(map[string]string)
	for _, info := range strings.Split(message, ";") {
		keyValue := strings.Split(strings.Trim(info, " "), " ")
		infos[keyValue[0]] = strings.TrimSpace(keyValue[1])
	}
	return infos
}

func ExtractBusAndAddress(message string) (bus int, address int) {
	bus = -1
	address = -1
	result := busAndAddressPattern.FindStringSubmatch(message)
	if result != nil {
		bus, _ = strconv.Atoi(result[1])
		address, _ = strconv.Atoi(result[2])
	}
	return bus, address
}

func ExtractDeviceGroup(message string) string {
	result := deviceGroupPattern.FindStringSubmatch(message)
	if result != nil {
		return result[1]
	} else {
		return ""
	}
}

func UpdateGeneralLoco(code int, message string, gl *GeneralLoco) {
	if code == 101 {
		if result := glInitPattern.FindStringSubmatch(message); result != nil {
			gl.Bus, _ = strconv.Atoi(result[1])
			gl.Address, _ = strconv.Atoi(result[2])
			gl.Protocol = result[3]
			gl.ProtocolVersion, _ = strconv.Atoi(result[4])
			gl.DecoderSpeedSteps, _ = strconv.Atoi(result[5])
			gl.NumberOfDecoderFunctions, _ = strconv.Atoi(result[6])
		}
	} else if code == 100 {
		if result := glDescriptionPattern.FindStringSubmatch(message); result != nil {
			gl.Bus, _ = strconv.Atoi(result[1])
			gl.Address, _ = strconv.Atoi(result[2])
			gl.Drivemode, _ = strconv.Atoi(result[3])
			gl.V, _ = strconv.Atoi(result[4])
			gl.Vmax, _ = strconv.Atoi(result[5])
			functions := strings.Split(strings.Trim(result[6], " "), " ")
			if gl.Function == nil {
				gl.Function = make([]int, len(functions))
			}
			for i, function := range functions {
				gl.Function[i], _ = strconv.Atoi(function)
			}
		}
	}
}

func ExtractDeviceGroups(message string) []string {
	var deviceGroups []string
	for index, deviceGroup := range strings.Split(message, " ") {
		if index > 1 {
			deviceGroups = append(deviceGroups, deviceGroup)
		}
	}
	return deviceGroups
}
