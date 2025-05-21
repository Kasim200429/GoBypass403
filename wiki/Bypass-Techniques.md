# Bypass Techniques Reference

This document provides comprehensive technical documentation on the bypass methodologies implemented within GoBypass403.

## Methodology Classification

GoBypass403 employs a systematic categorization of bypass techniques, organized according to their technical implementation characteristics:

| Category | Implementation Focus | Technical Basis |
|----------|----------------------|----------------|
| Request Method | HTTP verb manipulation | RFC 7231 compliance variations |
| Path Manipulation | URL path structure | Path normalization algorithms |
| Header Manipulation | HTTP header fields | Request header processing logic |
| IP Spoofing | Origin identification | Source IP trust mechanisms |
| URL Encoding | Character encoding | URL parser implementation variances |
| Protocol Manipulation | Transport layer | Protocol handler inconsistencies |
| Path Traversal | Directory navigation | Path canonicalization logic |
| Proxy Bypass | Caching mechanisms | Intermediary processing behaviors |
| Specialized Vectors | Application-specific | Custom processing edge cases |
| Wordlist-based | Enumeration | Resource discovery patterns |
| Combined Techniques | Multi-vector approach | Defense depth circumvention |

## Technical Implementation Details

### 1. Request Method Manipulation

This technique exploits inconsistencies in HTTP method handling between access control mechanisms and resource handlers.

**Implementation Variants:**

```go
methodVariants := []string{"GET", "POST", "HEAD", "OPTIONS", "PUT", "DELETE", "TRACE", "CONNECT", "PATCH"}
```

**Technical Basis:**
- RFC 7231 defines standard HTTP methods, but implementations vary in their handling
- Access control configurations often focus on GET/POST while neglecting other methods
- Method overriding techniques (X-HTTP-Method-Override) may bypass method-based restrictions

**Example Implementation:**
```go
func TestMethodManipulation(baseURL string, client *http.Client, config Config) ([]Result, error) {
    var results []Result
    methodVariants := []string{"GET", "POST", "HEAD", "OPTIONS", "PUT", "DELETE", "TRACE", "CONNECT", "PATCH"}
    
    for _, method := range methodVariants {
        req, err := http.NewRequest(method, baseURL, nil)
        if err != nil {
            continue
        }
        
    }
    
    return results, nil
}
```

### 2. Path Manipulation

This technique exploits inconsistencies in URL path parsing and normalization between security controls and application servers.

**Implementation Variants:**
- Directory self-reference (`/./`)
- Case manipulation (`/ADMIN/` vs `/admin/`)
- Path parameter injection (`/admin;foo=bar/`)
- Trailing characters (`/admin//`, `/admin/./`)

**Technical Basis:**
- URL paths undergo multiple normalization processes across the request chain
- Different components may interpret path elements differently
- RFC 3986 compliance varies across implementations

### 3. Header Manipulation

This technique targets inconsistencies in HTTP header processing logic between security controls and application servers.

**Key Headers Utilized:**
- `X-Original-URL`
- `X-Rewrite-URL`
- `X-Forwarded-Host`
- `X-Host`
- `X-Custom-IP-Authorization`

**Technical Implementation Example:**
```go
headers := []struct {
    Header string
    Value  string
}{
    {"X-Original-URL", targetPath},
    {"X-Rewrite-URL", targetPath},
    {"X-Forwarded-Host", parsedURL.Host},
    {"X-Host", parsedURL.Host},
    {"X-Custom-IP-Authorization", "127.0.0.1"},

}
```

## Unicode Normalization Exploitation

This advanced technique targets inconsistencies in Unicode character handling across different system components.

**Technical Basis:**
- Unicode normalization forms: NFD, NFC, NFKD, NFKC
- Canonical vs. compatibility equivalence
- UTF-8 encoding variances

**Implementation:**
- Unicode character substitution
- Multi-byte character sequences
- Right-to-left overrides
- Zero-width characters

**Example Payload:**
```
/%uff0e%uff0e/%uff0e%uff0e/%uff0e%uff0e/etc/passwd
```

## Technical Analysis Methodology

When developing and implementing bypass techniques, GoBypass403 follows a systematic approach:

1. **HTTP Request Analysis**: Examining how various components process HTTP requests
2. **Path Processing Evaluation**: Analyzing path normalization across components
3. **Header Processing Assessment**: Studying header handling inconsistencies
4. **Response Differentiation**: Categorizing responses based on bypass effectiveness

## Efficacy Evaluation

Each bypass technique is evaluated against the following criteria:

- **Success Rate**: Percentage of successful bypasses across tested environments
- **WAF Detection Avoidance**: Capability to avoid detection by common WAF implementations
- **Implementation Consistency**: Reliability across different target environments

## References

- [RFC 7231: HTTP/1.1 Semantics and Content](https://tools.ietf.org/html/rfc7231)
- [RFC 3986: URI Generic Syntax](https://tools.ietf.org/html/rfc3986)
- [OWASP: Path Traversal](https://owasp.org/www-community/attacks/Path_Traversal)
- [HackTricks: 403 Bypass](https://book.hacktricks.xyz/pentesting-web/403-bypass-forbidden) 