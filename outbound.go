package cloudmailin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

// OutboundMail represents an email message ready to be sent.
// The ID will be populated by the API call once the message has been sent.
type OutboundMail struct {
	From        string                   `json:"from"`
	To          []string                 `json:"to,omitempty"`
	CC          []string                 `json:"cc,omitempty"`
	Headers     map[string][]string      `json:"headers,omitempty"`
	Subject     string                   `json:"subject,omitempty"`
	Plain       string                   `json:"plain,omitempty"`
	HTML        string                   `json:"html,omitempty"`
	Priority    string                   `json:"priorty,omitempty"`
	Tags        []string                 `json:"tags,omitempty"`
	Attachments []OutboundMailAttachment `json:"attachments,omitempty"`

	ID string `json:"id,omitempty"`
}

// OutboundMailAttachment represents the format of attachments to be sent
// in an OutboundMail. Content must be a Base64 encoded string.
type OutboundMailAttachment struct {
	// The Base64 encoded representation of the content.
	Content string `json:"content"`

	// An optional content id for the embedded attachment
	ContentID string `json:"content_id,omitempty"`

	// The mime content type of the file such as `image/jpeg`
	ContentType string `json:"content_type"`

	// The file name of the attachment
	FileName string `json:"file_name"`
}

// SendMail will make a POST to send the OutboundMail email via the HTTP API.
func (client Client) SendMail(message *OutboundMail) (res *http.Response, err error) {
	res, err = client.Do("POST", "/messages", message, RequestTypeSMTP)
	if err != nil {
		return
	}

	if res.StatusCode != 202 {
		body, _ := ioutil.ReadAll(res.Body)
		err = fmt.Errorf("could not send message (%d): %s", res.StatusCode, body)
	} else {
		defer res.Body.Close()
		err = json.NewDecoder(res.Body).Decode(message)
	}

	return
}

// AttachmentFromFile is a convenience function to prepare an OutboundMailAttachment
// from a local file (given as the filepath argument).
// The content will be Base64 encoded automatically and the filename included.
// In future this might be extended to help with content types.
func AttachmentFromFile(filepath string) (att OutboundMailAttachment,
	err error) {

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}

	// encoded := make([]byte)
	encoded := base64.StdEncoding.EncodeToString(bytes)
	filename := path.Base(filepath)

	att = OutboundMailAttachment{
		Content:     encoded,
		ContentID:   "",
		ContentType: "",
		FileName:    filename,
	}

	return
}
