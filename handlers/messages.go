package handlers

type Data struct {
	Id         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes"`
}

type Wrapper struct {
	Data `json:"data"`
}

type InfoMessage struct {
	Message string `json:"message"`
	Action  string `json:"action"`
	Data    Data   `json:"data"`
}

type Session struct {
	SessionId int               `json:"session-id"`
	Mode      string            `json:"mode"`
	Infos     map[string]string `json:"infos"`
}

type SrcpError struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Text   string `json:"text"`
}
