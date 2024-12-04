package redditscraper

import (
	"time"
)

type RedditScraper struct {
	client        *Client
	subreddit     string
	maxPosts      int
	waitTime      time.Duration
	lastTimestamp time.Time
}

func (r RedditScraper) New(info ClientInfo, subreddit string, maxPosts, waitMilliseconds int) (*RedditScraper, error){
	c, err := Client{}.New(info)

	if err != nil {
		return nil, err
	}

	waitTime := time.Duration(waitMilliseconds) * time.Microsecond

	rs := RedditScraper{
		client:        c,
		subreddit:     subreddit,
		maxPosts:      maxPosts,
		waitTime:      waitTime,
		lastTimestamp: time.Now(),
	}

	return &rs, nil
}

func (r RedditScraper) FromEnv(subreddit string, maxPosts, waitMilliseconds int) (*RedditScraper, error) {
	i := ClientInfo{}.fromEnv()

	return r.New(i, subreddit, maxPosts, waitMilliseconds)
}

func (r RedditScraper) Scrape(c chan<- *CachedPosts) error {
	listing := PostListing{
		Latest: true,
		Limit:  r.maxPosts,
	}

	for {
		cachedPosts, err := r.getPosts(listing)

		if err != nil {
			return err
		}

		time.Sleep(r.waitTime)

		if cachedPosts != nil {
			listing.Id = cachedPosts.GetNextId()
			listing.Latest = false

			c <- cachedPosts
		}
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
