package srcp

import (
	"errors"
	"fmt"
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

type GMMessage struct {
	Bus         int
	SendTo      int
	ReplyTo     int
	MessageType string
	Message     string
}

var messagePattern = regexp.MustCompile(`(\d{10}\.\d{3}) (\d{3}) ([A-Z]+)[ ]{0,1}(.*)`)
var sessionIdPattern = regexp.MustCompile(`GO (\d+)`)
var deviceGroupPattern = regexp.MustCompile(`\d+ (\w+)`)
var busAndAddressPattern = regexp.MustCompile(`(\d+) \w+ (\d+)`)
var glInitPattern = regexp.MustCompile(`(\d+) GL (\d+) ([\w]) (\d+) (\d+) (\d+)`)
var glDescriptionPattern = regexp.MustCompile(`(\d+) GL (\d+) (\d) (\d+) (\d+)([ \d]*)`)
var gmPattern = regexp.MustCompile(`(\d+) GM (\d+) (\d+) (\w+) (.*)`)

func Parse(message string) SrcpMessage {
	var srcpMessage SrcpMessage
	if result := messagePattern.FindStringSubmatch(message); result != nil {
		srcpMessage.Time = result[1]
		srcpMessage.Code, _ = strconv.Atoi(result[2])
		srcpMessage.Status = result[3]
		srcpMessage.Message = result[4]
	} else {
		srcpMessage.Message = message
	}
	return srcpMessage
}

func (message *SrcpMessage) ExtractSessionId() int {
	sessionId, _ := strconv.Atoi(sessionIdPattern.FindStringSubmatch(message.Message)[1])
	return sessionId
}

func (message *SrcpMessage) ExtractSessionInfos() map[string]string {
	infos := make(map[string]string)
	for _, info := range strings.Split(message.Message, ";") {
		keyValue := strings.Split(strings.Trim(info, " "), " ")
		infos[keyValue[0]] = strings.TrimSpace(keyValue[1])
	}
	return infos
}

func (message *SrcpMessage) ExtractBusAndAddress() (bus int, address int) {
	bus = -1
	address = -1
	result := busAndAddressPattern.FindStringSubmatch(message.Message)
	if result != nil {
		bus, _ = strconv.Atoi(result[1])
		address, _ = strconv.Atoi(result[2])
	}
	return bus, address
}

func (message *SrcpMessage) ExtractDeviceGroup() string {
	result := deviceGroupPattern.FindStringSubmatch(message.Message)
	if result != nil {
		return result[1]
	} else {
		return ""
	}
}

func (message *SrcpMessage) ExtractGLInitValues() (bus int, address int, protocol string, protocolVersion int, decoderSpeedSteps int, numberOfDecoderFunctions int, err error) {
	if result := glInitPattern.FindStringSubmatch(message.Message); result != nil {
		bus, _ := strconv.Atoi(result[1])
		address, _ := strconv.Atoi(result[2])
		protocol := result[3]
		protocolVersion, _ := strconv.Atoi(result[4])
		decoderSpeedSteps, _ := strconv.Atoi(result[5])
		numberOfDecoderFunctions, _ := strconv.Atoi(result[6])
		return bus, address, protocol, protocolVersion, decoderSpeedSteps, numberOfDecoderFunctions, nil
	} else {
		return 0, 0, "", 0, 0, 0, errors.New(fmt.Sprintf("Unable to parse: %s", message))
	}
}

func (message *SrcpMessage) ExtractGLDescriptionValues() (bus int, address int, drivemode int, V int, Vmax int, function []int, err error) {
	if result := glDescriptionPattern.FindStringSubmatch(message.Message); result != nil {
		bus, _ = strconv.Atoi(result[1])
		address, _ = strconv.Atoi(result[2])
		drivemode, _ = strconv.Atoi(result[3])
		V, _ = strconv.Atoi(result[4])
		Vmax, _ = strconv.Atoi(result[5])
		function := make([]int, 0)
		for _, value := range strings.Split(strings.Trim(result[6], " "), " ") {
			f, _ := strconv.Atoi(value)
			function = append(function, f)
		}
		return bus, address, drivemode, V, Vmax, function, nil
	} else {
		return 0, 0, 0, 0, 0, nil, errors.New(fmt.Sprintf("Unable to parse: %s", message))
	}
}

func (message *SrcpMessage) ExtractDeviceGroups() []string {
	var deviceGroups []string
	for index, deviceGroup := range strings.Split(message.Message, " ") {
		if index > 1 {
			deviceGroups = append(deviceGroups, deviceGroup)
		}
	}
	return deviceGroups
}

func ParseGM(message string) (GMMessage, error) {
	var gm GMMessage
	if result := gmPattern.FindStringSubmatch(message); result != nil {
		gm.Bus, _ = strconv.Atoi(result[1])
		gm.SendTo, _ = strconv.Atoi(result[2])
		gm.ReplyTo, _ = strconv.Atoi(result[3])
		gm.MessageType = result[4]
		gm.Message = result[5]
		return gm, nil
	}
	return gm, errors.New(fmt.Sprintf("Unable to parse: %s", message))
}
