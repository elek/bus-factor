package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestDevPower(t *testing.T) {
	h := createHistogram(10, 10)
	assert.Equal(t, float64(2), h.DevPower())
	assert.Equal(t, float64(2), h.DevPower())

	h = createHistogram(10, 10, 5)
	assert.Equal(t, float64(2.25), h.DevPower())


	h = createHistogram(10, 10, 5, 5)
	assert.Equal(t, float64(2.5), h.DevPower())
	assert.Equal(t, float64(3), h.RawDevPower())



}

func createHistogram(shares ...int) Histogram {
	h := NewHistogram()

	for i := 0; i < len(shares); i++ {
		for j := 0; j < shares[i]; j++ {
			h.Increment("user" + strconv.Itoa(i))
		}
	}
	return h
}
