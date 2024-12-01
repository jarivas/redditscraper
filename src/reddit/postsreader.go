package redditscraper

import (
	"context"
	"errors"
	"os"
	"strconv"
)

const MAX_POSTS = 25

func getPosts(subreddit string, listOptions reddit.ListOptions) ([]Post, error) {
	options := getListPostOptions(listOptions)

	return getTopPosts(subreddit, &options)
}

func getPreviousPost(subreddit string, limit int, fullID string) ([]Post, error) {
	options, err := getListPostOptionsEnv(limit)

	if err != nil {
		return nil, err
	}

	options.Before = fullID

	return getTopPosts(subreddit, &options)
}

func getLatestPosts(subreddit string, limit int, fullID string) ([]Post, error) {
	options, err := getListPostOptionsEnv(limit)

	if err != nil {
		return nil, err
	}

	options.After = fullID

	return getTopPosts(subreddit, &options)
}

func getTopPosts(subreddit string, options *reddit.ListPostOptions) ([]Post, error) {
	ctx := context.Background()
	client := reddit.DefaultClient()

	posts, _, err := client.Subreddit.TopPosts(ctx, subreddit, options)

	if err != nil {
		return nil, err
	}

	return convertRedditPosts(posts)
}

func getListPostOptions(listOptions reddit.ListOptions) reddit.ListPostOptions {
	return reddit.ListPostOptions{
		ListOptions: listOptions,
		Time:        "all",
	}
}

func getListPostOptionsEnv(limit int) (reddit.ListPostOptions, error) {
	listOptions := reddit.ListOptions{
		Limit: limit,
	}
	options := getListPostOptions(listOptions)
	err := godotenv.Load()

	if err != nil {
		return options, errors.New("error loading .env file")
	}

	err = getListPostOptionsLimitEnvHelper(&listOptions)

	if err != nil {
		return options, err
	}

	return options, nil
}

func getListPostOptionsLimitEnvHelper(listOptions *reddit.ListOptions) error {
	if listOptions.Limit > 0 {
		return nil
	}

	maxPost, ok := os.LookupEnv("REDDIT_MAX_POSTS")

	if !ok {
		return errors.New("REDDIT_MAX_POSTS not present on .env file or env var")
	}

	limit, err := strconv.Atoi(maxPost)

	if err != nil {
		return errors.New("REDDIT_MAX_POSTS is not an integer")
	}

	if limit < 0 || limit > MAX_POSTS {
		return errors.New("REDDIT_MAX_POSTS value is invalid")
	}

	listOptions.Limit = limit

	return nil
}

func convertRedditPosts(posts []*reddit.Post) ([]Post, error) {
	result := []Post{}

	for index, post := range posts {
		result[index] = convertRedditPost(post)		
	}

	return result
}

func convertRedditPost(post *reddit.Post) Post {
	return Post{
		FullID: post.FullID,
		Title:  post.Title,
		Body:   post.Body,
	}
}
