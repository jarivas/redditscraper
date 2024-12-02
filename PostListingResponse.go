package redditscraper

type PostListingChild struct {
	Data *Post `json:"data"`
}

type PostListingData struct {
	Children []PostListingChild `json:"children"`
}

type PostListingResponse struct {
	Data PostListingData `json:"data"`
}
