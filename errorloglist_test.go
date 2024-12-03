package redditscraper

import (
	"testing"
	"time"
)

func TestCheckCapacity(t *testing.T) {
	e := ErrorLogList{
		MaxEntries:      10,
		MaxMilliseconds: 999,
	}

	if !e.hasCapacity() {
		t.Errorf("The count is greater than MaxEntries %v", e.Count)
	}

	e.Count = 10

	if e.hasCapacity() {
		t.Errorf("The count is less than MaxEntries %v", e.Count)
	}
}

func TestRemoveOlderErrors(t *testing.T) {
	d, _ := time.ParseDuration("-10s")
	now := time.Now()
	past := now.Add(d)

	errorList := []ErrorEntry{
		{
			Timestamp: past,
		},
	}

	e := ErrorLogList{
		ErrorList:       errorList,
		MaxEntries:      10,
		MaxMilliseconds: 999,
	}

	e.removeOlderErrors(now)

	if e.Count > 0 {
		t.Errorf("Older Entries where not removed %v", e.Count)
	}
}

func TestGetCompiledError(t *testing.T) {
	errMsg := "Hello"
	errMsg2 := "Bye"

	errorList := []ErrorEntry{
		{
			Message: errMsg,
		},
		{
			Message: errMsg2,
		},
	}

	errorLogList := ErrorLogList{
		ErrorList: errorList,
	}

	msg := errorLogList.getCompiledError().Error()

	if msg != "Error Messages, Hello, Bye"{
		t.Errorf("There is a problem with the message | %v", msg)
	}

}