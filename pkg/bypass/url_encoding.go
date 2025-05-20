package bypass

import (
	"net/http"
	"net/url"
	"strings"
)

// TestURLEncodingBypass tests URL encoding techniques to bypass 403 responses
func TestURLEncodingBypass(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	originalPath := parsedURL.Path
	encodedPaths := []string{
		strings.ReplaceAll(originalPath, "/", "%2f"),
		strings.ReplaceAll(originalPath, "/", "%252f"),
		strings.ReplaceAll(originalPath, "/", "%2F"),
		encodePath(originalPath, "a", "%61"),
		encodePath(originalPath, "A", "%41"),
		encodePath(originalPath, "s", "%73"),
		encodePath(originalPath, "S", "%53"),
		encodePath(originalPath, "/", "%2f%2f"),
		encodePath(originalPath, ".", "%2e"),
		// Double encoding
		strings.ReplaceAll(originalPath, "/", "%25%32%66"),
		strings.ReplaceAll(originalPath, "/", "%25%32%46"),
		strings.ReplaceAll(originalPath, ".", "%25%32%65"),
		strings.ReplaceAll(originalPath, ".", "%25%32%45"),
		// Triple encoding
		strings.ReplaceAll(originalPath, "/", "%25%25%33%32%25%36%36"),
		// Mixed encoding
		mixEncode(originalPath),
	}

	for _, path := range encodedPaths {
		manipulatedURL := *parsedURL
		manipulatedURL.Path = path

		req, err := http.NewRequest("GET", manipulatedURL.String(), nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", config.UserAgent)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, Result{
			URL:        manipulatedURL.String(),
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "URL Encoding",
		})
	}

	return results, nil
}

// encodePath replaces a specific character with its encoded version
func encodePath(path, char, encoded string) string {
	return strings.ReplaceAll(path, char, encoded)
}

// mixEncode creates a mixed encoding pattern for path elements
func mixEncode(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if i%2 == 0 && part != "" {
			// URL encode even-indexed parts
			parts[i] = url.QueryEscape(part)
		}
	}
	return strings.Join(parts, "/")
}
