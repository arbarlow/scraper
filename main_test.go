package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	urls, err := fetchBaseDocument(BaseURL)

	assert.Nil(t, err)
	assert.Equal(t, len(urls), 7)

	for _, url := range urls {
		res := make(chan Product)
		errs := make(chan error)
		go fetchProduct(url, res, errs) // start a goroutine
		select {
		case product := <-res:
			assert.NotNil(t, product.Title)
		case err := <-errs:
			assert.Fail(t, err.Error())
		}
	}
}
