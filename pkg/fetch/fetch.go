package fetch

import (
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	TIMEOUT = 10
)

type NewsSource struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type NewsFeed struct {
	Source *NewsSource
	Feed   *gofeed.Feed
}

// Fetch fetches the feed from the source and returns the Feed. Note that some sources may respond with 403 (Forbidden) due to strict bot policy.
func (source NewsSource) Fetch() (*NewsFeed, error) {
	fp := gofeed.NewParser()
	var feed *gofeed.Feed
	var err error
	succeed := make(chan struct{})
	// Fetch with timeout
	go func() {
		feed, err = fp.ParseURL(source.URL)
		succeed <- struct{}{}
	}()

	select {
	case <-time.After(TIMEOUT * time.Second):
		return nil, fmt.Errorf("fetch: timeout source: %s; %w", source.Name, err)
	case <-succeed:
		return &NewsFeed{Source: &source, Feed: feed}, nil
	}
}

func FailedFetchHandler(source NewsSource) {
	log.Println("Failed to fetch from source: ", source.Name)
}

func SuccessFetchHandler(source NewsSource, feed *gofeed.Feed) {
	log.Println("Successfully fetched from source: ", source.Name)
	log.Println("Feed: ", feed)
}
