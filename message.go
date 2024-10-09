package main

type Message struct {
	UserId  string
	Message string
}

type WSMessage struct {
	Text    string      `json:"text"`
	Headers interface{} `json:"Headers"`
}
