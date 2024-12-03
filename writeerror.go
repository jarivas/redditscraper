package redditscraper

import (
	"log"
	"os"
)

func writeError(err error) {
	file, err2 := os.OpenFile("./logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err2 != nil {
		log.Fatal(err)
		log.Fatal(err2)
	}

	defer file.Close()

	log.SetOutput(file)

	log.Fatal(err)
}