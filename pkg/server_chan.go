package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ServerChanPayload struct {
	Title string `json:"title"`
	Desp  string `json:"desp"`
}

type ServerChan struct {
	sendKey string
}

func NewServerChan(sendKey string) *ServerChan {
	return &ServerChan{
		sendKey: sendKey,
	}
}

func (s *ServerChan) Push(title string, msg string) {
	payload := &ServerChanPayload{
		Title: title,
		Desp:  msg,
	}
	payloadJSON, _ := json.Marshal(payload)
	resp, err := http.Post(
		fmt.Sprintf("https://sctapi.ftqq.com/%s.send", s.sendKey),
		"application/json",
		strings.NewReader(string(payloadJSON)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Push: payload: %v, resp: %v\n", payload, resp)
}
