package redditscraper

import (
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	tokenResponse := tokenResponse{
		Token:   "asd",
		Expires: 10,
	}

	token, err := tokenResponse.convert()

	if err != nil {
		t.Error(err)
	}

	if token.at != tokenResponse.Token {
		t.Errorf("Invalid access token %v", token.at)
	}

	now := time.Now()

	if token.expires.Before(now) {
		t.Errorf("Invalid expires a %v", token.expires)
	}
}
