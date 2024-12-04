package redditscraper

import (
	"errors"
	"fmt"
	"time"
)

type ErrorEntry struct {
	timestamp time.Time
	err       error
}

type ErrorLogList struct {
	errorList       []ErrorEntry
	count           int
	maxEntries      int
	maxMilliseconds int64
}

func (e ErrorLogList) New(maxEntries, maxMilliseconds int) ErrorLogList {
	list := make([]ErrorEntry, 0, maxEntries)

	return ErrorLogList{
		errorList:       list,
		maxEntries:      maxEntries,
		maxMilliseconds: int64(maxMilliseconds),
	}
}

func (e *ErrorLogList) LogError(err error) error {
	if e.maxEntries == 0 {
		return errors.New("please use newErrorLogList before call LogError")
	}

	now := time.Now()

	WriteError(err)

	e.removeOlderErrors(now)

	if e.hasCapacity() {
		e.errorList = append(e.errorList, ErrorEntry{
			timestamp: now,
			err:       err,
		})

		return nil
	}

	return e.getCompiledError()
}

func (e ErrorLogList) hasCapacity() bool {
	return e.count < e.maxEntries
}

func (e *ErrorLogList) removeOlderErrors(now time.Time) {
	newList := make([]ErrorEntry, 0, e.maxEntries)

	for _, item := range e.errorList {
		d := -1 * item.timestamp.Sub(now).Milliseconds()

		if d < e.maxMilliseconds {
			newList = append(newList, item)
		}
	}

	e.errorList = newList
	e.count = len(newList)
}

func (e ErrorLogList) getCompiledError() error {
	msg := fmt.Sprintf("more than %v (Max Entries) errors occurred", e.maxEntries)
	err := errors.New(msg)

	for _, item := range e.errorList {
		err = errors.Join(err, item.err)
	}

	return err
}
