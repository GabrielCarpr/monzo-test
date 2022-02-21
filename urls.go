package monzo

import (
	"context"
	"sync"
	"time"
)

func NewURL(url string) URL {
	return URL{url}
}

// URL is a target to crawl
type URL struct {
	string
}

func (url URL) String() string {
	return url.string
}

// URLQueue queues URLs for crawling
type URLQueue struct {
	queue []URL // The queue of URLs to process
	front chan URL

	config Config // The system config
	lock   sync.Mutex
}

// NewURLQueue returns a new URLQueue
func NewURLQueue(c Config) *URLQueue {
	return &URLQueue{make([]URL, 0), make(chan URL, 1), c, sync.Mutex{}}
}

// Run runs the queue, blocking.
// Enforces politeness
func (q *URLQueue) Run(ctx context.Context) error {
	for {
		err := q.frontQueue(ctx)
		if err != nil {
			return err
		}

		select {
		case <-time.After(time.Millisecond * time.Duration(q.config.PolitenessInterval)):
			break
		case <-ctx.Done():
			close(q.front)
			return ctx.Err()
		}
	}

}

func (q *URLQueue) frontQueue(ctx context.Context) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	var next URL
	var queue []URL
	switch len(q.queue) {
	case 0:
		return nil
	case 1:
		next = q.queue[0]
		queue = []URL{}
		break
	default:
		next, queue = q.queue[0], q.queue[1:]
	}

	q.queue = queue
	select {
	case q.front <- next:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Queue adds a URL to the queue to be crawled
func (q *URLQueue) Queue(url URL) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue = append(q.queue, url)
}

// Dequeue removes a URL to be crawled
func (q *URLQueue) Dequeue() <-chan URL {
	return q.front
}
