package runner

import (
	"fmt"
	"os"
	"sync"

	"github.com/ibrahimsql/bypass403/pkg/bypass"
	"github.com/ibrahimsql/bypass403/pkg/config"
	"github.com/ibrahimsql/bypass403/pkg/http"
	"github.com/ibrahimsql/bypass403/pkg/output"
	"github.com/ibrahimsql/bypass403/pkg/useragent"
	"github.com/ibrahimsql/bypass403/pkg/utils"
)

// Runner encapsulates the bypass403 execution
type Runner struct {
	config *config.Config
	client *http.Client
}

// New creates a new Runner instance
func New(cfg *config.Config) *Runner {
	return &Runner{
		config: cfg,
	}
}

// Run executes the bypass techniques
func (r *Runner) Run() {
	// Initialize HTTP client
	r.client = http.NewClient(r.config.Timeout, r.config.UserAgent)

	// Handle random user agent if enabled
	if r.config.RandomUserAgent {
		if r.config.UserAgentType != "" {
			r.config.UserAgent = useragent.GetRandomByCategory(r.config.UserAgentType)
		} else {
			r.config.UserAgent = useragent.GetRandom()
		}

		if r.config.Verbose {
			fmt.Printf("Using random User-Agent: %s\n", r.config.UserAgent)
		}
	}

	// Verify the URL returns 403
	if err := http.VerifyURL(r.config.URL, r.client); err != nil {
		fmt.Printf("Warning: %s. Continue anyway? (y/n): ", err)
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" && answer != "Y" {
			os.Exit(0)
		}
	}

	// Initialize bypass configuration
	bypassConfig := bypass.Config{
		URL:          r.config.URL,
		UserAgent:    r.config.UserAgent,
		WordlistPath: r.config.WordlistPath,
		Verbose:      r.config.Verbose,
		RandomUA:     r.config.RandomUserAgent,
	}

	fmt.Printf("Starting 403 bypass attempts on %s\n", r.config.URL)
	fmt.Println("============================================")

	// Setup concurrency handling
	var wg sync.WaitGroup
	resultChan := make(chan bypass.Result)
	semaphore := make(chan struct{}, r.config.Threads)

	// Process results in background
	var successfulResults []bypass.Result
	go func() {
		for result := range resultChan {
			if result.StatusCode != 403 && result.StatusCode != 0 && result.StatusCode != 404 {
				fmt.Printf("[+] BYPASS FOUND! %s (%d) - Technique: %s/%s\n",
					result.URL, result.StatusCode, result.Technique, result.Method)
				successfulResults = append(successfulResults, result)

				// Save successful bypass to separate file
				err := utils.SaveForbiddenBypass(result.URL)
				if err != nil && r.config.Verbose {
					fmt.Printf("Warning: Could not save bypass to file: %s\n", err)
				}
			} else if r.config.Verbose {
				fmt.Printf("[-] Failed: %s (%d) - Technique: %s/%s\n",
					result.URL, result.StatusCode, result.Technique, result.Method)
			}
		}
	}()

	// Run selected techniques
	techniques := bypass.GetTechniques()
	for _, technique := range techniques {
		if r.shouldRunTechnique(technique.Category) {
			wg.Add(1)
			semaphore <- struct{}{} // Acquire semaphore

			go func(t bypass.Technique) {
				defer wg.Done()
				defer func() { <-semaphore }() // Release semaphore

				if r.config.Verbose {
					fmt.Printf("Trying %s techniques...\n", t.Name)
				}

				// Use a new random user agent for each technique if enabled
				if r.config.RandomUserAgent {
					if r.config.UserAgentType != "" {
						bypassConfig.UserAgent = useragent.GetRandomByCategory(r.config.UserAgentType)
					} else {
						bypassConfig.UserAgent = useragent.GetRandom()
					}
				}

				// Use the standard http.Client from our custom client
				results, err := t.Test(r.config.URL, r.client.Client, bypassConfig)
				if err != nil && r.config.Verbose {
					fmt.Printf("Error with %s technique: %s\n", t.Name, err)
				}

				for _, result := range results {
					resultChan <- result
				}
			}(technique)
		}
	}

	// Wait for all tests to complete
	wg.Wait()
	close(resultChan)

	// Show summary
	r.showSummary(successfulResults)

	// Generate Burp Suite project if requested
	if r.config.BurpOutput != "" && len(successfulResults) > 0 {
		if err := output.GenerateBurpSuiteProject(successfulResults, r.config.BurpOutput); err != nil {
			fmt.Printf("Error generating Burp Suite project: %s\n", err)
		} else {
			fmt.Printf("Burp Suite project saved to %s\n", r.config.BurpOutput)
		}
	}
}

// shouldRunTechnique determines if a technique should be run based on the configuration
func (r *Runner) shouldRunTechnique(techniqueCategory string) bool {
	if r.config.AllTechniques {
		return true
	}

	if r.config.Category == "" {
		return true
	}

	return utils.ContainsCategory(techniqueCategory, r.config.Category)
}

// showSummary displays a summary of the results
func (r *Runner) showSummary(results []bypass.Result) {
	fmt.Println("\n============= RESULTS =============")
	if len(results) > 0 {
		fmt.Printf("Found %d potential bypasses:\n", len(results))
		for i, result := range results {
			fmt.Printf("%d. %s (%d) - Technique: %s/%s\n",
				i+1, result.URL, result.StatusCode, result.Technique, result.Method)
		}

		// Save results to file if requested
		if r.config.OutputFile != "" {
			utils.SaveResultsToFile(results, r.config.OutputFile, r.config.URL)
		}

		fmt.Println("\nSuccessful bypasses have been saved to forbidden_bypass.txt")
		fmt.Println("\nTips:")
		fmt.Println("* Try combined techniques for better results")
		fmt.Println("* Check each bypass manually to confirm access")
		fmt.Println("* Different status codes may indicate different levels of access")
		fmt.Println("* Consider using a custom wordlist with `-w` option")

		// Show example curl command for the first successful bypass
		if len(results) > 0 {
			fmt.Println("\nExample curl command for first successful bypass:")
			fmt.Println(utils.GenerateCurlCommand(results[0]))

			fmt.Println("\nExample Python request for first successful bypass:")
			fmt.Println(utils.GeneratePythonRequest(results[0]))
		}
	} else {
		fmt.Println("No bypasses found for the given URL.")
		fmt.Println("Try with different techniques or check if the protection can be bypassed.")
		fmt.Println("Consider using a custom wordlist with `-w` option or try the combined techniques category.")
	}
}
