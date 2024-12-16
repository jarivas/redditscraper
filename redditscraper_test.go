package redditscraper

import (
	"testing"
)

func TestNewScraper(t *testing.T) {
	scraper, err := RedditScraper{}.FromEnv("AmItheAsshole")

	if err != nil {
		t.Error(err)
	}

	if scraper == nil {
		t.Error("Scraper is nil")
	}
}

func TestScrape(t *testing.T) {
	scraper, err := RedditScraper{}.FromEnv("AmItheAsshole")

	if err != nil {
		t.Error(err)
	}

	if scraper == nil {
		t.Error("Scraper is nil")
	}

	c := make(chan *CachedPosts)
	e := make(chan error)

	go scraper.Scrape(c, e, "", SubredditBest)

	for {
		select {
		case cachedPosts := <-c:
			if cachedPosts == nil {
				t.Error("cachedPosts is nil")
			}
			return
		case err = <-e:
			t.Error(err)
			return
		}
	}
}

func TestScrapeAll(t *testing.T) {
	scraper, err := RedditScraper{}.FromEnv("AmItheAsshole")

	if err != nil {
		t.Error(err)
	}

	if scraper == nil {
		t.Error("Scraper is nil")
	}

	c := make(chan *CachedPosts)
	e := make(chan error)

	go scraper.ScrapeAll(c, e)

	for {
		select {
		case cachedPosts := <-c:
			if cachedPosts == nil {
				t.Error("cachedPosts is nil")
			}
			return
		case err = <-e:
			t.Error(err)
			return
		}
	}
}