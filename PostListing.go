package redditscraper

import (
	"fmt"
	"slices"
	"strconv"
)

var validPeriods = []string{"hour", "day", "week", "month", "year", "all"}

type PostListing struct {
	Id              string `example:"1h3wrtm"`
	Limit           int    `default:"25"`
	Count           int
	Latest          bool
	Period          string `example:"one of (hour, day, week, month, year, all)"`
	Show            bool
	SubredditDetail bool
}

func (l PostListing) String() string {
	result := l.limitParam()

	result += "&" + l.countParam()

	result += "&" + l.periodParam()

	dummy := l.idParam()

	if dummy != "" {
		result += "&" + dummy
	}

	if l.Show {
		result += "&" + l.showParam()
	}

	if l.SubredditDetail {
		result += "&" + l.subredditDetailsParam()
	}

	return result
}

func (l PostListing) limitParam() string {
	if l.Limit == 0 || l.Limit > maxPosts {
		l.Limit = maxPosts
	}

	return "limit=" + strconv.Itoa(l.Limit)
}

func (l PostListing) countParam() string {
	return "count=" + strconv.Itoa(l.Count)
}

func (l PostListing) periodParam() string {
	if ! slices.Contains(validPeriods, l.Period) {
		l.Period = "all"
	}

	return "t"+l.Period
}

func (l PostListing) idParam() string {
	result := ""

	if l.Id != "" {
		if l.Latest {
			result += "before="
		} else {
			result += "after="
		}
		result += l.Id
	}

	return result
}

func (l PostListing) showParam() string {
	return "show=all"
}

func (l PostListing) subredditDetailsParam() string {
	return "sr_details="
}

func (l PostListing) getUrl(subreddit, sort string) string {
	return fmt.Sprintf("/%v/%v?%v", subreddit, sort, l)
}