package bypass

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ibrahimsql/bypass403/pkg/wordlist"
)

// TestCombinedBypass tests combined techniques for bypassing 403 responses
func TestCombinedBypass(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Load wordlist for path manipulations
	payloads, err := wordlist.Load(config.WordlistPath)
	if err != nil {
		if config.Verbose {
			payloads = wordlist.GetDefaultPayloads()
		} else {
			// Use a minimal set of payloads
			payloads = []string{
				"/",
				"//",
				"/./",
				"/%2e/",
				"/%20",
				"/..;/",
			}
		}
	}

	// Headers to try
	headers := []map[string]string{
		{"X-Forwarded-For": "127.0.0.1"},
		{"X-Custom-IP-Authorization": "127.0.0.1"},
		{"X-Original-URL": "/admin"},
		{"X-Rewrite-URL": "/admin"},
		{"X-Forwarded-Host": "127.0.0.1"},
		{"X-Host": "127.0.0.1"},
		{"X-Remote-IP": "127.0.0.1"},
		{"User-Agent": "Googlebot/2.1 (+http://www.google.com/bot.html)"},
	}

	// Methods to try
	methods := []string{"GET", "POST", "HEAD", "OPTIONS"}

	// Limit number of paths to avoid excessive requests
	maxPaths := 10
	if len(payloads) > maxPaths {
		payloads = payloads[:maxPaths]
	}

	// Combined tests
	for _, payload := range payloads {
		for _, header := range headers {
			for _, method := range methods {
				manipulatedURL := *parsedURL
				// Apply path manipulation
				manipulatedURL.Path = filepath.Join(filepath.Dir(manipulatedURL.Path), payload)

				req, err := http.NewRequest(method, manipulatedURL.String(), nil)
				if err != nil {
					continue
				}

				// Set the user agent, unless we're specifically testing a different user agent
				if _, ok := header["User-Agent"]; !ok {
					req.Header.Set("User-Agent", config.UserAgent)
				}

				// Apply headers
				for key, value := range header {
					req.Header.Set(key, value)
				}

				resp, err := client.Do(req)
				if err != nil {
					continue
				}
				resp.Body.Close()

				headerNames := make([]string, 0)
				for key := range header {
					headerNames = append(headerNames, key)
				}

				results = append(results, Result{
					URL:        manipulatedURL.String(),
					StatusCode: resp.StatusCode,
					Method:     method,
					Technique:  "Combined: " + strings.Join(headerNames, "+") + " + " + payload,
				})

				// If we found a successful bypass, try adding query parameters
				if resp.StatusCode != 403 && resp.StatusCode != 404 {
					queryManipulations := []string{
						"?id=1",
						"?admin=true",
						"?debug=true",
					}

					for _, query := range queryManipulations {
						queryURL := manipulatedURL
						queryURL.RawQuery = query[1:]

						queryReq, err := http.NewRequest(method, queryURL.String(), nil)
						if err != nil {
							continue
						}

						// Set headers
						if _, ok := header["User-Agent"]; !ok {
							queryReq.Header.Set("User-Agent", config.UserAgent)
						}

						for key, value := range header {
							queryReq.Header.Set(key, value)
						}

						queryResp, err := client.Do(queryReq)
						if err != nil {
							continue
						}
						queryResp.Body.Close()

						results = append(results, Result{
							URL:        queryURL.String(),
							StatusCode: queryResp.StatusCode,
							Method:     method,
							Technique:  "Combined: " + strings.Join(headerNames, "+") + " + " + payload + " + " + query,
						})
					}
				}
			}
		}
	}

	return results, nil
}
