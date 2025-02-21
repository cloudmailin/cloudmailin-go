package cloudmailin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	defaultBaseURL = "https://api.cloudmailin.com/api/v0.1"
)

// Client is the base of an API Client for CloudMailin. It provides a struct to
// put credentials, base URL and the HTTP Client.
type Client struct {
	HTTPClient    http.Client
	BaseURL       string
	SMTPToken     string
	SMTPAccountID string

	// For future use with the API
	AccountID    string
	AccountToken string
}

// NewClientFromURL returns an instance of the Client using the provided SMTP URL.
// The BaseURL will be set from CLOUDMAILIN_API_BASE_URL environment variable or
// fall back to the default if not set.
func NewClientFromURL(smtpURL string) (client Client, err error) {
	var smtpToken string
	var smtpAccount string

	if smtpURL != "" {
		u, err := url.Parse(smtpURL)
		if err != nil {
			return client, err
		}
		smtpAccount = u.User.Username()
		smtpToken, _ = u.User.Password()
	}

	baseURL := os.Getenv("CLOUDMAILIN_API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	client = Client{
		HTTPClient:    *http.DefaultClient,
		BaseURL:       baseURL,
		SMTPToken:     smtpToken,
		SMTPAccountID: smtpAccount,
		AccountID:     "",
		AccountToken:  "",
	}

	if client.SMTPAccountID == "" || client.SMTPToken == "" {
		err = fmt.Errorf("missing client values: %v", client)
	}

	return
}

// NewClient returns an instance of the DefaultClient. DefaultClient assumes
// that the ENV variable CLOUDMAILIN_SMTP_URL is set containing the SMTP
// credentials and sets the BaseURL to CLOUDMAILIN_API_BASE_URL if present
// or the default if not set.
//
// If the CLOUDMAILIN_SMTP_URL is not set either use NewClientFromURL or
// create a client instance manually and set the required credentials.
func NewClient() (client Client, err error) {
	smtpURL := os.Getenv("CLOUDMAILIN_SMTP_URL")
	return NewClientFromURL(smtpURL)
}

// RequestType is used to differentiate between account and SMTP requests.
type RequestType string

const (
	// RequestTypeAccount is used for account API calls
	// (currently not implemented).
	RequestTypeAccount RequestType = "account"

	// RequestTypeSMTP is used for SMTP related API calls.
	RequestTypeSMTP RequestType = "smtp"
)

// Do Performs the HTTP request with SMTP or Account credentials,
// SMTP credentials work in a really similar way to normal API credentials
// but use the SMTP Server ID instead of the AccountID and the SMTP Password
// as the Authorization Bearer Token.
func (client Client) Do(method string, path string, body interface{}, kind RequestType) (
	res *http.Response, err error) {

	var account, token string

	if kind == RequestTypeSMTP {
		account = client.SMTPAccountID
		token = client.SMTPToken
	} else {
		account = client.AccountID
		token = client.AccountToken
	}

	if account == "" || token == "" {
		err = fmt.Errorf("account or token credentials are missing for request (%s, %d)",
			account, len(token))
		return
	}

	url := strings.TrimSuffix(client.BaseURL, "/") + "/" + account + "/" +
		strings.TrimPrefix(path, "/")

	msgJSON, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msgJSON))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err = client.HTTPClient.Do(req)
	if err != nil {
		return
	}

	return
}
