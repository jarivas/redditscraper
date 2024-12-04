package redditscraper

import (
	"time"
	"encoding/json"
	"net/http"
)

type Client struct {
	info  *ClientInfo
	token *oauthToken
}

type cachedPosts struct {
	posts     []*Post
	expiresAt time.Time
}

var cache = make(map[string]*cachedPosts)

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
	cachedPost := cache[url]

	if cachedPost == nil {
		return nil, nil
	}

	now := time.Now()

	if now.After(cachedPost.expiresAt) {
		cache[url] = nil

		return nil, nil
	}

	return cachedPost.posts, nil
}

func (c Client) getPostsHelper(url string) ([]*Post, error) {
	request, err := http.NewRequest(
		"GET",
		apiBaseUrl + url + "&raw_json=1.json",
		nil,
	)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Authentication", "Bearer " + c.token.accessToken)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	return c.convertResponseToPosts(response)
}

func (c Client) convertResponseToPosts(response *http.Response) ([]*Post, error) {
	var body PostListingResponse
	posts := []*Post{}

	defer response.Body.Close()

	err := json.NewDecoder(response.Body).Decode(body)

	if err != nil {
		return nil, err
	}

	for index, item := range(body.Data.Children) {
		posts[index] = item.Data
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

	cache[url] = &cachedPosts{
		posts: posts,
		expiresAt: expiresAt,
	}
}