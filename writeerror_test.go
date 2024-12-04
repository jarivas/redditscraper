package redditscraper

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestWriteError(t *testing.T) {
	err := os.Remove(writeErrorPath)

	if err != nil {
		t.Errorf("Problem removing the log file %v", err.Error())
	}

	WriteError(errors.New("test"))

	b, err := os.ReadFile(writeErrorPath)

	if err != nil {
		t.Errorf("Problem reading the log file %v", err.Error())
	}

	content := string(b)

	if !strings.Contains(content, "test") {
		t.Error("the content of log file does not have test")
	}
}
