package cloudmailin_test

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
func ExampleParseIncoming() {
	http.HandleFunc("/", handleIncomingPOST)

	// Start the HTTP server to listen for HTTP POST
	http.ListenAndServe(":8080", nil)
}
