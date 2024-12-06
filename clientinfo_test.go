package redditscraper

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestFromEnvClientInfo(t *testing.T) {
	err := godotenv.Load()

	if err != nil {
		t.Errorf("problem reading .env file, %v", err.Error())
	}

	i, err := clientInfo{}.new()

	if err != nil {
		t.Errorf("problem creating %v", err.Error())
	}

	if i.username != os.Getenv("REDDIT_USERNAME") {
		t.Errorf("Invalid username: %v", i.username)
	}

	if i.password != os.Getenv("REDDIT_PASSWORD") {
		t.Errorf("Invalid password: %v", i.password)
	}

	if i.clientId != os.Getenv("REDDIT_CLIENT_ID") {
		t.Errorf("Invalid client id: %v", i.clientId)
	}

	if i.appSecret != os.Getenv("REDDIT_APP_SECRET") {
		t.Errorf("Invalid app secret: %v", i.appSecret)
	}
}

func TestGetToken(t *testing.T) {
	i, err := clientInfo{}.new()
	
	if err != nil {
		t.Errorf("problem creating, %v", err.Error())
	}

	token, err := i.getToken()

	if err != nil {
		t.Errorf("problem getting the token, %v", err.Error())
	}

	if token.accessToken == "" {
		t.Error("access token is empty")
	}

	if token.expiresAt.Before(time.Now()) {
		t.Errorf("expires in is unvalid, %v", token.expiresAt)
	}
}
