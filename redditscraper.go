package redditscraper

import (
	"time"
)

type RedditScraper struct {
	client        *RedditClient
	subreddit     string
	maxPosts      int
	lastTimestamp time.Time
}

func (r RedditScraper) New(subreddit string) (*RedditScraper, error) {
	c, err := RedditClient{}.FromEnv("AmItheAsshole")

	if err != nil {
		return nil, err
	}

	rs := RedditScraper{
		client:        c,
		subreddit:     subreddit,
		maxPosts:      maxPosts,
		lastTimestamp: time.Now(),
	}

	return &rs, nil
}

func (r RedditScraper) Scrape(c chan<- *CachedPosts, e chan<- error, nextId string) {
	listing := PostListing{
		Id:    nextId,
		Limit: r.maxPosts,
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
