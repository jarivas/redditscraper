package redditscraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ClientInfo struct {
	Username  string
	Password  string
	ClientId  string
	AppSecret string
}

func (i ClientInfo) fromEnv() ClientInfo {
	err := godotenv.Load()

	if err != nil {
		writeError(err)
	}

	return ClientInfo{
		Username:  os.Getenv("REDDIT_USERNAME"),
		Password:  os.Getenv("REDDIT_PASSWORD"),
		ClientId:  os.Getenv("REDDIT_CLIENT_ID"),
		AppSecret: os.Getenv("REDDIT_APP_SECRET"),
	}
}

func (i ClientInfo) getToken() (*oauthToken, error) {
	response, err := i.getTokenResponse()

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	defer response.Body.Close()

	token := tokenResponse{}
	json.NewDecoder(response.Body).Decode(&token)

	return token.convert()
}

func (i ClientInfo) getTokenResponse() (*http.Response, error) {

	body := fmt.Sprintf(
		"grant_type=password&username=%v&password=%v",
		i.Username,
		i.Password,
	)

	request, err := http.NewRequest(
		"POST",
		apiTokenUrl,
		strings.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(i.ClientId, i.AppSecret)

	return http.DefaultClient.Do(request)
}
