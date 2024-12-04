# Reddit scraper
## Description
Scraps a particular subreddit

## Install
```go get https://github.com/jarivas/redditscraper```

## Usage
**.env**
```bash
REDDIT_USERNAME=reddit_bot
REDDIT_PASSWORD=snoo
REDDIT_CLIENT_ID=p-jcoLKBynTLew
REDDIT_APP_SECRET=gko_LXELoV07ZBNUXrvWZfzE3aI
```
**scraper.go**
```golang
package demo

import (
	"github.com/jarivas/redditscraper"
    "fmt"
)

func main() {
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
}
```

For more flexibility please check:
- **scraper.go**
- **scraper_test.go**
- **client.go**
- **client_test.go**
- **clientinfo.go**
- **clientinfo_test.go**