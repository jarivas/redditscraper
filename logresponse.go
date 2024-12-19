package redditscraper

import (
	"io"
	"log"
)

func logResponse(body io.ReadCloser) {
	bytes, err := io.ReadAll(body)

	if err != nil {
		log.Printf("something happened when logging the response: %v", err)
	}

	r := string(bytes)

	log.Println(r)
}