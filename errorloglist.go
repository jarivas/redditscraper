package redditscraper

import (
	"errors"
	"time"
)

type ErrorEntry struct {
	Timestamp time.Time
	Message   string
}

type ErrorLogList struct {
	ErrorList       []ErrorEntry
	Count           int
	MaxEntries      int
	MaxMilliseconds int64
}

func (e ErrorLogList) hasCapacity() bool {
	return e.Count < e.MaxEntries
}

func (e *ErrorLogList) LogError(err error) error {
	now := time.Now()

	writeError(err)

	e.removeOlderErrors(now)

	if e.hasCapacity() {
		e.ErrorList = append(e.ErrorList, ErrorEntry{
			Timestamp: now,
			Message: err.Error(),
		})

		return nil
	}

	return e.getCompiledError()
}

func (e *ErrorLogList) removeOlderErrors(now time.Time) {
	newList := make([]ErrorEntry, 0, e.MaxEntries)

	if len(e.ErrorList) == 0 {
		if cap(e.ErrorList) == 0 {
			e.ErrorList = newList
		}
	}

	for _, item := range e.ErrorList {
		d := -1 * item.Timestamp.Sub(now).Milliseconds()

		if d < e.MaxMilliseconds {
			newList = append(newList, item)
		}
	}

	e.ErrorList = newList
	e.Count = len(newList)
}

func (e ErrorLogList) getCompiledError() error {
	msg := "Error Messages"

	for _, item := range e.ErrorList {
		msg += ", " + item.Message
	}

	return errors.New(msg)
}
