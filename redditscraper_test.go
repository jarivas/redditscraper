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

	go rs.Listen(SubredditBest, p, e)

	for {
		select{
		case post := <- p:
			log.Println(post)
			return
		case err = <- e:
			t.Error(err)
		}
	}
}