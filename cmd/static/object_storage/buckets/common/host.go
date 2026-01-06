package common

import (
	"net/url"
	"strings"
)

type location struct {
	Host string
	Path string
}

func BuildHost(bucketName string, region string) (string, error) {
	baseURL := strings.ReplaceAll(TemplateURL, "{{region}}", region)

	parts := parseLocation(bucketName)

	bucketURL, err := url.JoinPath(baseURL, parts.Host, parts.Path)
	if err != nil {
		return "", err
	}

	bucketURL = strings.TrimSuffix(bucketURL, "/")

	return bucketURL, nil
}

func parseLocation(input string) location {
	u, err := url.Parse(input)
	if err != nil {
		return location{Host: input}
	}

	if host := u.Hostname(); host != "" {
		return location{
			Host: host,
			Path: strings.Trim(u.Path, "/"),
		}
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) > 0 {
		return location{
			Host: parts[0],
			Path: strings.Join(parts[1:], "/"),
		}
	}

	return location{Host: input}
}
