package pkg

import "testing"

func TestBuildWebhookURLReplacesMessagePlaceholder(t *testing.T) {
	got, err := buildWebhookURL("https://example.com/api/messages?token=token&message={message}", "first paragraph\n\nsecond paragraph")
	if err != nil {
		t.Fatalf("buildWebhookURL returned error: %v", err)
	}

	want := "https://example.com/api/messages?token=token&message=first+paragraph%0A%0Asecond+paragraph"
	if got != want {
		t.Fatalf("buildWebhookURL = %q, want %q", got, want)
	}
}

func TestBuildWebhookURLRequiresMessagePlaceholder(t *testing.T) {
	_, err := buildWebhookURL("https://example.com/api/messages?token=token", "message")
	if err == nil {
		t.Fatal("buildWebhookURL returned nil error")
	}
}

func TestMergeTitleAndMessage(t *testing.T) {
	got := mergeTitleAndMessage("IPMONITOR: INIT", "- lo: 127.0.0.1")
	want := "IPMONITOR: INIT\n- lo: 127.0.0.1"
	if got != want {
		t.Fatalf("mergeTitleAndMessage = %q, want %q", got, want)
	}
}
