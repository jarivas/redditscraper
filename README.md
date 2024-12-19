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
REDDIT_TIME_SLEEP=30s
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
   rs, err := redditscraper.RedditScraper{}.FromEnv("redditdev")

	if err != nil {
		log.Fatal(err)
	}

	p := make(chan *Post)
	e := make(chan error)

	listing := redditscraper.PostListing{
		Limit: redditscraper.MaxPosts,
		Id: "1h3wrtm"
	}

	go rs.Listen(SubredditBest, listing, p, e)

	for {
		select{
		case post := <- p:
			log.Println(post)
			return
		case err = <- e:
			log.Fatal(err)
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
