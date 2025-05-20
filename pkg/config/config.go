package config

import (
	"errors"
	"net/url"
)

// Config holds all configuration options for bypass403
type Config struct {
	// Required parameters
	URL string

	// Optional parameters
	Threads         int
	OutputFile      string
	Timeout         int
	Verbose         bool
	AllTechniques   bool
	Category        string
	UserAgent       string
	WordlistPath    string
	RandomUserAgent bool
	UserAgentType   string
	BurpOutput      string
	Version         bool
}

// NewDefaultConfig returns a Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		Threads:         10,
		Timeout:         10,
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		WordlistPath:    "payloads/bypasses.txt",
		RandomUserAgent: false,
		AllTechniques:   false,
		Verbose:         false,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Check if URL is provided
	if c.URL == "" && !c.Version {
		return errors.New("URL is required")
	}

	// If URL is provided, validate it
	if c.URL != "" {
		_, err := url.Parse(c.URL)
		if err != nil {
			return errors.New("invalid URL format")
		}
	}

	// Validate threads
	if c.Threads < 1 {
		return errors.New("threads must be at least 1")
	}

	// Validate timeout
	if c.Timeout < 1 {
		return errors.New("timeout must be at least 1 second")
	}

	return nil
}
