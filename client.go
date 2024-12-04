package redditscraper

import (
	"encoding/json"
	"net/http"
	"time"
	"errors"
)

type Client struct {
	info  ClientInfo
	token *oauthToken
}

type cachedPosts struct {
	posts     []*Post
	expiresAt time.Time
}

var postsCache = map[string]*cachedPosts{}

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

func (c Client) GetBestPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditBest, subredditCacheLong)
}

func (c Client) GetNewPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditNew, subredditCache1Short)
}

func (c Client) GetRandomPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditRandom, subredditCache1Short)
}

func (c Client) GetRisingPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditRising, subredditCacheLong)
}

func (c Client) GetTopPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditTop, subredditCacheLong)
}

func (c Client) GetControversialPosts(subreddit string, listing PostListing) ([]*Post, error) {
	return c.getPosts(listing, subreddit, subredditControversial, subredditCacheLong)
}

func (c Client) getPosts(listing PostListing, subreddit, sort, duration string) ([]*Post, error) {
	url := listing.getUrl(subreddit, sort)

	posts, err := c.getCachedPosts(url)

	if err != nil {
		writeError(err)
	}

	if posts != nil {
		return posts, nil
	}

	posts, err = c.getPostsHelper(url)

	if err != nil {
		return nil, err
	}

	c.cachePosts(url, duration, posts)

	return posts, nil
}

func (c Client) getCachedPosts(url string) ([]*Post, error) {
	cachedPosts := postsCache[url]

	if cachedPosts == nil {
		return nil, nil
	}

	now := time.Now()

	if now.After(cachedPosts.expiresAt) {
		postsCache[url] = nil

		return nil, nil
	}

	return cachedPosts.posts, nil
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

func (c Client) cachePosts(url, duration string, posts []*Post) {
	d, err := time.ParseDuration(duration)

	if err != nil {
		writeError(err)
		return
	}

	expiresAt := time.Now().Add(d)

	postsCache[url] = &cachedPosts{
		posts:     posts,
		expiresAt: expiresAt,
	}
}
