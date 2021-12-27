package cloudmailin_test

import (
	"fmt"

	"github.com/cloudmailin/cloudmailin-go"
)

func ExampleClient_SendMail() {
	// Create the default CloudMailin Client. This example will
	// panic if there are any failures at all.
	client, err := cloudmailin.NewClient()
	if err != nil {
		panic(err)
	}

	// SMTP Settings will be taken from CLOUDMAILIN_SMTP_URL env variable by
	// default but they can be overridden.
	// client.SMTPAccountID = ""
	// client.SMTPToken = ""

	// Create an instance of cloudmailin.OutboundMailAttachment
	attachment, err := cloudmailin.AttachmentFromFile("./test/fixtures/pixel.png")
	if err != nil {
		panic(err)
	}

	// Generate an example email
	message := cloudmailin.OutboundMail{
		From:        "sender@example.com",
		To:          []string{"debug@example.net"},
		CC:          []string{"carbon@example.net"},
		Headers:     map[string][]string{"x-agent": {"cloudmailin-go"}},
		Subject:     "Hello From Go",
		Plain:       "Hello World",
		HTML:        "<h1>Hello!</h1>\nWorld",
		Priority:    "",
		Tags:        []string{"go"},
		Attachments: []cloudmailin.OutboundMailAttachment{attachment},
		TestMode:    true,
	}

	// This will re-write the message struct based on the
	// JSON returned from the call if successful.
	_, err = client.SendMail(&message)
	if err != nil {
		panic(err)
	}

	// The message.ID should now be populated
	fmt.Printf("ID: %t, Tags: %s", message.ID != "", message.Tags)
	// Output: ID: true, Tags: [go api]
}
