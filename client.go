package redditscraper

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Client struct {
	info  ClientInfo
	token *oauthToken
}

var postsCache = map[string]*CachedPosts{}

func (c Client) New(info ClientInfo) (*Client, error) {
	token, err := info.getToken()

	if err != nil {
		return nil, err
	}

	client := Client{
		info:  info,
		token: token,
	}

	return &client, nil
}

func (c Client) FromEnv() (*Client, error) {
	info := ClientInfo{}.fromEnv()

	return c.New(info)
}

func (c Client) GetBestPosts(subreddit string, listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subreddit, subredditBest, subredditCacheLong)
}

func (c Client) GetNewPosts(subreddit string, listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subreddit, subredditNew, subredditCache1Short)
}

func (c Client) GetRandomPosts(subreddit string, listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subreddit, subredditRandom, subredditCache1Short)
}

func (c Client) GetRisingPosts(subreddit string, listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subreddit, subredditRising, subredditCacheLong)
}

func (c Client) GetTopPosts(subreddit string, listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subreddit, subredditTop, subredditCacheLong)
}

func (c Client) GetControversialPosts(subreddit string, listing PostListing) (*CachedPosts, error) {
	return c.getPosts(listing, subreddit, subredditControversial, subredditCacheLong)
}

func (c Client) getPosts(listing PostListing, subreddit, sort, duration string) (*CachedPosts, error) {
	url := listing.getUrl(subreddit, sort)

	cachedPosts, err := c.getCachedPosts(url)

	if err != nil {
		writeError(err)

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

func (c Client) getCachedPosts(url string) (*CachedPosts, error) {
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

func (c Client) getPostsHelper(url string) ([]*Post, error) {
	request, err := http.NewRequest(
		"GET",
		apiPostsBaseUrl+url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Authentication", "Bearer "+c.token.accessToken)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	return c.convertResponseToPosts(response)
}

func (c Client) convertResponseToPosts(response *http.Response) ([]*Post, error) {
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

func (c Client) cachePosts(url, duration string, posts []*Post) (*CachedPosts, error) {
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
