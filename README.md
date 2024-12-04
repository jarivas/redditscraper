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
    "fmt"
)

scraper, err := RedditScraper{}.New("AmItheAsshole", 10, 999)

if err != nil {
    fmt.Errorf("There was an error: %v", err)
}

if scraper == nil {
    fmt.Error("Scraper is nil")
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