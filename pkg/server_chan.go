package pkg

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const messagePlaceholder = "{message}"

type ServerChan struct {
	webhookURL string
}

func NewServerChan(webhookURL string) *ServerChan {
	return &ServerChan{
		webhookURL: webhookURL,
	}
}

func (s *ServerChan) Push(title string, msg string) {
	body := mergeTitleAndMessage(title, msg)
	webhookURL, err := buildWebhookURL(s.webhookURL, body)
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(webhookURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf("Push: resp: %v\n", resp)
}

func buildWebhookURL(webhookURL string, message string) (string, error) {
	if !strings.Contains(webhookURL, messagePlaceholder) {
		return "", fmt.Errorf("webhook-url must contain %s placeholder", messagePlaceholder)
	}
	return strings.ReplaceAll(webhookURL, messagePlaceholder, url.QueryEscape(message)), nil
}

func mergeTitleAndMessage(title string, msg string) string {
	if title == "" {
		return msg
	}
	if msg == "" {
		return title
	}
	return fmt.Sprintf("%s\n%s", title, msg)
}
