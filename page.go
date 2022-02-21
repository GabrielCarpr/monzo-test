package monzo

// Page is the result of a HTTP request to a target
type Page struct {
	URL     URL
	Content string
}

// PageQueue queues and fans out downloaded pages
type PageQueue struct {
	subscribers []chan Page
}

// NewPageQueue builds a new PageQueue
func NewPageQueue() *PageQueue {
	return &PageQueue{[]chan Page{}}
}

// Publish publishes a page to the queue
func (q *PageQueue) Publish(page Page) error {
	for _, queue := range q.subscribers {
		queue <- page
	}

	return nil
}

// Subscribe registers a subscriber of the page queue
func (q *PageQueue) Subscribe(queue chan Page) {
	q.subscribers = append(q.subscribers, queue)
}
