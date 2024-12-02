package redditscraper

import (
	"time"
)

type Client struct {
	info  *ClientInfo
	token *oauthToken
}

type cachedPosts struct {
	posts     []*Post
	expiresAt time.Time
}

var cache = make(map[string]*cachedPosts)

func (c Client) GetBestPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditBest, subredditCacheLong)
}

func (c Client) GetNewPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditNew, subredditCache1Short)
}

func (c Client) GetRandomPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditRandom, subredditCache1Short)
}

func (c Client) GetRisingPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditRising, subredditCacheLong)
}

func (c Client) GetTopPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditTop, subredditCacheLong)
}

func (c Client) GetControversialPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditControversial, subredditCacheLong)
}

func (c Client) getPosts(listing PostListing, subreddit, sort, duration string) ([]*Post, error) {
	url := listing.getUrl(subreddit, sort)

	posts, err := c.getCachedPosts(url, sort, duration)

	if err != nil {
		writeError(err)
	}
}

func (c Client) getCachedPosts(url string, subreddit, sort string) ([]*Post, error) {
	cachedPost := cache[url]

	if cachedPost == nil {
		return nil, nil
	}

	now := time.Now()

	if now.After(cachedPost.expiresAt) {
		cache[url] = nil

		return nil, nil
	}

	return cachedPost.posts, nil
}

func (c Client) cachePosts(url, duration string, posts []*Post) {
	d, err := time.ParseDuration(duration)

	if err != nil {
		writeError(err)
		return
	}

	expiresIn := time.Now().Add(d)

}