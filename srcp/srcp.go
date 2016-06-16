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

type GLValues struct {
	Bus       int
	Address   int
	Drivemode int
	V         int
	Vmax      int
	Function  []int
}

var messagePattern = regexp.MustCompile(`(\d{10}\.\d{3}) (\d{3}) (\w+)[ ]{0,1}([\w ]*)`)
var sessionIdPattern = regexp.MustCompile(`GO (\d+)`)
var glValuesPattern = regexp.MustCompile(`(\d+) GL (\d+) (\d) (\d+) (\d+)([ \d]*)`)

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

func ExtractGLValues(message string) GLValues {
	result := glValuesPattern.FindStringSubmatch(message)
	var glValues GLValues
	glValues.Drivemode, _ = strconv.Atoi(result[3])
	glValues.V, _ = strconv.Atoi(result[4])
	glValues.Vmax, _ = strconv.Atoi(result[5])
	functions := strings.Split(strings.Trim(result[6], " "), " ")
	glValues.Function = make([]int, len(functions))
	for i, function := range functions {
		glValues.Function[i], _ = strconv.Atoi(function)
	}
	return glValues
}
