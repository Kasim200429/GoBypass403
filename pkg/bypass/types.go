package bypass

import (
	"net/http"
)

// Result represents the result of a bypass attempt
type Result struct {
	URL        string
	StatusCode int
	Method     string
	Technique  string
}

// Config represents configuration options for bypass techniques
type Config struct {
	URL          string
	UserAgent    string
	WordlistPath string
	Verbose      bool
	RandomUA     bool
}

// Technique represents a bypass technique
type Technique struct {
	Name     string
	Test     func(string, *http.Client, Config) ([]Result, error)
	Category string
}

// GetTechniques returns all available bypass techniques
func GetTechniques() []Technique {
	return []Technique{
		{"Method Manipulation", TestMethodManipulation, "Request Method"},
		{"URL Path Manipulation", TestURLPathManipulation, "URL Path"},
		{"Header Manipulation", TestHeaderManipulation, "Headers"},
		{"IP Spoofing Headers", TestIPSpoofingHeaders, "IP Spoofing"},
		{"URL Encoding Bypass", TestURLEncodingBypass, "URL Encoding"},
		{"Protocol Bypass", TestProtocolBypass, "Protocol"},
		{"Path Traversal", TestPathTraversal, "Path Traversal"},
		{"Caching Proxy Bypass", TestCachingProxyBypass, "Proxy"},
		{"Specialized Payloads", TestPayloads, "Specialized"},
		{"Wordlist Path Bypass", TestWordlistPathBypass, "Wordlist"},
		{"Combined Technique Bypass", TestCombinedBypass, "Combined"},
	}
}
