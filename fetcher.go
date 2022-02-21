package monzo

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
)

var (
	// ErrPageUnavailable means the page could not be retrieved
	ErrPageUnavailable = errors.New("page unavailable")
)

// Fetcher is a component that retrieves URLs
// and downloads the pages
type Fetcher struct {
	getURL   func() <-chan URL
	complete func(Page) error

	client Client
}

// NewFetcher creates a new Fetcher
func NewFetcher(
	getURL func() <-chan URL,
	complete func(Page) error,
	client Client,
) *Fetcher {
	return &Fetcher{getURL, complete, client}
}

// Crawl blocks and crawls pages from the queue
func (f *Fetcher) Crawl(ctx context.Context) error {
	for {
		var url URL
		select {
		case <-ctx.Done():
			return nil
		case url = <-f.getURL():
		}

		err := f.fetch(ctx, url)
		if err != nil {
			log.Printf("Page failed to crawl: %s, %s", url.String(), err)
			continue
		}

		log.Printf("Crawled page: %s", url.String())
	}
}

func (f *Fetcher) fetch(ctx context.Context, url URL) error {
	content, err := f.client.Get(url)
	if err != nil {
		return err
	}

	page := Page{url, content}
	return f.complete(page)
}

// Client converts a URL into page content
type Client interface {
	// Get converts a URL into a string
	Get(URL) (string, error)
}

// HTTPClient implements Client using Htt[p]
type HTTPClient struct {
	client http.Client
}

var _ Client = (*HTTPClient)(nil)

// NewHTTPClient returns a HttpClient
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

// Get implements the Get Method
func (c *HTTPClient) Get(url URL) (string, error) {
	resp, err := c.client.Get(url.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ErrPageUnavailable
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	return string(bodyBytes), err
}

// FakeClient implements Client for testing purposes
type FakeClient struct {
	Response string
	Err      error
}

var _ Client = (*FakeClient)(nil)

// Get returns the pre-set client's responses
func (c *FakeClient) Get(url URL) (string, error) {
	return c.Response, c.Err
}
