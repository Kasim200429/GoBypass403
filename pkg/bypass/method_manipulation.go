package bypass

import (
	"net/http"
)

// TestMethodManipulation tests different HTTP methods to bypass 403 responses
func TestMethodManipulation(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result
	methods := []string{
		"GET", "POST", "HEAD", "OPTIONS", "PUT", "DELETE", "TRACE", "CONNECT", "PATCH",
		"PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK", "FAKE-METHOD",
		"REPORT", "CHECKOUT", "CHECKIN", "SEARCH", "SUBSCRIBE", "UNSUBSCRIBE", "NOTIFY",
		"BREW", "BASELINE-CONTROL", "ACL", "VERSION-CONTROL", "MKWORKSPACE", "UPDATE",
		"LABEL", "MERGE", "PURGE", "DEBUG", "FOO", "BAR", "BATMAN", "ADMIN",
	}

	for _, method := range methods {
		req, err := http.NewRequest(method, baseURL, nil)
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
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     method,
			Technique:  "Method Manipulation",
		})
	}

	return results, nil
}
