package cloudmailin

import (
	"fmt"
	"strings"
	"testing"
)

func buildMessage() OutboundMail {
	attachmentContent := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP0rdr1HwAFHwKCk87e6gAAAABJRU5ErkJggg=="
	return OutboundMail{
		From:     "sender@example.com",
		To:       []string{"debug@example.net"},
		CC:       []string{"carbon@example.net"},
		Headers:  map[string][]string{"x-agent": {"cloudmailin-go"}},
		Subject:  "Hello From Go",
		Plain:    "Hello World",
		HTML:     "<h1>Hello!</h1>\nWorld",
		Priority: "",
		Tags:     []string{"go"},
		Attachments: []OutboundMailAttachment{
			{
				Content:     attachmentContent,
				ContentType: "image/png",
				FileName:    "pixel.png",
				ContentID:   "<pix12>",
			},
		},
	}
}

func TestClient_SendMessage(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		panic(err)
	}

	message := buildMessage()
	res, err := client.SendMail(&message)

	t.Run("Send", func(t *testing.T) {
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Returns the HTTP res", func(t *testing.T) {
		if res == nil || res.StatusCode != 202 {
			t.Errorf("Expected HTTP response got {%v}", res)
		}
	})

	t.Run("Sets the ID", func(t *testing.T) {
		if message.ID == "" {
			t.Errorf("Expected message-id to be present got {%v}", message.ID)
		}
	})
}

func TestClient_SendMessage_Fail(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		panic(err)
	}

	client.SMTPToken = ""

	message := buildMessage()
	res, err := client.SendMail(&message)

	if err == nil || !strings.Contains(err.Error(), "credentials") {
		t.Errorf("Expected error, got: {%v}, {%v}", res, err)
	}
}

func TestClient_SendMessage_404(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		panic(err)
	}

	client.BaseURL = "https://www.cloudmailin.com/target/404"

	message := buildMessage()
	res, err := client.SendMail(&message)

	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("Expected error, got: {%v}, {%v}", res, err)
	}
}

func TestOutboundMailAttachment_AttachmentFromFile(t *testing.T) {
	t.Run("Invalid file", func(t *testing.T) {
		path := "test/fixtures/missing.png"
		_, err := AttachmentFromFile(path)
		if err == nil {
			t.Errorf("Err should be present but was nil")
		}
	})

	t.Run("Valid file", func(t *testing.T) {
		path := "test/fixtures/pixel.png"
		att, err := AttachmentFromFile(path)
		if err != nil {
			fmt.Println(err)
		}

		expected := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP0rdr1HwAFHwKCk87e6gAAAABJRU5ErkJggg=="
		if att.Content != expected {
			t.Errorf("Expected base64 string got: {%v}", att.Content)
		}
	})
}
