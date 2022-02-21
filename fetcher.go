package monzo

type Fetcher struct {
	getURL func() <-chan URL

	complete func(Page) error
}

func NewFetcher(getURL func() <-chan URL, complete func(Page) error) *Fetcher {
	return &Fetcher{getURL, complete}
}
