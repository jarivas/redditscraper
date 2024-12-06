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
	scraper, err := RedditScraper{}.New("AmItheAsshole", 10, 1000)

	if err != nil {
		t.Error(err)
	}

	if scraper == nil {
		t.Error("Scraper is nil")
	}

	c := make(chan *CachedPosts)
	e := make(chan error)

	go func() {
		scraper.Scrape(c, e, "")
	}()

	for {
		select{
			case cachedPosts := <- c: 
				if cachedPosts == nil {
					t.Error("cachedPosts is nil")
				}
				return
			case err = <- e: 
				t.Error(err)
				return
		}	
	}
}