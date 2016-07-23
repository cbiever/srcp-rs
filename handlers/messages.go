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
	SessionId int               `json:"session-id"`
	Mode      string            `json:"mode"`
	Infos     map[string]string `json:"infos"`
}

type GeneralLoco struct {
	Address                  int    `json:"address"`
	Protocol                 string `json:"protocol"`
	ProtocalVersion          int    `json:"protocol-version"`
	DecoderSpeedSteps        int    `json:"decoder-speed-steps"`
	NumberOfDecoderFunctions int    `json:"number-of-decoder-functions"`
	Drivemode                int    `json:"drivemode"`
	V                        int    `json:"v"`
	Vmax                     int    `json:"v-max"`
	Function                 []int  `json:"functions"`
}
