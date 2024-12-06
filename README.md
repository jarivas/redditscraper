# Reddit scraper
## Description
Scraps a particular subreddit

## Install
```go get github.com/jarivas/redditscraper```

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
    "log"
)

func main() {
    scraper, err := RedditScraper{}.New("AmItheAsshole", 10, 999)

	if err != nil {
		log.Fatal(err)
	}

	if scraper == nil {
		log.Fatal("Scraper is nil")
	}

	c := make(chan *CachedPosts)
	e := make(chan error)

	go func() {
		scraper.Scrape(c, e)
	}()

	for {
		select{
			case cachedPosts := <- c: 
				if cachedPosts == nil {
					log.Fatal("cachedPosts is nil")
				}
				return
			case err = <- e: 
				log.Fatal(err)
				return
		}	
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
