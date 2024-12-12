package redditscraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type RedditInfo struct {
	username  string
	password  string
	clientId  string
	appSecret string
}

func (ri RedditInfo) New(username, password, clientId, appSecret string) (*RedditInfo, error) {
	if username == "" {
		return nil, errors.New("redditinfo: invalid username")
	}

	if password == "" {
		return nil, errors.New("redditinfo: invalid password")
	}

	if clientId == "" {
		return nil, errors.New("redditinfo: invalid clientId")
	}

	if appSecret == "" {
		return nil, errors.New("redditinfo: invalid appSecret")
	}

	return &RedditInfo{
		username:  username,
		password:  password,
		clientId:  clientId,
		appSecret: appSecret,
	}, nil
}

func (ri RedditInfo) FromEnv() (*RedditInfo, error) {
	return ri.New(
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		os.Getenv("REDDIT_CLIENT_ID"),
		os.Getenv("REDDIT_APP_SECRET"),
	)
}

func (ri RedditInfo) getToken() (*oauthToken, error) {
	response, err := ri.getTokenResponse()

	if err != nil {
		return nil, err
	}

	if response.StatusCode == 429 {
		ri.sleep()
		return ri.getToken()
	}

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	defer response.Body.Close()

	token := tokenResponse{}
	json.NewDecoder(response.Body).Decode(&token)

	return token.convert()
}

func (ri RedditInfo) getTokenResponse() (*http.Response, error) {
	body := fmt.Sprintf(
		"grant_type=password&username=%v&password=%v",
		ri.username,
		ri.password,
	)

	request, err := http.NewRequest(
		"POST",
		apiTokenUrl,
		strings.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(ri.clientId, ri.appSecret)

	return http.DefaultClient.Do(request)
}

func (ri RedditInfo) sleep() {
	d, _ := time.ParseDuration("30")

	time.Sleep(d)
}