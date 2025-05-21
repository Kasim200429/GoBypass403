package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

// The default banner as a string in case the banner file cannot be loaded
const defaultBanner = `███████████████████████████████████████████████████████████████████████████
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
 ║           github.com/ibrahimsql/bypass403             ║
 ╚═══════════════════════════════════════════════════════╝`

// PrintBanner loads and prints the banner from file or uses the default banner
func PrintBanner() {
	bannerText := loadBannerFile()

	// Print the banner in color
	color.Cyan(bannerText)
	fmt.Println()
}

// loadBannerFile attempts to load the banner from the banner.txt file
// If it can't find or read the file, it falls back to the default banner
func loadBannerFile() string {
	// Try to find the banner file in various locations
	possiblePaths := []string{
		"banner.txt",
		"./banner.txt",
		"../banner.txt",
		filepath.Join(getExecutablePath(), "banner.txt"),
	}

	for _, path := range possiblePaths {
		content, err := os.ReadFile(path)
		if err == nil {
			return string(content)
		}
	}

	// Fall back to the default banner
	return defaultBanner
}

// getExecutablePath returns the directory of the current executable
func getExecutablePath() string {
	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}

	// For symlinks, get the real path
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return filepath.Dir(exePath)
	}

	return filepath.Dir(realPath)
}

// GetVersion returns the current version of the tool
func GetVersion() string {
	return "v1.0.0"
}

// PrintInfo prints information about the tool
func PrintInfo() {
	fmt.Printf("bypass403 %s\n", GetVersion())
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Author: Ibrahim SQL")
	fmt.Println("GitHub: https://github.com/ibrahimsql/bypass403")
	fmt.Println()
}
