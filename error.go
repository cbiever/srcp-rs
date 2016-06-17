package main

type SrcpError struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Text   string `json:"text"`
}
