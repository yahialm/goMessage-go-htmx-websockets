package main

type Message struct {
	userId  string
	message string
}

type WSMessage struct {
	Text    string      `json:"text"`
	Headers interface{} `json:"Headers"`
}
