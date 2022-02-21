package monzo_test

import (
	"context"
	"testing"
	"time"

	"github.com/GabrielCarpr/monzo"
	"github.com/stretchr/testify/assert"
)

func TestGetsPages(t *testing.T) {
	pageQueue := make([]monzo.Page, 0)
	urlQueue := make(chan monzo.URL, 2)
	urlQueue <- monzo.NewURL("www.monzo.com")
	urlQueue <- monzo.NewURL("www.google.com")
	queue := func() <-chan monzo.URL {
		return urlQueue
	}
	pages := func(p monzo.Page) error {
		pageQueue = append(pageQueue, p)
		return nil
	}

	f := monzo.NewFetcher(
		queue,
		pages,
		&monzo.FakeClient{"<html>hello</html>", nil},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	assert.NoError(t, f.Crawl(ctx))

	assert.Len(t, pageQueue, 2)
	page := pageQueue[0]
	assert.Contains(t, page.Content, "hello")
	assert.Equal(t, page.URL.String(), "www.monzo.com")
}

func TestPageErrors(t *testing.T) {
	pageQueue := make([]monzo.Page, 0)
	urlQueue := make(chan monzo.URL, 2)
	urlQueue <- monzo.NewURL("www.monzo.com")
	urlQueue <- monzo.NewURL("www.google.com")
	queue := func() <-chan monzo.URL {
		return urlQueue
	}
	pages := func(p monzo.Page) error {
		pageQueue = append(pageQueue, p)
		return nil
	}

	f := monzo.NewFetcher(
		queue,
		pages,
		&monzo.FakeClient{"", monzo.ErrPageUnavailable},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	assert.NoError(t, f.Crawl(ctx))

	assert.Len(t, pageQueue, 0)
}
