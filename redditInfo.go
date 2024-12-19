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
	timeSleep time.Duration
}

func (ri RedditInfo) New(username, password, clientId, appSecret, timeSleep string) (*RedditInfo, error) {
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

	if timeSleep == "" {
		return nil, errors.New("redditinfo: invalid timeSleep")
	}

	d, err := time.ParseDuration(timeSleep)

	if err != nil {
		return nil, errors.New("redditinfo: invalid timeSleep format")
	}

	return &RedditInfo{
		username:  username,
		password:  password,
		clientId:  clientId,
		appSecret: appSecret,
		timeSleep: d,
	}, nil
}

func (ri RedditInfo) FromEnv() (*RedditInfo, error) {
	return ri.New(
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		os.Getenv("REDDIT_CLIENT_ID"),
		os.Getenv("REDDIT_APP_SECRET"),
		os.Getenv("REDDIT_TIME_SLEEP"),
	)
}

func (ri RedditInfo) getToken() (*oauthToken, error) {
	r, err := ri.requestToken()

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	t := tokenResponse{}
	json.NewDecoder(r.Body).Decode(&t)

	ot, err := t.convert()

	if err != nil {
		logResponse(r.Body)
		return nil, err
	}

	if ot.at == "" || ot.expires.Before(time.Now()) {
		logResponse(r.Body)
		return nil, errors.New("invalid oauth token")
	}

	return ot, nil
}

func (ri RedditInfo) requestToken() (*http.Response, error) {
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

	r, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		logResponse(r.Body)
		return nil, errors.New(r.Status)
	}

	return r, nil
}