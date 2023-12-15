package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/log"
)

var TestApp App

func init() {
	TestApp.DbPath = "shrtnr_test.db"
	err := TestApp.createDatabaseIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

}

func TestShortening(t *testing.T) {

	var testingUrls []string = []string{
		"google.com",
		"https://youtu.be/dQw4w9WgXcQ?si=d1lqXPpbH0N8Pq2h",
	}

	for _, url := range testingUrls {
		shortenedUrl, err := TestApp.shortenUrl(url)
		if len(shortenedUrl) == 0 {
			t.Errorf("shortening url '%s' failed", url)
		}
		if err != nil {
			t.Error(err)
		}

		split := strings.Split(shortenedUrl, "/")

		loadedUrl, err := TestApp.getFromDatabase(split[len(split)-1])
		if err != nil {
			t.Error(err)
		}

		if loadedUrl != url {
			t.Errorf("loaded url and url do not match, expected %s, got %s", url, loadedUrl)
		}
	}

}
