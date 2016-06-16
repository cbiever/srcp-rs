package main

type SessionRequest struct {
	Mode string `json:"mode"`
}

type SessionResponse struct {
	SessionId int               `json:"sessionId"`
	Infos     map[string]string `json:"infos"`
}

type GeneralLocoCreateRequest struct {
	Bus                      int    `json:"bus"`
	Address                  int    `json:"address"`
	Protocol                 string `json:"protocol"`
	ProtocalVersion          int    `json:"protocolVersion"`
	DecoderSpeedSteps        int    `json:"decoderSpeedSteps"`
	NumberOfDecoderFunctions int    `json:"numberOfDecoderFunctions"`
}

type GeneralLocoCreateResponse struct {
	Time string `json:"time"`
}

type GeneralLocoGetRequest struct {
	SessionId int `json:"sessionId"`
	Bus       int `json:"bus"`
	Address   int `json:"address"`
}

type GeneralLocoGetResponse struct {
	Time      string `json:"time"`
	Drivemode int    `json:"drivemode"`
	V         int    `json:"v"`
	Vmax      int    `json:"vMax"`
	Function  []int  `json:"functions"`
}

type GeneralLocoUpdateRequest struct {
	Drivemode int   `json:"drivemode"`
	V         int   `json:"v"`
	Vmax      int   `json:"vMax"`
	Function  []int `json:"functions"`
}

type GeneralLocoUpdateResponse struct {
	Time string `json:"time"`
}
