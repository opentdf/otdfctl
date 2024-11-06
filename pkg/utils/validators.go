package utils

import (
	"errors"
	"net/url"
	"strings"
)

// NormalizeEndpoint validates, defaults, and standardizes an endpoint URL by setting default scheme/port and removing trailing slashes.
func NormalizeEndpoint(endpoint string) (*url.URL, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint is required")
	}

	// Add default scheme if missing
	if !strings.Contains(endpoint, "://") {
		endpoint = "https://" + endpoint
	}

	// Parse the URL
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// Validate scheme and set default ports for http/https
	switch u.Scheme {
	case "http":
		if u.Port() == "" {
			u.Host = u.Hostname() + ":80"
		}
	case "https":
		if u.Port() == "" {
			u.Host = u.Hostname() + ":443"
		}
	default:
		return nil, errors.New("invalid scheme: only http and https are supported")
	}

	// Trim trailing slashes from path
	u.Path = strings.TrimRight(u.Path, "/")

	return u, nil
}
