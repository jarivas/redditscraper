package redditscraper

import (
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	tokenResponse := tokenResponse{
		AccessToken: "asd",
		ExpiresIn: 10,
	}

	token, err := tokenResponse.convert()

	if err != nil {
		t.Error(err)
	}

	if token.accessToken != tokenResponse.AccessToken {
		t.Errorf("Invalid access token %v", token.accessToken)
	}

	now := time.Now()

	if token.expiresAt.Before(now) {
		t.Errorf("Invalid expires a %v", token.expiresAt)
	}
}