package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ibrahimsql/bypass403/pkg/bypass"
)

// SaveForbiddenBypass saves a successful bypass URL to a file
func SaveForbiddenBypass(url string) error {
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

// SaveResultsToFile saves bypass results to a file
func SaveResultsToFile(results []bypass.Result, filename string, targetURL string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("=== 403 Bypass Results ===\n")
	writer.WriteString(fmt.Sprintf("Target URL: %s\n", targetURL))
	writer.WriteString(fmt.Sprintf("Date: %s\n\n", time.Now().Format(time.RFC1123)))

	for i, r := range results {
		writer.WriteString(fmt.Sprintf("%d. %s (%d) - Technique: %s/%s\n",
			i+1, r.URL, r.StatusCode, r.Technique, r.Method))
	}

	writer.WriteString("\n=== Tips ===\n")
	writer.WriteString("1. Check all successful responses manually to confirm they provide actual access\n")
	writer.WriteString("2. Some bypasses might only provide partial access or different content\n")
	writer.WriteString("3. Combine multiple techniques for better results\n")
	writer.WriteString("4. Try advanced mutations and custom wordlists for better coverage\n")

	// Add curl examples
	writer.WriteString("\n=== Examples for successful bypasses ===\n")
	for _, r := range results {
		if r.StatusCode != 403 && r.StatusCode != 404 {
			writer.WriteString(fmt.Sprintf("CURL: %s\n", GenerateCurlCommand(r)))
			writer.WriteString(fmt.Sprintf("Python: %s\n\n", GeneratePythonRequest(r)))
		}
	}

	writer.Flush()

	return nil
}

// ContainsCategory checks if a technique category matches the user-specified category
func ContainsCategory(techniqueCategory, userCategory string) bool {
	return strings.Contains(
		strings.ToLower(techniqueCategory),
		strings.ToLower(userCategory),
	)
}

// GenerateCurlCommand generates a curl command for a successful bypass
func GenerateCurlCommand(result bypass.Result) string {
	// Basic curl command
	curlCmd := fmt.Sprintf("curl -X %s '%s'", result.Method, result.URL)

	// Add additional headers based on technique
	if strings.Contains(result.Technique, "Header") {
		headerName := strings.TrimPrefix(strings.Split(result.Technique, ":")[1], " ")
		if headerName == "X-Original-URL" {
			curlCmd += " -H 'X-Original-URL: /'"
		} else if headerName == "X-Rewrite-URL" {
			curlCmd += " -H 'X-Rewrite-URL: /'"
		} else if strings.Contains(headerName, "Forwarded-For") {
			curlCmd += " -H '" + headerName + ": 127.0.0.1'"
		} else if headerName == "User-Agent" {
			curlCmd += " -H 'User-Agent: Googlebot/2.1 (+http://www.google.com/bot.html)'"
		}
	}

	// Add -k for insecure SSL
	curlCmd += " -k"

	return curlCmd
}

// GeneratePythonRequest generates a Python request code for a successful bypass
func GeneratePythonRequest(result bypass.Result) string {
	// Basic Python code
	pythonCode := "import requests\n\n"

	// Add headers if needed
	if strings.Contains(result.Technique, "Header") {
		headerName := strings.TrimPrefix(strings.Split(result.Technique, ":")[1], " ")
		pythonCode += "headers = {\n"
		if headerName == "X-Original-URL" {
			pythonCode += "    'X-Original-URL': '/',\n"
		} else if headerName == "X-Rewrite-URL" {
			pythonCode += "    'X-Rewrite-URL': '/',\n"
		} else if strings.Contains(headerName, "Forwarded-For") {
			pythonCode += "    '" + headerName + "': '127.0.0.1',\n"
		} else if headerName == "User-Agent" {
			pythonCode += "    'User-Agent': 'Googlebot/2.1 (+http://www.google.com/bot.html)',\n"
		}
		pythonCode += "}\n\n"
		pythonCode += fmt.Sprintf("response = requests.%s('%s', headers=headers, verify=False)\n",
			strings.ToLower(result.Method), result.URL)
	} else {
		pythonCode += fmt.Sprintf("response = requests.%s('%s', verify=False)\n",
			strings.ToLower(result.Method), result.URL)
	}

	pythonCode += "print(response.status_code)\n"
	pythonCode += "print(response.text)\n"

	return pythonCode
}
