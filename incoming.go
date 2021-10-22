package cloudmailin

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// IncomingMail represents an email received via HTTP from the
// CloudMailin email to Webhook.
type IncomingMail struct {
	Envelope    IncomingMailEnvelope     `json:"envelope"`
	Headers     IncomingMailHeaders      `json:"headers"`
	Plain       string                   `json:"plain"`
	HTML        string                   `json:"html"`
	ReplyPlain  string                   `json:"reply_plain"`
	Attachments []IncomingMailAttachment `json:"attachments"`
}

// ParseIncoming parses an IO reader and emits an IncomingMail.
// This is most useful to parse the http.Request.Body.
func ParseIncoming(data io.Reader) (mail IncomingMail, err error) {
	err = json.NewDecoder(data).Decode(&mail)
	return
}

// ParseIncomingBytes parses a []byte and returns an IncomingMail.
func ParseIncomingBytes(data []byte) (mail IncomingMail, err error) {
	err = json.Unmarshal([]byte(data), &mail)
	return
}

// IncomingMailEnvelope represents the CloudMailin Email Envelope containing
// information passed to the SMTP server.
type IncomingMailEnvelope struct {
	To         string                    `json:"to"`
	From       string                    `json:"from"`
	Recipients []string                  `json:"recipients"`
	HeloDomain string                    `json:"helo_domain"`
	RemoteIP   string                    `json:"remote_ip"`
	TLS        bool                      `json:"tls"`
	TLSCipher  string                    `json:"tls_cipher"`
	MD5        string                    `json:"md5"`
	StoreURL   string                    `json:"store_url"`
	SPF        IncomingMailEnvelopeSPF   `json:"spf"`
	SPAMD      IncomingMailEnvelopeSPAMD `json:"spamd,omitempty"`
}

// IncomingMailEnvelopeSPF contains the SPF result and domain for the Envelope.
type IncomingMailEnvelopeSPF struct {
	Result string `json:"result"`
	Domain string `json:"domain"`
}

// IncomingMailEnvelopeSPAMD contains the result of the SpamAssassin scan of the
// email (if enabled).
type IncomingMailEnvelopeSPAMD struct {
	Score       uint     `json:"score,string"`
	Symbols     []string `json:"symbols"`
	Success     bool     `json:"success"`
	Description string   `json:"description"`
}

// IncomingMailHeaders contains the received email headers.
type IncomingMailHeaders map[string]IncomingMailHeader

//IncomingMailHeader contains a single header value. This is an array because
// each header can be received more than once.
type IncomingMailHeader []string

// UnmarshalJSON takes the CloudMailin JSON format and handles the parsing
// into the struct. It ensures that every response is an array of strings.
func (i *IncomingMailHeader) UnmarshalJSON(b []byte) error {
	var x interface{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	switch v := x.(type) {
	case string:
		*i = IncomingMailHeader([]string{v})
	case []interface{}:
		items := make([]string, 0)
		for _, v2 := range v {
			items = append(items, v2.(string))
		}
		*i = items
	default:
		return fmt.Errorf("unexpected type %v for value: {%v}", reflect.TypeOf(v), v)
	}

	return nil
}

// Find will return a header by it's name
// Email headers are ordered bottom up so they will be top
// first in this array.
func (i IncomingMailHeaders) Find(key string) IncomingMailHeader {
	return i[key]
}

// First will return the first entry for a header by it's name.
func (i IncomingMailHeaders) First(key string) string {
	return i.Find(key).First()
}

// Last will find and return the most recent entry for a header by it's name.
func (i IncomingMailHeaders) Last(key string) string {
	return i.Find(key).Last()
}

// To is a helper function to find the To Header.
func (i IncomingMailHeaders) To() string {
	return i.First("to")
}

// From is a helper function to find the From Header.
func (i IncomingMailHeaders) From() string {
	return i.First("from")
}

// Subject is a helper function to find the Subject Header.
func (i IncomingMailHeaders) Subject() string {
	return i.First("subject")
}

// MessageID is a helper function to find the Message-ID Header.
func (i IncomingMailHeaders) MessageID() string {
	return i.First("message_id")
}

// Last will return the occurrence of this header appearing at the Top
// of the headers i.e. the most recent header value.
func (i IncomingMailHeader) Last() string {
	if len(i) == 0 {
		return ""
	}
	return i[0]
}

// First will return the occurrence of this header appearing at the Bottom
// of the headers i.e. the first header value.
func (i IncomingMailHeader) First() string {
	if len(i) == 0 {
		return ""
	}
	return i[len(i)-1]
}

// IncomingMailAttachment represents email attachments.
type IncomingMailAttachment struct {
	// Base64 encoded string containing the attachment content
	Content     string `json:"content"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	Size        uint64 `json:"size,string"`
	Disposition string
	ContentID   string `json:"content_id"`
	URL         string
	Scan        IncomingMailAttachmentScan
}

// IncomingMailAttachmentScan represents the result of virus scanning the email
// attachments (if enabled).
type IncomingMailAttachmentScan struct {
	Status  string
	ID      string
	Matches []string
}
