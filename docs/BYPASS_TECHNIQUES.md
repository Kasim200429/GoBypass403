# Bypass Techniques

This document details the various bypass techniques implemented in GoBypass403.

## Method Manipulation
Tests different HTTP methods to bypass access controls:
- GET, POST, HEAD, OPTIONS, PUT, DELETE
- Method override techniques
- Custom method testing

## Path Manipulation
Modifies URL paths using various techniques:
- Path normalization
- URL encoding variations
- Path traversal patterns
- Directory listing attempts

## Header Manipulation
Tests HTTP headers that might affect access control:
- X-Forwarded-For
- X-Original-URL
- X-Rewrite-URL
- Referer
- User-Agent
- Custom header testing

## IP Spoofing
Tests IP-based access controls:
- Localhost spoofing
- Trusted IP ranges
- X-Forwarded-For variations
- Client IP header manipulation

## URL Encoding
Applies various encoding techniques:
- URL encoding
- Double encoding
- Unicode encoding
- Special character handling

## Protocol Testing
Tests different protocol behaviors:
- HTTP/HTTPS switching
- Protocol downgrade
- Port variations
- Protocol-specific headers

## Path Traversal
Tests directory traversal patterns:
- Basic traversal (`../`)
- Encoded traversal
- Nested traversal
- Custom traversal patterns

## Future Techniques

The following techniques are planned for future releases:

- Advanced payload combinations
- Caching proxy bypass methods
- Wordlist-based path discovery
- Combined technique automation
- Custom success criteria
- Evidence collection
- Authentication scripting 