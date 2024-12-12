package redditscraper

import (
	"fmt"
	"time"
)

type RedditScraper struct {
	client        *RedditClient
	subreddit     string
	maxPosts      int
	waitTime      time.Duration
	lastTimestamp time.Time
}

func (r RedditScraper) New(subreddit string, maxPosts, waitMilliseconds int) (*RedditScraper, error) {
	c, err := RedditClient{}.FromEnv()

	if err != nil {
		return nil, err
	}

	dummy := fmt.Sprintf("%vms", waitMilliseconds)
	waitTime, err := time.ParseDuration(dummy)

	if err != nil {
		return nil, err
	}

	rs := RedditScraper{
		client:        c,
		subreddit:     subreddit,
		maxPosts:      maxPosts,
		waitTime:      waitTime,
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

		time.Sleep(r.waitTime)
	}
}

func (r *RedditScraper) getPosts(listing PostListing) (*CachedPosts, error) {
	cachedPosts, err := r.client.GetTopPosts(r.subreddit, listing)

	if err != nil {
		return nil, err
	}

	if r.lastTimestamp.After(cachedPosts.timestamp) {
		return nil, nil
	}

	r.lastTimestamp = cachedPosts.timestamp

	return cachedPosts, nil
}
