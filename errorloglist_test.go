package redditscraper

import (
	"errors"
	"testing"
	"time"
)

func TestCheckCapacity(t *testing.T) {
	e := ErrorLogList{
		maxEntries:      10,
		maxMilliseconds: 999,
	}

	if !e.hasCapacity() {
		t.Errorf("The count is greater than MaxEntries %v", e.count)
	}

	e.count = 10

	if e.hasCapacity() {
		t.Errorf("The count is less than MaxEntries %v", e.count)
	}
}

func TestRemoveOlderErrors(t *testing.T) {
	d, _ := time.ParseDuration("-10s")
	now := time.Now()
	past := now.Add(d)

	errorList := []ErrorEntry{
		{
			timestamp: past,
		},
	}

	e := ErrorLogList{
		errorList:       errorList,
		maxEntries:      10,
		maxMilliseconds: 999,
	}

	e.removeOlderErrors(now)

	if e.count > 0 {
		t.Errorf("Older Entries where not removed %v", e.count)
	}
}

func TestGetCompiledError(t *testing.T) {
	errMsg := "Hello"
	errMsg2 := "Bye"

	errorList := []ErrorEntry{
		{
			err: errors.New(errMsg),
		},
		{
			err: errors.New(errMsg2),
		},
	}

	errorLogList := ErrorLogList{
		errorList:  errorList,
		maxEntries: 1,
	}

	msg := errorLogList.getCompiledError().Error()

	if msg != "more than 1 (Max Entries) errors occurred\nHello\nBye" {
		t.Errorf("There is a problem with the message | %v", msg)
	}
}

func TestNew(t *testing.T) {
	errorLogList := ErrorLogList{}.New(1, 999)

	if len(errorLogList.errorList) != 0 {
		t.Errorf("errorList length is invalid")
	}

	if cap(errorLogList.errorList) != 1 {
		t.Errorf("errorList capacity is invalid")
	}

	if errorLogList.maxMilliseconds != 999 {
		t.Errorf("errorList max milliseconds is invalid")
	}
}

func TestLogError(t *testing.T) {
	errorLogList := ErrorLogList{}.New(1, 999)

	err := errorLogList.LogError(errors.New("Test"))

	if err != nil {
		t.Errorf("problem with errorlist capacity: %v", err.Error())
	}

	err = errorLogList.LogError(errors.New("Test2"))

	if err == nil {
		t.Error("problem with errorlist capacity")
	}
}
