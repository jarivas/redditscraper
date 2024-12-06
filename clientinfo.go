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

type clientInfo struct {
	username  string
	password  string
	clientId  string
	appSecret string
}

func (i clientInfo) new() (*clientInfo, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	return &clientInfo{
		username:  os.Getenv("REDDIT_USERNAME"),
		password:  os.Getenv("REDDIT_PASSWORD"),
		clientId:  os.Getenv("REDDIT_CLIENT_ID"),
		appSecret: os.Getenv("REDDIT_APP_SECRET"),
	}, nil
}

func (i clientInfo) getToken() (*oauthToken, error) {
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

func (i clientInfo) getTokenResponse() (*http.Response, error) {

	body := fmt.Sprintf(
		"grant_type=password&username=%v&password=%v",
		i.username,
		i.password,
	)

	request, err := http.NewRequest(
		"POST",
		apiTokenUrl,
		strings.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(i.clientId, i.appSecret)

	return http.DefaultClient.Do(request)
}
