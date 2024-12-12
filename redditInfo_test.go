package redditscraper

import (
	"testing"
)

func TestFromEnv(t *testing.T) {
	_, err := RedditInfo{}.FromEnv()

	if err != nil {
		t.Error(err)
	}
}