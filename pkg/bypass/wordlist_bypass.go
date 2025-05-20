package bypass

import (
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/ibrahimsql/bypass403/pkg/wordlist"
)

// TestWordlistPathBypass tests bypass paths from a wordlist
func TestWordlistPathBypass(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Load wordlist
	payloads, err := wordlist.Load(config.WordlistPath)
	if err != nil {
		if config.Verbose {
			// Fall back to default payloads
			payloads = wordlist.GetDefaultPayloads()
		} else {
			// Return empty results but don't break execution
			return results, nil
		}
	}

	baseDir := filepath.Dir(parsedURL.Path)
	if baseDir == "." || baseDir == "/" {
		baseDir = ""
	}

	for _, payload := range payloads {
		// Try with base path + payload
		manipulatedURL := *parsedURL
		manipulatedURL.Path = filepath.Join(baseDir, payload)

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
			Technique:  "Wordlist Path: " + payload,
		})

		// If we found a successful bypass and it's not a 403 or 404 response, try with POST too
		if resp.StatusCode != 403 && resp.StatusCode != 404 {
			postReq, err := http.NewRequest("POST", manipulatedURL.String(), nil)
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
				URL:        manipulatedURL.String(),
				StatusCode: postResp.StatusCode,
				Method:     "POST",
				Technique:  "Wordlist Path: " + payload,
			})
		}

		// Try also adding query parameters and fragments
		queryManipulations := []string{
			"?id=1",
			"?admin=true",
			"?debug=true",
			"?access=true",
			"?token=1",
			"?_=" + payload,
		}

		for _, queryParam := range queryManipulations {
			queryURL := manipulatedURL
			queryURL.RawQuery = queryParam[1:]

			queryReq, err := http.NewRequest("GET", queryURL.String(), nil)
			if err != nil {
				continue
			}

			queryReq.Header.Set("User-Agent", config.UserAgent)

			queryResp, err := client.Do(queryReq)
			if err != nil {
				continue
			}
			queryResp.Body.Close()

			results = append(results, Result{
				URL:        queryURL.String(),
				StatusCode: queryResp.StatusCode,
				Method:     "GET",
				Technique:  "Wordlist Path + Query: " + payload + queryParam,
			})
		}
	}

	return results, nil
}
