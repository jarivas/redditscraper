package redditscraper

type Post struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Body  string `json:"selftext"`
}