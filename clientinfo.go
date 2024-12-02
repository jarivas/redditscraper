package redditscraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ClientInfo struct {
	Username string
	Password string
	ClientId string
	AppSecret string
}

func (info *ClientInfo) GetClient() (*Client, error) {
	token, err := info.getToken()

	if err != nil {
		return nil, err
	}

	client := Client{
		info:  info,
		token: token,
	}

	return &client, nil
}

func (info ClientInfo) getToken() (*oauthToken, error) {
	response, err := info.getTokenResponse()

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	token := tokenResponse{}
	json.NewDecoder(response.Body).Decode(token)

	return token.convert()
}

func (info ClientInfo) getTokenResponse() (*http.Response, error) {

	body := fmt.Sprintf(
		"grant_type=password&username=%v&password=%v",
		info.Username,
		info.Password,
	)

	request, err := http.NewRequest(
		"POST",
		apiTokenUrl,
		strings.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(info.Username, info.Password)

	return http.DefaultClient.Do(request)
}
