package redditscraper

import (
	"time"
)

type RedditScraper struct {
	client        *RedditClient
	lastTimestamp time.Time
}

func (r RedditScraper) New(c *RedditClient) (*RedditScraper, error) {
	rs := RedditScraper{
		client:        c,
		lastTimestamp: time.Now(),
	}

	return &rs, nil
}

func (r RedditScraper) FromEnv(subreddit string) (*RedditScraper, error) {
	c, err := RedditClient{}.FromEnv("AmItheAsshole")

	if err != nil {
		return nil, err
	}

	return r.New(c)
}

func (r *RedditScraper) Scrape(c chan<- *CachedPosts, e chan<- error, nextId string) {
	listing := PostListing{
		Id:    nextId,
		Limit: maxPosts,
	}

	if nextId == "" {
		listing.Latest = true
	}

	for {
		cachedPosts, err := r.getPosts(listing)

		if err == nil {
			if cachedPosts != nil && cachedPosts.HasPost() {
				listing.Id = cachedPosts.GetNextId()
				listing.Latest = false

				c <- cachedPosts
			}
		} else {
			e <- err
		}

		time.Sleep(waitTime)
	}
}

func (r *RedditScraper) getPosts(listing PostListing) (*CachedPosts, error) {
	cachedPosts, err := r.client.GetTopPosts(listing)

	if err != nil {
		return nil, err
	}

	if r.lastTimestamp.After(cachedPosts.timestamp) {
		return nil, nil
	}

	r.lastTimestamp = cachedPosts.timestamp

	return cachedPosts, nil
}
