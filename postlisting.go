package redditscraper

import (
	"fmt"
	"slices"
	"strconv"
)

var validPeriods = []string{"hour", "day", "week", "month", "year", "all"}

type PostListing struct {
	Id              string `example:"1h3wrtm"`
	Limit           int    `example:"25"`
	Count           int    `example:"25"`
	Latest          bool
	Period          string `example:"one of (hour, day, week, month, year, all)"`
	Show            bool
	SubredditDetail bool
}

func (l PostListing) String() string {
	result := l.limitParam()

	result += "&" + l.periodParam()

	if l.Count > 0 {
		result += "&" + l.countParam()
	}

	if l.Show {
		result += "&" + l.showParam()
	}

	if l.SubredditDetail {
		result += "&" + l.subredditDetailsParam()
	}

	dummy := l.idParam()

	if dummy != "" {
		result += "&" + dummy
	}

	return result
}

func (l PostListing) limitParam() string {
	if l.Limit == 0 || l.Limit > MaxPosts {
		l.Limit = MaxPosts
	}

	return "limit=" + strconv.Itoa(l.Limit)
}

func (l PostListing) countParam() string {
	return "count=" + strconv.Itoa(l.Count)
}

func (l PostListing) periodParam() string {
	if !slices.Contains(validPeriods, l.Period) {
		l.Period = "all"
	}

	return "t=" + l.Period
}

func (l PostListing) idParam() string {
	result := ""

	if l.Id == "" {
		return result
	}

	if l.Latest {
		result += "before="
	} else {
		result += "after="
	}

	result += l.Id

	return result
}

func (l PostListing) showParam() string {
	return "show=all"
}

func (l PostListing) subredditDetailsParam() string {
	return "sr_details="
}

func (l PostListing) getUrl(subreddit, sort string) string {
	return fmt.Sprintf("%v/%v.json?json_raw=1&%v", subreddit, sort, l)
}
