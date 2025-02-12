package utils

import (
	"errors"
	"net/url"
	"strings"

	"github.com/google/uuid"
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
	for strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}
	return u, nil
}

func IsUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsURI(u string) bool {
	ur, err := url.Parse(u)
	if err != nil {
		return false
	}
	return ur.IsAbs()
}
