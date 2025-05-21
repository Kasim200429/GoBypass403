# Command Line Interface Reference

This document provides a comprehensive reference for the GoBypass403 command line interface, including all available commands, options, and usage patterns.

## Command Syntax

The general syntax for GoBypass403 commands follows this pattern:

```
gobypass403 [global options] command [command options] [arguments]
```

## Global Options

These options apply to all commands and control the general behavior of the tool.

| Option | Format | Description | Default |
|--------|--------|-------------|---------|
| `-u`, `--url` | `<URL>` | Target URL that returns 403 Forbidden | None (Required) |
| `-t`, `--threads` | `<int>` | Number of concurrent execution threads | 10 |
| `-o`, `--output` | `<file>` | Output file path for results | None (stdout only) |
| `-timeout` | `<int>` | HTTP request timeout in seconds | 10 |
| `-v`, `--verbose` | | Enable verbose output mode | false |
| `--version` | | Display version information and exit | |
| `-h`, `--help` | | Display help information and exit | |

## Core Commands

### Default Command (No Command Specified)

When executed without a specific command, GoBypass403 runs with the default bypass operation.

```bash
gobypass403 -u https://example.com/admin -v
```

### Test Command

```bash
gobypass403 test -u https://example.com/admin -c Headers
```

The test command executes bypass attempts with specific techniques.

| Option | Format | Description | Default |
|--------|--------|-------------|---------|
| `-c`, `--category` | `<string>` | Category of bypass techniques to attempt | (All categories) |
| `-w`, `--wordlist` | `<path>` | Custom wordlist file path | payloads/bypasses.txt |
| `--all` | | Try all bypass techniques | false |

## Technique Selection Options

GoBypass403 provides options to select specific bypass techniques for testing:

| Option | Description |
|--------|-------------|
| `--method` | Enable method manipulation techniques |
| `--path` | Enable path manipulation techniques |
| `--headers` | Enable header manipulation techniques |
| `--ip` | Enable IP spoofing techniques |
| `--encoding` | Enable URL encoding techniques |
| `--protocol` | Enable protocol switching techniques |
| `--traversal` | Enable path traversal techniques |
| `--proxy` | Enable proxy bypass techniques |
| `--payloads` | Enable specialized payloads |
| `--wordlist` | Enable wordlist-based techniques |
| `--combined` | Enable combined techniques |

## Advanced Options

These options control more specialized aspects of the tool's behavior.

| Option | Format | Description | Default |
|--------|--------|-------------|---------|
| `--user-agent`, `-ua` | `<string>` | User-Agent string to use in requests | Mozilla/5.0... |
| `--random-ua` | | Randomize User-Agent for each request | false |
| `--ua-type` | `<string>` | Type of User-Agent to use (mobile, desktop, etc.) | |
| `--proxy` | `<url>` | Proxy URL (e.g., http://127.0.0.1:8080) | None |
| `--no-verify` | | Disable TLS certificate verification | false |
| `--follow-redirects` | | Follow HTTP redirects | false |
| `--max-redirects` | `<int>` | Maximum number of redirects to follow | 10 |
| `--burp` | `<file>` | Generate Burp Suite project file | None |

## Output Control Options

These options control how results are presented and saved.

| Option | Format | Description | Default |
|--------|--------|-------------|---------|
| `--no-color` | | Disable colored output | false |
| `--json` | | Output results in JSON format | false |
| `--silent` | | Suppress all output except results | false |
| `--show-headers` | | Show response headers in output | false |
| `--show-body` | | Show response body in output | false |
| `--save-responses` | | Save full HTTP responses to disk | false |
| `--response-dir` | `<dir>` | Directory to save responses | ./responses/ |

## Exit Codes

GoBypass403 uses the following exit codes to indicate different execution outcomes:

| Exit Code | Description |
|-----------|-------------|
| 0 | Successful execution, bypasses found |
| 1 | Successful execution, no bypasses found |
| 2 | Command line argument error |
| 3 | Network error or connection issues |
| 4 | File I/O error |
| 5 | Configuration error |
| 6 | Unexpected internal error |

## Environment Variables

GoBypass403 recognizes the following environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| GOBYPASS_CONFIG | Path to configuration file | /path/to/config.yaml |
| GOBYPASS_LOG_LEVEL | Logging verbosity level | debug |
| HTTP_PROXY | HTTP proxy URL | http://proxy:8080 |
| HTTPS_PROXY | HTTPS proxy URL | http://proxy:8080 |
| NO_PROXY | Comma-separated list of hosts to exclude from proxy | localhost,127.0.0.1 |

## Usage Examples

### Basic Scan

```bash
# Basic scan with verbose output
gobypass403 -u https://example.com/admin -v

# Scan with increased thread count
gobypass403 -u https://example.com/admin -t 20

# Scan with output file
gobypass403 -u https://example.com/admin -o results.txt
```

### Technique Selection

```bash
# Try only header manipulation techniques
gobypass403 -u https://example.com/admin -c Headers -v

# Try path and protocol techniques
gobypass403 -u https://example.com/admin -c Path,Protocol -v

# Run all techniques
gobypass403 -u https://example.com/admin --all
```

### Advanced Usage

```bash
# Use custom wordlist with path traversal techniques
gobypass403 -u https://example.com/admin -c Traversal -w custom_paths.txt

# Use random user agent with specific category
gobypass403 -u https://example.com/admin --random-ua --ua-type mobile

# Generate Burp Suite project
gobypass403 -u https://example.com/admin --burp project.burp

# Use proxy and custom timeout
gobypass403 -u https://example.com/admin --proxy http://127.0.0.1:8080 -timeout 30
```

### Output Options

```bash
# Output in JSON format
gobypass403 -u https://example.com/admin --json > results.json

# Silent mode with only successful results
gobypass403 -u https://example.com/admin --silent -o successful.txt

# Save full HTTP responses
gobypass403 -u https://example.com/admin --save-responses --response-dir ./evidence/
```

## Configuration File Reference

GoBypass403 can be configured using a YAML configuration file. See the [Configuration Reference](./Configuration.md) for details.

## See Also

* [Configuration Reference](./Configuration.md) - Detailed configuration options
* [Bypass Techniques](./Bypass-Techniques.md) - Technical details of bypass methods
* [Advanced Usage](./Advanced-Usage.md) - Complex usage scenarios 