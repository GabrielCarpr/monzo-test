package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GabrielCarpr/monzo"
	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Seed URL not provided")
		log.Fatal()
	}

	seed := os.Args[1]
	log.Printf("Target: %s", seed)

	crawler := Build(monzo.GetConfig())
	err := crawler.Crawl(monzo.NewURL(seed))
	if err != nil {
		log.Printf("Crawl failed: %s", err)
	}
}

func Build(config monzo.Config) *System {
	urls := monzo.NewURLQueue(config)
	pages := monzo.NewPageQueue()
	fetcher := monzo.NewFetcher(
		func() <-chan monzo.URL { return urls.Dequeue() },
		func(p monzo.Page) error { return pages.Publish(p) },
		monzo.NewHTTPClient(),
	)

	return &System{urls, fetcher, pages, config}
}

type System struct {
	urls    *monzo.URLQueue
	fetcher *monzo.Fetcher
	pages   *monzo.PageQueue

	config monzo.Config
}

func (s *System) Crawl(seed monzo.URL) error {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	s.urls.Queue(seed)

	for i := 0; i < int(s.config.ConcurrentFetchers); i++ {
		g.Go(func() error {
			return s.fetcher.Crawl(ctx)
		})
	}
	g.Go(func() error {
		return s.urls.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	log.Print("Crawl complete")
	return nil
}
