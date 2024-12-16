package redditscraper

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var waitTime time.Duration = 1000000000
var nextRequestWait time.Duration = 5000000000
var nextRequest time.Time = time.Now()
var postsCache = map[string]*CachedPosts{}

type RedditClient struct {
	ri  *RedditInfo
	subreddit string
}

func (c RedditClient) New(subreddit string, ri *RedditInfo) (*RedditClient, error) {
	err := c.refreshToken(ri)

	if err != nil {
		return nil, err
	}

	client := RedditClient{
		ri:  ri,
		subreddit: subreddit,
	}

	return &client, nil
}

func (c RedditClient) FromEnv(subreddit string) (*RedditClient, error) {
	ri, err := RedditInfo{}.FromEnv()

	if err != nil {
		return nil, err
	}

	return c.New(subreddit, ri)
}

func (c *RedditClient) GetBestPosts(listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subredditBest, subredditCacheLong)
}

func (c *RedditClient) GetNewPosts(listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subredditNew, subredditCache1Short)
}

func (c *RedditClient) GetRandomPosts(listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subredditRandom, subredditCache1Short)
}

func (c *RedditClient) GetRisingPosts(listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subredditRising, subredditCacheLong)
}

func (c *RedditClient) GetTopPosts(listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subredditTop, subredditCacheLong)
}

func (c *RedditClient) GetControversialPosts(listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subredditControversial, subredditCacheLong)
}

func (c *RedditClient) getPosts(listing PostListing, sort, duration string) (*CachedPosts, error) {
	url := listing.getUrl(c.subreddit, sort)

	cachedPosts, err := c.getCachedPosts(url)

	if err != nil {
		return nil, err
	}

	if cachedPosts != nil {
		return cachedPosts, nil
	}

	posts, err := c.getPostsHelper(url)

	if err != nil {
		return nil, err
	}

	return c.cachePosts(url, duration, posts)
}

func (c *RedditClient) getCachedPosts(url string) (*CachedPosts, error) {
	cachedPosts := postsCache[url]

	if cachedPosts == nil {
		return nil, nil
	}

	now := time.Now()

	if now.After(cachedPosts.expiresAt) {
		postsCache[url] = nil

		return nil, nil
	}

	return cachedPosts, nil
}

func (c *RedditClient) getPostsHelper(url string) ([]*Post, error) {
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

	request.Header.Set("Authentication", "Bearer "+currentToken.accessToken )

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == 429 {
		c.ri.sleep()
		return c.getPostsHelper(url)
	}

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	return c.convertResponseToPosts(response)
}

func (c *RedditClient) refreshToken(ri *RedditInfo) error {
	if currentToken != nil && currentToken.expiresAt.After(time.Now()) {
		return nil
	}
	
	c.wait()

	t, err := ri.getToken()

	if err != nil {
		return err
	}

	currentToken = t

	return nil
}

func (c *RedditClient) wait() {
	for ; nextRequest.After(time.Now()); {
		time.Sleep(waitTime)
	}
	
	nextRequest = time.Now().Add(nextRequestWait)
}

func (c *RedditClient) convertResponseToPosts(response *http.Response) ([]*Post, error) {
	var body PostListingResponse
	posts := []*Post{}

	defer response.Body.Close()

	err := json.NewDecoder(response.Body).Decode(&body)

	if err != nil {
		return nil, err
	}

	for _, item := range body.Data.Children {
		posts = append(posts, item.Data)
	}

	return posts, nil
}

func (c *RedditClient) cachePosts(url, duration string, posts []*Post) (*CachedPosts, error) {
	d, err := time.ParseDuration(duration)

	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(d)

	cachedPosts := &CachedPosts{
		posts:     posts,
		timestamp: time.Now(),
		expiresAt: expiresAt,
	}

	postsCache[url] = cachedPosts

	return cachedPosts, nil
}