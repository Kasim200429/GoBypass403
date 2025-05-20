package mutation

import (
	"fmt"
	"net/url"
	"strings"
)

// Mutator represents a function that mutates a URL path
type Mutator func(string) []string

// GetAllMutators returns all available path mutation functions
func GetAllMutators() map[string]Mutator {
	return map[string]Mutator{
		"Case Manipulation":   CaseManipulation,
		"URL Encoding":        URLEncoding,
		"Double Encoding":     DoubleEncoding,
		"Path Traversal":      PathTraversal,
		"Slash Manipulation":  SlashManipulation,
		"Extension Addition":  ExtensionAddition,
		"Special Characters":  SpecialCharacters,
		"Parameter Injection": ParameterInjection,
	}
}

// CaseManipulation performs case mutations on a path
func CaseManipulation(path string) []string {
	var results []string

	// Original path
	results = append(results, path)

	// Uppercase
	results = append(results, strings.ToUpper(path))

	// Lowercase
	results = append(results, strings.ToLower(path))

	// Mixed case - toggle every other character
	var mixed strings.Builder
	for i, c := range path {
		if i%2 == 0 {
			mixed.WriteString(strings.ToUpper(string(c)))
		} else {
			mixed.WriteString(strings.ToLower(string(c)))
		}
	}
	results = append(results, mixed.String())

	return results
}

// URLEncoding performs URL encoding mutations on a path
func URLEncoding(path string) []string {
	var results []string

	// Original path
	results = append(results, path)

	// Full path URL encoding
	results = append(results, url.PathEscape(path))

	// Encode slashes
	results = append(results, strings.ReplaceAll(path, "/", "%2f"))
	results = append(results, strings.ReplaceAll(path, "/", "%2F"))

	// Encode dots
	results = append(results, strings.ReplaceAll(path, ".", "%2e"))
	results = append(results, strings.ReplaceAll(path, ".", "%2E"))

	return results
}

// DoubleEncoding performs double URL encoding mutations on a path
func DoubleEncoding(path string) []string {
	var results []string

	// Original path
	results = append(results, path)

	// Double encode slashes
	results = append(results, strings.ReplaceAll(path, "/", "%252f"))
	results = append(results, strings.ReplaceAll(path, "/", "%252F"))

	// Double encode dots
	results = append(results, strings.ReplaceAll(path, ".", "%252e"))
	results = append(results, strings.ReplaceAll(path, ".", "%252E"))

	return results
}

// PathTraversal performs path traversal mutations on a path
func PathTraversal(path string) []string {
	var results []string

	// Original path
	results = append(results, path)

	// Add path traversal sequences
	results = append(results, path+"/..")
	results = append(results, path+"/../")
	results = append(results, path+"/.././")
	results = append(results, path+"/../../")
	results = append(results, path+"/../../../")

	// URL encoded traversal
	results = append(results, path+"/%2e%2e")
	results = append(results, path+"/%2e%2e/")
	results = append(results, path+"/%2e%2e%2f")

	// Double encoded traversal
	results = append(results, path+"/%252e%252e")
	results = append(results, path+"/%252e%252e/")
	results = append(results, path+"/%252e%252e%252f")

	// Null byte injection
	results = append(results, path+"/../%00")
	results = append(results, path+"/.%00./")

	return results
}

// SlashManipulation performs slash manipulation mutations on a path
func SlashManipulation(path string) []string {
	var results []string

	// Original path
	results = append(results, path)

	// Add extra slashes
	results = append(results, strings.ReplaceAll(path, "/", "//"))
	results = append(results, strings.ReplaceAll(path, "/", "///"))

	// Add trailing slash
	if !strings.HasSuffix(path, "/") {
		results = append(results, path+"/")
	}

	// Backslash instead of forward slash
	results = append(results, strings.ReplaceAll(path, "/", "\\"))
	results = append(results, strings.ReplaceAll(path, "/", "\\/"))
	results = append(results, strings.ReplaceAll(path, "/", "/\\"))

	return results
}

// ExtensionAddition adds various file extensions to a path
func ExtensionAddition(path string) []string {
	var results []string

	// Original path
	results = append(results, path)

	// Common extensions
	extensions := []string{".html", ".php", ".asp", ".aspx", ".jsp", ".json", ".xml", ".txt", ".bak", ".old", ".swp", "~"}

	for _, ext := range extensions {
		results = append(results, path+ext)
	}

	// Null byte + extension
	results = append(results, path+"%00.html")
	results = append(results, path+"%00.php")
	results = append(results, path+"%00.asp")

	return results
}

// SpecialCharacters adds special characters to a path
func SpecialCharacters(path string) []string {
	var results []string

	// Original path
	results = append(results, path)

	// Special characters
	specialChars := []string{"%09", "%0a", "%0d", "%20", "%23", "%25", "%26", "%2b", "%3f", "%5c", ";", ":", "!", "$", "^", "*"}

	for _, char := range specialChars {
		// Append special character
		results = append(results, path+char)

		// Insert special character after each slash
		parts := strings.Split(path, "/")
		var modified string
		for i, part := range parts {
			if i < len(parts)-1 {
				modified += part + "/" + char
			} else {
				modified += part
			}
		}
		results = append(results, modified)
	}

	return results
}

// ParameterInjection adds various URL parameters to a path
func ParameterInjection(path string) []string {
	var results []string

	// Original path
	results = append(results, path)

	// Common parameters
	params := []string{
		"?id=1",
		"?page=1",
		"?file=index",
		"?include=true",
		"?debug=true",
		"?test=1",
		"?admin=1",
		"?admin=true",
		"?access=1",
		"?access=true",
		"?show=1",
		"?s=1",
		"?p=1",
	}

	for _, param := range params {
		results = append(results, path+param)
	}

	return results
}

// MutateURL applies all mutation techniques to a URL
func MutateURL(urlStr string) ([]string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %s", err)
	}

	originalPath := parsedURL.Path
	var results []string

	// Get all mutators
	mutators := GetAllMutators()

	// Apply each mutator to the path
	for _, mutator := range mutators {
		mutatedPaths := mutator(originalPath)

		// Create new URLs with the mutated paths
		for _, path := range mutatedPaths {
			mutatedURL := *parsedURL
			mutatedURL.Path = path
			results = append(results, mutatedURL.String())
		}
	}

	return results, nil
}
