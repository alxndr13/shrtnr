package main

import (
	"math/rand"

	"github.com/sqids/sqids-go"
)

func (a *App) shortenUrl(url string) (shortenedUrl string, err error) {
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

	err = a.saveToDatabase(shortCode, url)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}
