package redditscraper

import (
	"testing"
)

func TestNewScraper(t *testing.T) {
	scraper, err := RedditScraper{}.New("AmItheAsshole", 10, 999)

	if err != nil {
		t.Error(err)
	}

	if scraper == nil {
		t.Error("Scraper is nil")
	}
}

func TestScrape(t *testing.T) {
	scraper, err := RedditScraper{}.New("AmItheAsshole", 10, 999)

	if err != nil {
		t.Error(err)
	}

	if scraper == nil {
		t.Error("Scraper is nil")
	}

	c := make(chan *CachedPosts)

	go func() {
		err := scraper.Scrape(c)

		close(c)

		t.Error(err)
	}()

	cachedPosts := <- c

	if cachedPosts == nil {
		t.Error("cachedPosts is nil")
	}
}