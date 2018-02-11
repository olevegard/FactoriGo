package main

import (
	"testing"

	"github.com/stvp/assert"
)

func TestThatMinReturnsMinimum(t *testing.T) {
	assert.Equal(t, 1, min(1, 2))
	assert.Equal(t, 1, min(2, 1))
	assert.Equal(t, 1, min(1, 1))
	assert.Equal(t, 0, min(0, 1))
	assert.Equal(t, 0, min(0, 1))
	assert.Equal(t, -1, min(0, -1))
}
