package redditscraper

import (
	"errors"
	"fmt"
	"time"
)

type oauthToken struct {
	accessToken string
	expiresAt   time.Time
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (t tokenResponse) convert() (*oauthToken, error) {
	if t.AccessToken == "" || t.ExpiresIn == 0 {
		return nil, errors.New("empty token")
	}

	format := fmt.Sprintf("%vs", t.ExpiresIn)
	duration, err := time.ParseDuration(format)

	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(duration)

	result := oauthToken{
		accessToken: t.AccessToken,
		expiresAt:   expiresAt,
	}

	return &result, nil
}
