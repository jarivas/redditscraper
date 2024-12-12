package redditscraper

import (
	"time"
)

var currentToken *oauthToken = nil;

type oauthToken struct {
	accessToken string
	expiresAt   time.Time
}