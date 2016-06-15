package main

type SessionRequest struct {
	Mode string `json:"mode"`
}

type SessionResponse struct {
	SessionId int               `json:"sessionId"`
	Infos     map[string]string `json:"infos"`
}

type GeneralLocoGetRequest struct {
	SessionId int `json:"sessionId"`
	Bus       int `json:"bus"`
	Address   int `json:"address"`
}

type GeneralLocoGetResponse struct {
	Time string `json:"time"`
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
