# Reddit scraper
## Description
Scraps a particular subreddit

## Install
```go get https://github.com/jarivas/redditscraper```

## Usage
```golang
package demo

import (
	"github.com/jarivas/redditscraper"
)

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
```