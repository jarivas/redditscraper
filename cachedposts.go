package redditscraper

import (
	"time"
)

type CachedPosts struct {
	posts     []*Post
	timestamp time.Time
	expiresAt time.Time
}

func (c CachedPosts) GetNextId() string {
	l := len(c.posts)

	return c.posts[l-1].Name
}

func (c CachedPosts) GetPreviousId() string {
	return c.posts[0].Name
}

func (c CachedPosts) GetPosts() []*Post {
	return c.posts
}