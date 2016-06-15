package srcp

import (
	"regexp"
	"strconv"
)

type SrcpMessage struct {
	Time    string
	Code    int
	Status  string
	Message string
}

var format = regexp.MustCompile(`(\d{10}\.\d{3}) (\d{3}) (\w+)[ ]{0,1}([\w ]*)`)
var sessionId = regexp.MustCompile(`GO (\d+)`)

func Parse(message string) SrcpMessage {
	result := format.FindStringSubmatch(message)
	var srcpMessage SrcpMessage
	srcpMessage.Time = result[1]
	srcpMessage.Code, _ = strconv.Atoi(result[2])
	srcpMessage.Status = result[3]
	srcpMessage.Message = result[4]
	return srcpMessage
}

func ExtractSessionId(message string) int {
	sessionId, _ := strconv.Atoi(sessionId.FindStringSubmatch(message)[1])
	return sessionId
}
