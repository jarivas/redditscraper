package redditscraper

import (
	"io"
	"log"
	"net/http"
)

func logResponse(response *http.Response) {
	log.Println(response.Request.URL.Path)
	log.Println(response.Status)
	logBody(response.Body)
}

func logBody(body io.ReadCloser) {
	bytes, err := io.ReadAll(body)

	if err != nil {
		log.Printf("something happened when logging the response: %v", err)
	}

	r := string(bytes)

	log.Println(r)
}