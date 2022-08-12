package fetch

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT*time.Second)
	defer cancel()
	feed, err = fp.ParseURLWithContext(source.URL, ctx)
	if err != nil {
		return nil, err
	}
	return &NewsFeed{Source: &source, Feed: feed}, nil
}
