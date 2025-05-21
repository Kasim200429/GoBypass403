# Installation Guide

This document provides instructions for installing and configuring GoBypass403.

## Prerequisites

Before installing GoBypass403, ensure your system meets the [system requirements](./System-Requirements.md).

## Installation Methods

### Method 1: Direct Binary Installation

#### From Releases

```bash
# Download latest release (Linux/macOS)
curl -sSL https://github.com/ibrahimsql/GoBypass403/releases/latest/download/gobypass403_$(uname -s)_$(uname -m).tar.gz | tar xz

# Move binary to executable path
sudo mv gobypass403 /usr/local/bin/
chmod +x /usr/local/bin/gobypass403
```

#### Windows Installation

1. Download the latest Windows binary from [GitHub Releases](https://github.com/ibrahimsql/GoBypass403/releases)
2. Extract the ZIP archive
3. Add the binary location to your PATH environment variable or execute directly

### Method 2: Go Installation

This method requires Go 1.21 or higher installed on your system.

```bash
# Install latest release version
go install github.com/ibrahimsql/GoBypass403@latest
```

The binary will be installed in your `$GOPATH/bin` directory.

### Method 3: Build from Source

This method provides the most control over the installation process.

```bash
# Clone repository
git clone https://github.com/ibrahimsql/GoBypass403.git
cd GoBypass403

# Install dependencies
go mod download

# Build binary
go build -o gobypass403 .

# Optional: Install system-wide
sudo mv gobypass403 /usr/local/bin/
```

## Basic Configuration

GoBypass403 is designed to work with minimal configuration. The essential configuration is specified through command-line parameters.

### Environment Variables

GoBypass403 currently supports the following environment variables:

| Variable | Purpose | Default Value |
|----------|---------|---------------|
| HTTP_PROXY / HTTPS_PROXY | Standard proxy environment variables | none |
| NO_PROXY | Skip proxy for specified hosts | none |

## Verification

To verify your installation:

```bash
# Check if the binary is accessible
gobypass403 --version
```

## Uninstallation

To remove GoBypass403 from your system:

```bash
# If installed via go install
rm $(which gobypass403)

# If installed manually
rm /usr/local/bin/gobypass403
```

## Next Steps

After successful installation, proceed to:

1. [Basic Usage](./CLI-Reference.md) - Learn command-line options
2. [Bypass Techniques](./Bypass-Techniques.md) - Learn about the bypass techniques

## Coming in Future Releases

> Note: The following features are planned for upcoming releases and are not yet implemented.

### Upcoming Installation Methods
- Docker container support
- Package manager installations (apt, brew, etc.)
- Cross-compilation for various platforms

### Upcoming Configuration Features
- YAML configuration file support
- Custom payload directory structure
- Shell completion scripts
- Self-update mechanism
- System-wide and user-specific configuration 