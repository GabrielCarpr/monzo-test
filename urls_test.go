package monzo_test

import (
	"context"
	"testing"
	"time"

	"github.com/GabrielCarpr/monzo"
	"github.com/stretchr/testify/assert"
)

func TestQueuesDequeues(t *testing.T) {
	q := monzo.NewURLQueue(monzo.Config{
		PolitenessInterval: 1,
	})

	q.Queue(monzo.NewURL("www.google.com"))

	// Run for 1 ms
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3)
	defer cancel()
	err := q.Run(ctx)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	url := <-q.Dequeue()

	assert.Equal(t, "www.google.com", url.String())
}

func TestQueuesDequeuesMultiple(t *testing.T) {
	q := monzo.NewURLQueue(monzo.Config{
		PolitenessInterval: 1,
	})

	q.Queue(monzo.NewURL("www.google.com"))
	q.Queue(monzo.NewURL("www.monzo.com"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	go func() {
		q.Run(ctx)
	}()

	url := <-q.Dequeue()
	assert.Equal(t, "www.google.com", url.String())

	url = <-q.Dequeue()
	assert.Equal(t, "www.monzo.com", url.String())
}

func TestLimitsRate(t *testing.T) {
	q := monzo.NewURLQueue(monzo.Config{
		PolitenessInterval: 1,
	})

	for i := 0; i < 100; i++ {
		q.Queue(monzo.NewURL("www.monzo.com"))
	}

	// In 4ms, max 5 should be possible
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*4)
	defer cancel()
	go func() {
		q.Run(ctx)
	}()

	var urls []monzo.URL
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case url := <-q.Dequeue():
			urls = append(urls, url)
		}
	}

	assert.LessOrEqual(t, len(urls), 5)
}
