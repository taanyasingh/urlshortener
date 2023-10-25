package helpers

import (
	"log"
	"net/url"
	"os"
	"strings"
)

// EnforceHTTP ...
func EnforceHTTP(url string) string {
	// make every url https
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// RemoveDomainError
func RemoveDomainError(url string) bool {
	// basically this functions removes all the commonly found
	// prefixes from URL such as http, https, www
	// then checks of the remaining string is the DOMAIN itself
	if url == os.Getenv("DOMAIN") {
		return false
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}
	return true
}

func GetDomainFromURL(u string) (string, error) {
	url, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	return hostname, nil
}
