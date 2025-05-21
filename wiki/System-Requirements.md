# System Requirements

This document outlines the technical specifications and system requirements for deploying and operating GoBypass403.

## Hardware Requirements

| Component | Minimum | Recommended | Notes |
|-----------|---------|-------------|-------|
| Processor | Dual-core 1.8GHz | Quad-core 2.5GHz+ | Additional cores improve concurrent operation performance |
| Memory | 1GB RAM | 4GB+ RAM | Higher memory capacity enables larger payload execution |
| Disk Space | 100MB | 500MB | Additional space required for scan results and logs |
| Network | 10 Mbps | 100+ Mbps | Bandwidth impacts concurrent request throughput |

## Software Requirements

### Operating System Compatibility

GoBypass403 is compatible with the following operating systems:

| Operating System | Version | Architecture | Notes |
|-----------------|---------|--------------|-------|
| Linux | Kernel 4.0+ | x86_64, ARM64 | Preferred development platform |
| macOS | 10.14+ | x86_64, ARM64 | Full compatibility, including Apple Silicon |
| Windows | 10/11, Server 2019+ | x86_64 | Confirmed functional via testing |
| FreeBSD | 12.0+ | x86_64 | Limited testing, report issues if encountered |

### Runtime Dependencies

#### Go Environment

- Go version 1.21 or higher

#### External Dependencies

| Dependency | Version | Purpose | Installation Method |
|------------|---------|---------|---------------------|
| CA Certificates | Current | Certificate validation | System package |

## Network Requirements

GoBypass403 requires the following network configurations:

- Outbound TCP connections allowed (default HTTP/HTTPS ports)
- DNS resolution capability
- IPv4 connectivity (IPv6 supported but not required)

### Proxy Configuration

For operation through proxy servers:

- HTTP/HTTPS proxy support via environment variables

## Permission Requirements

| Environment | Requirements | Notes |
|-------------|--------------|-------|
| Linux/macOS | Regular user privileges | Root not required for standard operation |
| Windows | Standard user account | Administrator not required |

## Additional Considerations

### Firewall Configuration

If operating behind a firewall, ensure the following:

- Outbound connections allowed to ports 80 and 443
- Outbound connections to custom ports if non-standard target ports are used

### Rate Limiting

To prevent target saturation:

- Default concurrency limited to 10 threads
- Manual adjustment possible with the `-t` flag

### System Load

GoBypass403 is designed to minimize system resource impact:

- CPU utilization scales with concurrent threads
- Memory utilization remains relatively constant

## Coming in Future Updates

> Note: The following features are planned for upcoming releases and are not yet implemented.

- Compatibility verification tool (`cmd/compat_check.go`)
- Advanced proxy chaining
- Custom success criteria
- Evidence collection features
- Authentication scripting 