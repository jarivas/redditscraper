package redditscraper

import (
	"testing"
	"time"
)

func TestCachePosts(t *testing.T) {
	client := client{}
	posts := cachePosts(client)

	client.cachePosts("asd", subredditCacheLong, posts)

	cachedPosts := *postsCache["asd"]

	if len(cachedPosts.posts) != 1 {
		t.Error("The invalid cache size")
	}

	if cachedPosts.expiresAt.Before(time.Now()) {
		t.Error("The invalid expiresAt")
	}
}

func TestGetCachedPosts(t *testing.T) {
	client := client{}

	postsCache = map[string]*CachedPosts{}

	cachedPosts, err := client.getCachedPosts("asd")

	if cachedPosts != nil {
		t.Error("the cache was not empty")
	}

	if err != nil {
		t.Errorf("an error happend %v", err.Error())
	}

	posts := cachePosts(client)

	client.cachePosts("asd", subredditCacheLong, posts)

	cachedPosts, err = client.getCachedPosts("asd")

	if err != nil {
		t.Errorf("error caching post %v", err)
	}

	if posts[0].Id != cachedPosts.posts[0].Id {
		t.Error("the cache was not empty")
	}
}

func TestNewClient(t *testing.T) {
	client, err := client{}.new()

	if err != nil {
		t.Errorf("error happened, %v", err.Error())
	}

	if client.token.accessToken == "" {
		t.Error("no access token")
	}
}

func TestGetPostsClient(t *testing.T) {
	client, err := client{}.new()

	if err != nil {
		t.Errorf("error happened, %v", err.Error())
	}

	listing := PostListing{
		Latest: true,
	}

	posts, err := client.GetTopPosts("AmItheAsshole", listing)

	if err != nil {
		t.Errorf("error happened, %v", err.Error())
	}

	if posts == nil || posts.posts == nil {
		t.Errorf("error happened, %v", err.Error())
	}

	if len(posts.posts) == 0 {
		t.Error("no posts token")
	}
}

func cachePosts(c client) []*Post {
	posts := []*Post{
		{
			Id:    "asd",
			Title: "asd",
			Body:  "asd",
		},
	}

	c.cachePosts("asd", subredditCacheLong, posts)

	return posts
}
