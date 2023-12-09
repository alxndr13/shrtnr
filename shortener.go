package main

import (
	"math/rand"
	"net/url"
	"strings"

	"github.com/sqids/sqids-go"
)

func (a *App) shortenUrl(url string, useTLD bool) (shortenedUrl string, err error) {
	s, err := sqids.New()
	if err != nil {
		return "", err
	}
	var randomIntegers []uint64
	for i := 0; i < 4; i++ {
		randomInt := rand.Intn(100)
		randomIntegers = append(randomIntegers, uint64(randomInt))
	}

	shortCode, err := s.Encode(randomIntegers)
	if err != nil {
		return "", err
	}

	if useTLD {
		topLevelDomain, err := getTopLevelDomain(url)
		if err != nil {
			return "", err
		}

		if topLevelDomain != "" {
			shortCode = topLevelDomain + "-" + shortCode
		}
	}

	err = a.saveToDatabase(shortCode, url)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func getTopLevelDomain(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(u.Host, "www.") {
		hostParts := strings.Split(u.Host, ".")
		if len(hostParts) > 1 {
			return hostParts[len(hostParts)-2], nil
		}
	} else {
		// Use the path if the URL doesn't start with "www."
		pathParts := strings.Split(u.Path, ".")
		if len(pathParts) > 1 {
			return pathParts[len(pathParts)-2], nil
		}
	}

	return "", nil
}
