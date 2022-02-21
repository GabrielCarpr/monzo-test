package monzo_test

import (
	"testing"

	"github.com/GabrielCarpr/monzo"
	"github.com/stretchr/testify/assert"
)

func TestFansOutQueue(t *testing.T) {
	q := monzo.NewPageQueue()

	one := make(chan monzo.Page, 100)
	two := make(chan monzo.Page, 100)
	three := make(chan monzo.Page, 100)

	q.Subscribe(one)
	q.Subscribe(two)
	q.Subscribe(three)

	q.Publish(monzo.Page{
		URL:     monzo.NewURL("www.monzo.com"),
		Content: "test",
	})

	page1 := <-one
	assert.Equal(t, "www.monzo.com", page1.URL.String())
	page2 := <-two
	assert.Equal(t, "www.monzo.com", page2.URL.String())
	page3 := <-three
	assert.Equal(t, "www.monzo.com", page3.URL.String())
}
