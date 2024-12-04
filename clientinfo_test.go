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

	i := ClientInfo{}.fromEnv()

	if i.Username != os.Getenv("REDDIT_USERNAME") {
		t.Errorf("Invalid username: %v", i.Username)
	}

	if i.Password != os.Getenv("REDDIT_PASSWORD") {
		t.Errorf("Invalid password: %v", i.Password)
	}

	if i.ClientId != os.Getenv("REDDIT_CLIENT_ID") {
		t.Errorf("Invalid client id: %v", i.ClientId)
	}

	if i.AppSecret != os.Getenv("REDDIT_APP_SECRET") {
		t.Errorf("Invalid app secret: %v", i.AppSecret)
	}
}


func TestGetToken(t *testing.T) {
	i := ClientInfo{}.fromEnv()

	token , err := i.getToken()

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