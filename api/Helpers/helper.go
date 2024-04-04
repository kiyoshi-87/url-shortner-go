package helpers

import (
	"os"
	"strings"
)

func EnforceHttp(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	domain := os.Getenv("DOMAIN")

	if url == domain {
		return false
	}

	//REMOVING COMMON PREFIXES
	newUrl := strings.Replace(url, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "https://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0] //DOUBT

	return newUrl != domain
}
