package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ibrahimsql/bypass403/pkg/config"
	"github.com/ibrahimsql/bypass403/pkg/runner"
	"github.com/ibrahimsql/bypass403/pkg/utils"
)

func main() {
	// Print the banner
	utils.PrintBanner()

	// Parse command line flags
	cfg := parseFlags()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Error: %s\n", err)
		printUsage()
		os.Exit(1)
	}

	// Start the bypass runner
	r := runner.New(cfg)
	r.Run()
}

func parseFlags() *config.Config {
	cfg := config.NewDefaultConfig()

	flag.StringVar(&cfg.URL, "u", "", "URL that returns 403 Forbidden")
	flag.IntVar(&cfg.Threads, "t", 10, "Number of concurrent threads")
	flag.StringVar(&cfg.OutputFile, "o", "", "Output file to save results")
	flag.IntVar(&cfg.Timeout, "timeout", 10, "HTTP request timeout in seconds")
	flag.BoolVar(&cfg.Verbose, "v", false, "Verbose mode")
	flag.BoolVar(&cfg.AllTechniques, "all", false, "Try all bypass techniques")
	flag.StringVar(&cfg.Category, "c", "", "Category of bypass techniques to try (Method, Path, Headers, IP, Encoding, Protocol, Traversal, Proxy, Advanced)")
	flag.StringVar(&cfg.UserAgent, "ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36", "User-Agent to use")
	flag.StringVar(&cfg.WordlistPath, "w", "payloads/bypasses.txt", "Path to wordlist file for bypass attempts")
	flag.BoolVar(&cfg.Version, "version", false, "Print version information and exit")

	flag.Parse()

	// If version flag is set, print version info and exit
	if cfg.Version {
		utils.PrintInfo()
		os.Exit(0)
	}

	return cfg
}

func printUsage() {
	fmt.Println("403 Bypass - A tool to bypass 403 Forbidden responses")
	fmt.Println("Usage: bypass403 -u https://example.com/forbidden")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  bypass403 -u https://example.com/admin -v -o results.txt")
	fmt.Println("  bypass403 -u https://example.com/admin -w payloads/bypasses.txt -all")
	fmt.Println("\nNote: Successful bypasses are automatically saved to forbidden_bypass.txt")
}
