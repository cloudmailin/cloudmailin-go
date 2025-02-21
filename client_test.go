package cloudmailin

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewClient_ENV_SMTP_URL(t *testing.T) {
	original := os.Getenv("CLOUDMAILIN_SMTP_URL")
	defer os.Setenv("CLOUDMAILIN_SMTP_URL", original)

	t.Run("Not a URL", func(t *testing.T) {
		os.Setenv("CLOUDMAILIN_SMTP_URL", "http://\\example.com/foo")
		_, err := NewClient()

		if err == nil {
			t.Error("Expected error got nil")
		}
	})

	t.Run("Missing params", func(t *testing.T) {
		os.Setenv("CLOUDMAILIN_SMTP_URL", "http://example.com/foo")
		_, err := NewClient()

		expected := "missing client values"
		if err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error got {%v}", err)
		}
	})

	t.Run("SMTPToken", func(t *testing.T) {
		os.Setenv("CLOUDMAILIN_SMTP_URL", "smtp://user:pass@localhost/path")
		client, _ := NewClient()

		if client.SMTPToken != "pass" {
			t.Error("Expected client token", client)
		}
	})

	t.Run("SMTPAccount", func(t *testing.T) {
		os.Setenv("CLOUDMAILIN_SMTP_URL", "smtp://user:pass@localhost/path")
		client, _ := NewClient()

		if client.SMTPAccountID != "user" {
			t.Error("Expected client account", client)
		}
	})
}

func TestNewClient_ENV_BASE_URL(t *testing.T) {
	original := os.Getenv("CLOUDMAILIN_API_BASE_URL")
	defer os.Setenv("CLOUDMAILIN_API_BASE_URL", original)

	t.Run("Default (no ENV)", func(t *testing.T) {
		expected := "https://api.cloudmailin.com/api/v0.1"
		os.Setenv("CLOUDMAILIN_API_BASE_URL", "")
		client, _ := NewClient()

		if client.BaseURL != expected {
			t.Errorf("Expected vs Got BaseURL {%s}",
				cmp.Diff(expected, client.BaseURL))
		}
	})

	t.Run("When set in ENV", func(t *testing.T) {
		expected := "http://localhost:3000/api/v0.1"
		os.Setenv("CLOUDMAILIN_API_BASE_URL", expected)
		client, _ := NewClient()

		if client.BaseURL != expected {
			t.Errorf("Expected vs Got BaseURL {%s}",
				cmp.Diff(expected, client.BaseURL))
		}
	})
}

func TestNewClientFromURL(t *testing.T) {
	original := os.Getenv("CLOUDMAILIN_API_BASE_URL")
	defer os.Setenv("CLOUDMAILIN_API_BASE_URL", original)

	t.Run("Not a URL", func(t *testing.T) {
		_, err := NewClientFromURL("http://\\example.com/foo")

		if err == nil {
			t.Error("Expected error got nil")
		}
	})

	t.Run("Empty URL", func(t *testing.T) {
		_, err := NewClientFromURL("")

		expected := "missing client values"
		if err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error got {%v}", err)
		}
	})

	t.Run("Missing credentials", func(t *testing.T) {
		_, err := NewClientFromURL("http://example.com/foo")

		expected := "missing client values"
		if err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error got {%v}", err)
		}
	})

	t.Run("Valid URL with credentials", func(t *testing.T) {
		client, err := NewClientFromURL("smtp://user:pass@localhost/path")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if client.SMTPToken != "pass" {
			t.Errorf("Expected SMTPToken to be 'pass', got %q", client.SMTPToken)
		}

		if client.SMTPAccountID != "user" {
			t.Errorf("Expected SMTPAccountID to be 'user', got %q", client.SMTPAccountID)
		}
	})

	t.Run("Uses default BaseURL when env not set", func(t *testing.T) {
		os.Setenv("CLOUDMAILIN_API_BASE_URL", "")
		client, err := NewClientFromURL("smtp://user:pass@localhost/path")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "https://api.cloudmailin.com/api/v0.1"
		if client.BaseURL != expected {
			t.Errorf("Expected BaseURL to be %q, got %q", expected, client.BaseURL)
		}
	})

	t.Run("Uses BaseURL from env when set", func(t *testing.T) {
		expected := "http://localhost:3000/api/v0.1"
		os.Setenv("CLOUDMAILIN_API_BASE_URL", expected)
		client, err := NewClientFromURL("smtp://user:pass@localhost/path")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if client.BaseURL != expected {
			t.Errorf("Expected BaseURL to be %q, got %q", expected, client.BaseURL)
		}
	})
}
