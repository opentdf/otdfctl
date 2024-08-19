package utils

import (
	"errors"
	"net/url"
)

func NormalizeEndpoint(endpoint string) (*url.URL, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint is required")
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "http":
		if u.Port() == "" {
			u.Host += ":80"
		}
	case "https":
		if u.Port() == "" {
			u.Host += ":443"
		}
	default:
		return nil, errors.New("invalid scheme")
	}
	return u, nil
}
