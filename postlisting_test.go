package redditscraper

import (
	"fmt"
	"testing"
)

func TestLimitParam(t *testing.T) {
	l := PostListing{}

	param := l.limitParam()

	if param != fmt.Sprintf("limit=%v", maxPosts) {
		t.Errorf("Invalid limit: %v", l.Limit)
	}

	l.Limit = 1

	param = l.limitParam()

	if param != "limit=1" {
		t.Errorf("Invalid limit: %v", l.Limit)
	}
}

func TestCount(t *testing.T) {
	l := PostListing{
		Count: 10,
	}

	param := l.countParam()

	if param != "count=10" {
		t.Errorf("Invalid count: %v", l.Count)
	}
}

func TestPeriod(t *testing.T) {
	l := PostListing{}

	param := l.periodParam()

	if param != "t=all" {
		t.Errorf("Invalid period: %v", l.Period)
	}

	for _, item := range validPeriods {
		l.Period = item

		param = l.periodParam()

		if param != "t="+item {
			t.Errorf("Invalid period: %v", l.Period)
		}
	}
}

func TestId(t *testing.T) {
	l := PostListing{}

	param := l.idParam()

	if param != "" {
		t.Errorf("Invalid id: %v", l.Id)
	}

	l.Id = "asd"

	param = l.idParam()

	if param != "after=asd" {
		t.Errorf("Invalid id: %v", l.Id)
	}

	l.Latest = true
	param = l.idParam()

	if param != "before=asd" {
		t.Errorf("Invalid id: %v", l.Id)
	}
}

func TestGetUrl(t *testing.T) {
	l := PostListing{
		Show:            true,
		SubredditDetail: true,
	}

	url := l.getUrl("AmItheAsshole", subredditTop)

	if url != "AmItheAsshole/top.json?json_raw=1&limit=100&t=all&show=all&sr_details=" {
		t.Errorf("Invalid url: %v", url)
	}
}
