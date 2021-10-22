<a href="https://www.cloudmailin.com">
  <img src="https://assets.cloudmailin.com/assets/favicon.png" alt="CloudMailin Logo" height="60" align="right" title="CloudMailin">
</a>

# CloudMailin Go Package

[![Go Reference](https://pkg.go.dev/badge/github.com/cloudmailin/cloudmailin-go.svg)](https://pkg.go.dev/github.com/cloudmailin/cloudmailin-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudmailin/cloudmailin-go)](https://goreportcard.com/report/github.com/cloudmailin/cloudmailin-go)

A Go SDK for CloudMailin incoming and outbound email via JSON HTTP POST.

Please see the [Documentation](https://docs.cloudmailin.com) for more details and examples.

## Usage

### Installation

You can install the package using:

```console
go get -u github.com/cloudmailin/cloudmailin-go
```

### Receiving Email

We recommend you take a look at our
[Documentation](https://docs.cloudmailin.com/receiving_email/examples/go/)
for a more detailed example but here's a snippet:

```go
import (
	"fmt"
	"net/http"

	"github.com/cloudmailin/cloudmailin-go"
)

func handleIncomingPOST(w http.ResponseWriter, req *http.Request) {
	// Parse the message from the request body
	message, err := cloudmailin.ParseIncoming(req.Body)
	if err != nil {
		// Print an error message if parsing fails
		fmt.Fprint(w, "Error parsing message: ", err)
	}

	// Output the first instance of the message-id in the headers to show
	// that we correctly parsed the message. We could also use the helper
	// message.Headers.MessageID().
	fmt.Fprint(w, "Thanks for message: ", message.Headers.First("message_id"))
}

// This example shows how to parse the incoming JSON sent by Cloudmailin
// and create an instance of the IncomingMail type.
func main() {
	http.HandleFunc("/", handleIncomingPOST)

	// Start the HTTP server to listen for HTTP POST
	http.ListenAndServe(":8080", nil)
}
```

### Sending Email

We recommend you take a look at our
[Documentation](https://docs.cloudmailin.com/outbound/examples/send_email_with_golang/)
for a more detailed example:

```go
package main

import (
  "log"
  "net/smtp"
)

func main() {
  // hostname is used by PlainAuth to validate the TLS certificate.
  hostname := "host from account"
  auth := smtp.PlainAuth("", "username from account", "password fromo account", hostname)

  msg := `To: to@example.net
From: from@example.com
Subject: Testing from GoLang

This is the message content!
Thanks
`
  err := smtp.SendMail(hostname+":587", auth, "from@example.com", []string{"to@example.net"},
    []byte(msg))
  if err != nil {
    log.Fatal(err)
  }
}
```

## Development

Although we experimented with code generators none of them provided the
experience that we desired.

This code was built inside a docker container in VSCode. Contact us if you
would like to make use of any of those tools.
