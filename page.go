package monzo

type Page struct {
	URL     string
	Content string
}

type PageQueue struct {
	subscribers []chan Page
}

func NewPageQueue() *PageQueue {
	return &PageQueue{[]chan Page{}}
}

func (q *PageQueue) Publish(page Page) error {
	for _, queue := range q.subscribers {
		queue <- page
	}

	return nil
}

func (q *PageQueue) Subscribe(queue chan Page) {
	q.subscribers = append(q.subscribers, queue)
}
