package bypass

import (
	"net/http"
	"net/url"
	"strings"
)

// TestProtocolBypass tests different protocol manipulations to bypass 403 responses
func TestProtocolBypass(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Try different protocols and protocol-related manipulations
	protocols := []string{
		"http://",
		"https://",
		"http:\\\\",
		"https:\\\\",
		"ftp://",
		"ftps://",
		"gopher://",
		"file://",
		"http:/",
		"https:/",
		"//",
	}

	for _, protocol := range protocols {
		// Skip the current protocol
		if strings.HasPrefix(baseURL, protocol) && protocol != "//" {
			continue
		}

		var manipulatedURL string

		if protocol == "//" {
			// Protocol-relative URL
			manipulatedURL = "//" + parsedURL.Host + parsedURL.Path
			if parsedURL.RawQuery != "" {
				manipulatedURL += "?" + parsedURL.RawQuery
			}
		} else {
			manipulatedURL = protocol + parsedURL.Host + parsedURL.Path
			if parsedURL.RawQuery != "" {
				manipulatedURL += "?" + parsedURL.RawQuery
			}
		}

		req, err := http.NewRequest("GET", manipulatedURL, nil)
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
			URL:        manipulatedURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Protocol Change: " + protocol,
		})

		// Also try with POST method
		postReq, err := http.NewRequest("POST", manipulatedURL, nil)
		if err != nil {
			continue
		}

		postReq.Header.Set("User-Agent", config.UserAgent)
		postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		postResp, err := client.Do(postReq)
		if err != nil {
			continue
		}
		postResp.Body.Close()

		results = append(results, Result{
			URL:        manipulatedURL,
			StatusCode: postResp.StatusCode,
			Method:     "POST",
			Technique:  "Protocol Change: " + protocol,
		})
	}

	return results, nil
}
