package redditscraper

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := RedditScraper{}.FromEnv("redditdev")

	if err != nil {
		t.Error(err)
	}
}

func TestListen(t *testing.T) {
	rs, err := RedditScraper{}.FromEnv("redditdev")

	if err != nil {
		t.Error(err)
	}

	p := make(chan *Post)
	e := make(chan error)

	listing := PostListing{
		Limit: MaxPosts,
	}

	go rs.Listen(SubredditBest, listing, p, e)

	for {
		select {
		case post := <-p:
			log.Println(post)
			return
		case err = <-e:
			t.Error(err)
		}
	}
}
