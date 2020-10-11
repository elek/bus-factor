package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDevPower(t *testing.T) {
	h := NewHistogram()
	for i := 0; i < 10; i++ {
		h.Increment("primary")
	}
	for i := 0; i < 10; i++ {
		h.Increment("secondary")
	}
	assert.Equal(t, float64(2), h.DevPower())
}

func TestDevPower2(t *testing.T) {
	h := NewHistogram()
	for i := 0; i < 10; i++ {
		h.Increment("primary")
	}
	for i := 0; i < 10; i++ {
		h.Increment("secondary")
	}
	for i := 0; i < 5; i++ {
		h.Increment("third")
	}
	assert.Equal(t, float64(2.7071067811865475), h.DevPower())
}
