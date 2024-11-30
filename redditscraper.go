package redditscraper

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/jarivas/redditreader"
	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

const MAX_MILLISECONDS = 999

type Post struct {
	FullID string
	Title  string
	Body   string
}

func Scrape(c chan<- Post) error {
	waitTime, err := getWaitTime()
	nextFullID := ""

	if err != nil {
		return err
	}

	for {
		posts, err := redditreader.GetNext(redditreader.MAX_POSTS, nextFullID)

		if err != nil {
			err = logError(err)

			if err != nil {
				return err
			}
		} else {
			nextFullID = posts[len(posts)-1].FullID

			for _, post := range posts {
				c <- convertToPost(post)
			}
		}

		time.Sleep(waitTime)
	}
}

func getWaitTime() (time.Duration, error) {
	err := godotenv.Load()

	if err != nil {
		return 0, errors.New("error loading .env file")
	}

	dummy, ok := os.LookupEnv("REDDIT_WAIT_MILLISECONDS")

	if !ok {
		return 0, errors.New("REDDIT_WAIT_MILLISECONDS not present on .env file or env var")
	}

	waitTime, err := strconv.Atoi(dummy)

	if err != nil {
		return 0, errors.New("REDDIT_WAIT_MILLISECONDS is not an integer")
	}

	if waitTime < 0 || waitTime > MAX_MILLISECONDS {
		return 0, errors.New("REDDIT_WAIT_MILLISECONDS value is invalid")
	}

	return time.Duration(waitTime) * time.Microsecond, nil
}

func convertToPost(post *reddit.Post) Post {
	return Post{
		FullID: post.FullID,
		Title:  post.Title,
		Body:   post.Body,
	}
}