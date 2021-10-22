package cloudmailin

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseIncomingReader(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		data, _ := os.Open("test/fixtures/post.json")
		defer data.Close()
		_, err := ParseIncoming(data)
		if err != nil {
			t.Error("Expecetd error to be nil but was {}", err)
		}
	})

	// Test reader with bad file
	t.Run("Invalid File", func(t *testing.T) {
		data, _ := os.Open("test/fixtures/nope.json")
		defer data.Close()
		_, err := ParseIncoming(data)
		if err == nil {
			t.Error("Expecetd error but was nil", err)
		}
	})

	// Test reader with invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		data := strings.NewReader("{ invalid")
		message, err := ParseIncoming(data)
		if err == nil {
			t.Error("Expecetd error but was nil", message)
		}
	})

	data, _ := os.Open("test/fixtures/post.json")
	defer data.Close()
	message, _ := ParseIncoming(data)

	tests := []struct {
		name     string
		item     interface{}
		expected interface{}
	}{
		{"Envelope", message.Envelope, IncomingMailEnvelope{
			To:         "postman@cloudmailin.net",
			From:       "from+test@cloudmailin.net",
			Recipients: []string{"postman@cloudmailin.net"},
			HeloDomain: "cloudmailin.net",
			RemoteIP:   "172.20.0.18",
			TLS:        true,
			TLSCipher:  "TLSv1.3",
			StoreURL:   "http://example.s3.amazonaws.com/store/2020_10_22_09_55_18_ce5f9a939358ba89b80acd97f737e0db.eml",
			SPF:        IncomingMailEnvelopeSPF{"fail", "cloudmailin.net"},
			SPAMD:      IncomingMailEnvelopeSPAMD{},
		}},
		// Headers tested in the other specs
		{"Plain", message.Plain, "Test Content\n\n\u003e On 08 Jul 2020 at 10:00, " +
			"example@cloudmailin.net wrote:\n\u003e \n\u003e Example message\n\u003e " +
			"Option: 2\n\u003e \n\n"},
		{"HTML", message.HTML, "<div dir=\"ltr\">Test Content<div><br></div></div>\n"},
		{"ReplyPlain", message.ReplyPlain, "Test Content\n"},
		{"Attachments", message.Attachments, []IncomingMailAttachment{
			{
				Content:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP0rdr1HwAFHwKCk87e6gAAAABJRU5ErkJggg==",
				FileName:    "pixel.png",
				ContentType: "image/png",
				Size:        37135,
				Disposition: "attachment",
				ContentID:   "<f_kcd6ejvs1>",
				URL:         "",
				Scan: IncomingMailAttachmentScan{
					Status:  "ok",
					ID:      "b1a187b7-bd22-406a-8c1a-f532a0ce5ba6",
					Matches: []string{"Example Match"},
				},
			},
		}},
	}

	for _, tt := range tests {
		t.Run("Struct "+tt.name, func(t *testing.T) {
			if !cmp.Equal(tt.item, tt.expected) {
				t.Errorf("Expected {%v} but was {%v}\n{%v}", tt.expected, tt.item,
					cmp.Diff(tt.expected, tt.item))
			}
		})
	}
}

func TestParseIncomingBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		data, _ := os.ReadFile("test/fixtures/post.json")
		_, err := ParseIncomingBytes(data)
		if err != nil {
			t.Error("Expecetd error to be nil but was {}", err)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		data := "{"
		_, err := ParseIncomingBytes([]byte(data))

		if err == nil {
			t.Error("Expecetd error but was nil", err)
		}
	})
}

func TestIncomingMailHeader_UnmarshalJSON(t *testing.T) {
	// Test unmarshall with invalid JSON
	t.Run("Invalid Type", func(t *testing.T) {
		data := `{`
		header := IncomingMailHeader{}
		err := header.UnmarshalJSON([]byte(data))
		if err == nil {
			t.Error("Expected error but was nil", header)
		}
	})

	// Test unmarshall with invalid type
	t.Run("Invalid Type", func(t *testing.T) {
		data := `1`
		header := IncomingMailHeader{}
		err := header.UnmarshalJSON([]byte(data))
		if err == nil {
			t.Error("Expected error but was nil", header)
		}
	})
}

func TestIncomingMailHeaders_Find(t *testing.T) {
	headers := IncomingMailHeaders{
		"to":       {"Example User <user@example.com>"},
		"received": {"Received 2", "Received 1"},
	}

	tests := []struct {
		name     string
		item     IncomingMailHeader
		expected IncomingMailHeader
	}{
		{"to Header", headers.Find("to"), IncomingMailHeader{"Example User <user@example.com>"}},
		// received 1 (first) bottom of email
		{"received Header", headers.Find("received"), IncomingMailHeader{"Received 2", "Received 1"}},
		{"Missing Header", headers.Find("foo"), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !cmp.Equal(tt.expected, tt.item) {
				t.Errorf("Expected {%v} but was {%v}\n{%v}", tt.expected, tt.item,
					cmp.Diff(tt.expected, tt.item))
			}
		})
	}
}

func TestIncomingMailHeaders_First(t *testing.T) {
	headers := IncomingMailHeaders{
		"to":       {"Example User <user@example.com>"},
		"received": {"Received 2", "Received 1"},
	}

	tests := []struct {
		name     string
		item     string
		expected string
	}{
		{"to Header", headers.First("to"), "Example User <user@example.com>"},
		{"received Header", headers.First("received"), "Received 1"},
		{"Missing Header", headers.First("foo"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !cmp.Equal(tt.expected, tt.item) {
				t.Errorf("Expected {%v} but was {%v}\n{%v}", tt.expected, tt.item,
					cmp.Diff(tt.expected, tt.item))
			}
		})
	}
}

func TestIncomingMailHeaders_Last(t *testing.T) {
	headers := IncomingMailHeaders{
		"to":       {"Example User <user@example.com>"},
		"received": {"Received 2", "Received 1"},
	}

	tests := []struct {
		name     string
		item     string
		expected string
	}{
		{"to Header", headers.Last("to"), "Example User <user@example.com>"},
		{"received Header", headers.Last("received"), "Received 2"},
		{"Missing Header", headers.Last("foo"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !cmp.Equal(tt.expected, tt.item) {
				t.Errorf("Expected {%v} but was {%v}\n{%v}", tt.expected, tt.item,
					cmp.Diff(tt.expected, tt.item))
			}
		})
	}
}

func TestIncomingMailHeaders_Helpers(t *testing.T) {
	data, _ := os.Open("test/fixtures/post.json")
	defer data.Close()
	message, _ := ParseIncoming(data)

	tests := []struct {
		name     string
		item     string
		expected string
	}{
		{"To", message.Headers.To(), "postman@cloudmailin.net"},
		{"From", message.Headers.From(), "Steve Smith <test@example.com>"},
		{"Subject", message.Headers.Subject(), "Test Email"},
		{
			"MessageID",
			message.Headers.MessageID(),
			"<CALazKR8Zr8Lsv+SUAeuaL-vrhWSCK36TRU8=7HjsenxwaP9ZbA@mail.gmail.com>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !cmp.Equal(tt.expected, tt.item) {
				t.Errorf("Expected {%v} but was {%v}\n{%v}", tt.expected, tt.item,
					cmp.Diff(tt.expected, tt.item))
			}
		})
	}
}
