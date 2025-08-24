package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_AddAndGet(t *testing.T) {
	c := New[string, int](3)

	c.Add("a", 1)
	val, ok := c.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, val)
}

func TestCache_Eviction(t *testing.T) {
	c := New[string, int](2)

	c.Add("a", 1)
	c.Add("b", 2)
	c.Add("c", 3)

	_, ok := c.Get("a")
	assert.False(t, ok)

	_, ok = c.Get("b")
	assert.True(t, ok)

	_, ok = c.Get("c")
	assert.True(t, ok)
}

func TestCache_Clear(t *testing.T) {
	c := New[string, int](3)

	c.Add("a", 1)
	c.Clear()

	_, ok := c.Get("a")
	assert.False(t, ok)
}

func TestCache_DuplicateKey(t *testing.T) {
	c := New[string, int](3)

	c.Add("a", 1)
	c.Add("a", 2)

	val, ok := c.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, val)
}
