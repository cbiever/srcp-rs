package main

type Data struct {
	Id         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes"`
}
type Wrapper struct {
	Data `json:"data"`
}

type Session struct {
	SessionId int               `json:"sessionId"`
	Mode      string            `json:"mode"`
	Infos     map[string]string `json:"infos"`
}

type GeneralLoco struct {
	Address                  int    `json:"address"`
	Protocol                 string `json:"protocol"`
	ProtocalVersion          int    `json:"protocolVersion"`
	DecoderSpeedSteps        int    `json:"decoderSpeedSteps"`
	NumberOfDecoderFunctions int    `json:"numberOfDecoderFunctions"`
	Drivemode                int    `json:"drivemode"`
	V                        int    `json:"v"`
	Vmax                     int    `json:"vMax"`
	Function                 []int  `json:"functions"`
}
