package redditscraper

import (
	"errors"
	"fmt"
	"time"
)

type tokenResponse struct {
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
}

func (t tokenResponse) convert() (*oauthToken, error) {
	if t.Token == "" || t.Expires == 0 {
		return nil, errors.New("empty token")
	}

	format := fmt.Sprintf("%vs", t.Expires)
	duration, err := time.ParseDuration(format)

	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(duration)

	result := oauthToken{
		at:      t.Token,
		expires: expiresAt,
	}

	return &result, nil
}
