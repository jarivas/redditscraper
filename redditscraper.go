package redditscraper

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"time"
)

var waitTime time.Duration = 1000000000
var nextRequest time.Time = time.Now()

type RedditScraper struct {
	ri        *RedditInfo
	subreddit string
}

var validSorts = []string{
	SubredditBest,
	SubredditControversial,
	SubredditNew,
	SubredditRandom,
	SubredditRising,
	SubredditTop,
}

func (c RedditScraper) New(subreddit string, ri *RedditInfo) (*RedditScraper, error) {
	err := c.refreshToken(ri)

	if err != nil {
		return nil, err
	}

	client := RedditScraper{
		ri:        ri,
		subreddit: subreddit,
	}

	return &client, nil
}

func (c RedditScraper) FromEnv(subreddit string) (*RedditScraper, error) {
	ri, err := RedditInfo{}.FromEnv()

	if err != nil {
		return nil, err
	}

	return c.New(subreddit, ri)
}

func (c *RedditScraper) Listen(sort string, listing PostListing, p chan<- *Post, e chan<- error) {
	if slices.Contains(validSorts, sort) {
		c.getPosts(sort,listing, p, e)
	}

	e <- errors.New("invalid subreddit sort for listening")
}

func (c *RedditScraper) getPosts(sort string, listing PostListing, p chan<- *Post, e chan<- error) {
	lastUrl := ""

	for {
		url := listing.getUrl(c.subreddit, sort)

		if url == lastUrl {
			e <- errors.New("repeated same request url: " + url)
			c.wait()
		}

		posts, err := c.getPostsHelper(url)
		l := len(posts)

		if err != nil {
			e <- err
		}

		if l == 0 {
			e <- errors.New("empty post on:" + url)
		} else {
			c.channelPosts(posts, p)
			listing.Id = posts[l-1].Id
		}
	}
}

func (c *RedditScraper) getPostsHelper(url string) ([]*Post, error) {
	err := c.refreshToken(c.ri)

	if err != nil {
		return nil, err
	}

	c.wait()

	request, err := http.NewRequest(
		"GET",
		apiPostsBaseUrl+url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Authentication", "Bearer "+currentToken.accessToken)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == 429 {
		c.ri.sleep()
		return c.getPostsHelper(url)
	}

	if response.StatusCode != 200 {
		logResponse(response.Body)
		return nil, errors.New(response.Status)
	}

	return c.convertResponseToPosts(response)
}

func (c *RedditScraper) refreshToken(ri *RedditInfo) error {
	if currentToken != nil && currentToken.expiresAt.After(time.Now()) {
		return nil
	}

	t, err := ri.getToken()

	if err != nil {
		return err
	}

	currentToken = t

	return nil
}

func (c *RedditScraper) wait() {
	for nextRequest.After(time.Now()) {
		time.Sleep(waitTime)
	}

	nextRequest = time.Now().Add(c.ri.timeSleep)
}

func (c *RedditScraper) convertResponseToPosts(response *http.Response) ([]*Post, error) {
	var body PostListingResponse
	posts := []*Post{}

	defer response.Body.Close()

	err := json.NewDecoder(response.Body).Decode(&body)

	if err != nil {
		logResponse(response.Body)
		return nil, err
	}

	if len(body.Data.Children) == 0 {
		logResponse(response.Body)
		return nil, errors.New("impossible to convert to posts")
	}

	for _, item := range body.Data.Children {
		if (item.Data.Id != "") {
			posts = append(posts, item.Data)
		}
	}

	return posts, nil
}

func (c *RedditScraper) channelPosts(posts []*Post, p chan<- *Post) {
	if len(posts) == 0 {
		return
	}

	for _, post := range posts {
		p <- post
	}
}
