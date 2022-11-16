package helpers

import (
	"os"
	"strings"
)

// RemoveDomainError returns true if it successfully removes all domain errors.
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(url, "https://", "", 1)
	newURL = strings.Replace(url, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}

	return true
}

// EnforceHTTP returns a string with https:// prefix.
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}

	return url
}
