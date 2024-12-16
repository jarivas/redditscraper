package redditscraper

import (
	"time"
)

type RedditScraper struct {
	client        *RedditClient
	lastTimestamp time.Time
}

func (r RedditScraper) New(c *RedditClient) *RedditScraper {
	return &RedditScraper{
		client:        c,
		lastTimestamp: time.Now(),
	}
}

func (r RedditScraper) FromEnv(subreddit string) (*RedditScraper, error) {
	c, err := RedditClient{}.FromEnv("AmItheAsshole")

	if err != nil {
		return nil, err
	}

	return r.New(c), nil
}

func (r *RedditScraper) GetSubreddit() string {
	return r.client.subreddit
}

func (r *RedditScraper) Scrape(c chan<- *CachedPosts, e chan<- error, nextId, sort string) {
	listing := PostListing{
		Id:    nextId,
		Limit: maxPosts,
	}

	if nextId == "" {
		listing.Latest = true
	}

	for {
		cachedPosts, err := r.getPosts(listing, sort)

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

func (r *RedditScraper) ScrapeAll(c chan<- *CachedPosts, e chan<- error) {
	r.Scrape(c, e, "", SubredditBest)
	r.Scrape(c, e, "", SubredditNew)
	r.Scrape(c, e, "", SubredditRandom)
	r.Scrape(c, e, "", SubredditRising)
	r.Scrape(c, e, "", SubredditTop)
	r.Scrape(c, e, "", SubredditControversial)
}

func (r *RedditScraper) getPosts(listing PostListing, sort string) (*CachedPosts, error) {
	cachedPosts, err := r.client.GetSortedPost(listing, sort)

	if err != nil {
		return nil, err
	}

	if r.lastTimestamp.After(cachedPosts.timestamp) {
		return nil, nil
	}

	r.lastTimestamp = cachedPosts.timestamp

	return cachedPosts, nil
}
