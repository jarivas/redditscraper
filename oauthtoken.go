package redditscraper

import (
	"time"
)

var currentToken *oauthToken = nil

type oauthToken struct {
	at      string
	expires time.Time
}
