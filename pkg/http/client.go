package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client is a custom HTTP client for bypass403
type Client struct {
	*http.Client
	UserAgent string
}

// NewClient creates a new HTTP client with custom settings
func NewClient(timeout int, userAgent string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	return &Client{
		Client:    client,
		UserAgent: userAgent,
	}
}

// VerifyURL checks if the URL returns a 403 Forbidden response
func VerifyURL(urlStr string, client *Client) error {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", client.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 403 {
		return errors.New(fmt.Sprintf("The provided URL returns %d, not 403 Forbidden", resp.StatusCode))
	}

	return nil
}
