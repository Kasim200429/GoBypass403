package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ibrahimsql/bypass403/pkg/bypass"
)

// GenerateBurpSuiteProject creates a Burp Suite project file for successful bypasses
func GenerateBurpSuiteProject(results []bypass.Result, filename string) error {
	if !strings.HasSuffix(filename, ".burp") {
		filename += ".burp"
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating Burp Suite project file: %s", err)
	}
	defer file.Close()

	// Write Burp Suite XML format
	file.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	file.WriteString("<items burpVersion=\"2023.1.2\" exportTime=\"" + time.Now().Format(time.RFC3339) + "\">\n")

	for _, result := range results {
		if result.StatusCode != 403 && result.StatusCode != 404 {
			// Only include successful bypasses
			file.WriteString(generateBurpItem(result))
		}
	}

	file.WriteString("</items>\n")
	return nil
}

// generateBurpItem creates a Burp Suite item for a bypass result
func generateBurpItem(result bypass.Result) string {
	parsedURL := parseURL(result.URL)

	var burpItem strings.Builder

	burpItem.WriteString("  <item>\n")
	burpItem.WriteString("    <time>" + time.Now().Format(time.RFC3339) + "</time>\n")
	burpItem.WriteString("    <url>" + escapeXML(result.URL) + "</url>\n")
	burpItem.WriteString("    <host>" + escapeXML(parsedURL.Host) + "</host>\n")
	burpItem.WriteString("    <port>" + getPort(parsedURL) + "</port>\n")
	burpItem.WriteString("    <protocol>" + escapeXML(parsedURL.Scheme) + "</protocol>\n")
	burpItem.WriteString("    <method>" + escapeXML(result.Method) + "</method>\n")
	burpItem.WriteString("    <path>" + escapeXML(parsedURL.Path) + "</path>\n")

	// Generate request
	request := generateRequest(result)
	burpItem.WriteString("    <request base64=\"false\">" + escapeXML(request) + "</request>\n")

	// Generate response placeholder
	burpItem.WriteString("    <response base64=\"false\">HTTP/1.1 " + fmt.Sprintf("%d", result.StatusCode) + " " + getStatusText(result.StatusCode) + "\r\n\r\n</response>\n")

	burpItem.WriteString("    <comment>" + escapeXML("403 Bypass: "+result.Technique) + "</comment>\n")
	burpItem.WriteString("    <highlight>green</highlight>\n")
	burpItem.WriteString("    <tags>\n")
	burpItem.WriteString("      <tag>403 Bypass</tag>\n")
	burpItem.WriteString("      <tag>" + escapeXML(result.Technique) + "</tag>\n")
	burpItem.WriteString("    </tags>\n")
	burpItem.WriteString("  </item>\n")

	return burpItem.String()
}

// generateRequest creates an HTTP request string for a bypass result
func generateRequest(result bypass.Result) string {
	parsedURL := parseURL(result.URL)

	var request strings.Builder

	// Request line
	path := parsedURL.Path
	if parsedURL.RawQuery != "" {
		path += "?" + parsedURL.RawQuery
	}
	if path == "" {
		path = "/"
	}
	request.WriteString(result.Method + " " + path + " HTTP/1.1\r\n")

	// Headers
	request.WriteString("Host: " + parsedURL.Host + "\r\n")

	// Add technique-specific headers
	if strings.Contains(result.Technique, "Header") {
		headerName := strings.TrimPrefix(strings.Split(result.Technique, ":")[1], " ")
		if headerName == "X-Original-URL" {
			request.WriteString("X-Original-URL: /\r\n")
		} else if headerName == "X-Rewrite-URL" {
			request.WriteString("X-Rewrite-URL: /\r\n")
		} else if strings.Contains(headerName, "Forwarded-For") {
			request.WriteString(headerName + ": 127.0.0.1\r\n")
		}
	}

	// Common headers
	request.WriteString("User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36\r\n")
	request.WriteString("Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\n")
	request.WriteString("Accept-Language: en-US,en;q=0.5\r\n")
	request.WriteString("Connection: close\r\n")
	request.WriteString("\r\n")

	return request.String()
}

// parseURL parses a URL string and returns its components
func parseURL(urlStr string) struct {
	Scheme   string
	Host     string
	Path     string
	RawQuery string
} {
	// Simple URL parsing
	scheme := "http"
	if strings.HasPrefix(urlStr, "https://") {
		scheme = "https"
		urlStr = strings.TrimPrefix(urlStr, "https://")
	} else if strings.HasPrefix(urlStr, "http://") {
		urlStr = strings.TrimPrefix(urlStr, "http://")
	}

	host := urlStr
	path := "/"
	rawQuery := ""

	if idx := strings.Index(urlStr, "/"); idx != -1 {
		host = urlStr[:idx]
		path = urlStr[idx:]
	}

	if idx := strings.Index(path, "?"); idx != -1 {
		rawQuery = path[idx+1:]
		path = path[:idx]
	}

	return struct {
		Scheme   string
		Host     string
		Path     string
		RawQuery string
	}{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: rawQuery,
	}
}

// getPort returns the port based on the URL scheme
func getPort(parsedURL struct {
	Scheme   string
	Host     string
	Path     string
	RawQuery string
}) string {
	if strings.Contains(parsedURL.Host, ":") {
		parts := strings.Split(parsedURL.Host, ":")
		return parts[1]
	}

	if parsedURL.Scheme == "https" {
		return "443"
	}
	return "80"
}

// getStatusText returns the text description for an HTTP status code
func getStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 301:
		return "Moved Permanently"
	case 302:
		return "Found"
	case 307:
		return "Temporary Redirect"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}

// escapeXML escapes special characters in XML
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
