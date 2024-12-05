package redditscraper

type Post struct {
	Id    string `json:"name"`
	Title string `json:"title"`
	Body  string `json:"selftext"`
}
