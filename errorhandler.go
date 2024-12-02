package redditscraper

import (
	"log"
	"os"
	"time"
	"errors"
)


type ErrorInfo struct {
	Timestamp time.Time
	Message   string
}

var errorList = []ErrorInfo{}

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

func logError(err error) error {
	now := time.Now()

	writeError(err)
	
	if checkErrorFrequency(now) {
		return nil
	}

	return getCompiledError()
}

func checkErrorFrequency(now time.Time) bool {
	if len(errorList) < maxErrors {
		return true
	}

	clearErrorList(now)

	return len(errorList) < maxErrors	
}

func clearErrorList(now time.Time) {
	newList := []ErrorInfo{}

	for _, item := range(errorList) {
		d := item.Timestamp.Sub(now).Milliseconds()

		if (d < maxMilliseconds) {
			l := len(newList)

			if l > 0 {
				l--
			}

			newList[l] = item
		}
	}

	errorList = newList
}

func getCompiledError() error {
	msg := "Error Messages: "

	for _, item := range(errorList) {
		msg += item.Message + ", "
	}

	return errors.New(msg)
}