package redditscraper

import (
	"testing"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var subreddit string = "AmItheAsshole"

func TestGetPosts(t *testing.T) {
	limit := 2
	listOptions := reddit.ListOptions{
		Limit: limit,
	}

	posts, err := getPosts(subreddit, listOptions)

	if err != nil {
		t.Error(err)
	}

	l := len(posts)

	if l != limit {
		t.Errorf("The length of the retrieved post is wrong: %d", l)
	}
}

func TestGetLatest(t *testing.T) {
	getLatestHelperTest(t, 2)
}

func TestGetPrevious(t *testing.T) {
	posts := getLatestHelperTest(t, 2)
	first := posts[0]
	second := posts[1]

	beforePosts, err := getPreviousPost(subreddit, 1, second.FullID)

	if err != nil {
		t.Error(err)
	}

	l := len(beforePosts)

	if l != 1 {
		t.Errorf("The length of the retrieved post is wrong: %d", l)
	}

	if beforePosts[0].FullID != first.FullID {
		t.Errorf("Fail on GetPrevious")
	}
}

func TestGetNext(t *testing.T) {
	posts := getLatestHelperTest(t, 2)
	first := posts[0]
	second := posts[1]

	afterPosts, err := getLatestPosts(subreddit, 1, first.FullID)

	if err != nil {
		t.Error(err)
	}

	l := len(afterPosts)

	if l != 1 {
		t.Errorf("The length of the retrieved post is wrong: %d", l)
	}

	if afterPosts[0].FullID != second.FullID {
		t.Errorf("Fail on GetNext")
	}
}

func getLatestHelperTest(t *testing.T, limit int) []*reddit.Post {
	posts, err := getLatestPosts(subreddit, limit, "")

	if err != nil {
		t.Error(err)
	}

	l := len(posts)

	if l != limit {
		t.Errorf("The length of the retrieved post is wrong: %d", l)
	}

	return posts
}
