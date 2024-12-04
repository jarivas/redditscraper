package redditscraper

import (
	"testing"
)

func TestGetNextId(t *testing.T) {
	cachedPosts := getCachedPosts()

	id := cachedPosts.GetNextId()

	if id != "02" {
		t.Error("wrong index")
	}
}

func TestGetPreviousId(t *testing.T) {
	cachedPosts := getCachedPosts()

	id := cachedPosts.GetPreviousId()

	if id != "01" {
		t.Error("wrong index")
	}
}

func TestGetPostsCachedPost(t *testing.T) {
	cachedPosts := getCachedPosts()

	posts := cachedPosts.GetPosts()

	if len(posts)!= 2 {
		t.Error("invalid posts")
	}
}

func getCachedPosts() CachedPosts {
	return CachedPosts{
		posts: []*Post{
			{
				Name:  "01",
				Title: "asd",
				Body:  "asd",
			},
			{
				Name:  "02",
				Title: "asd",
				Body:  "asd",
			},
		},
	}
}
