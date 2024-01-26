package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionalEqual(t *testing.T) {
	assert.Equal(t, ConditionalEqual(true, 1, 2), 1)
	assert.Equal(t, ConditionalEqual(false, 1, 2), 2)
}

func TestMatch(t *testing.T) {
	assert.True(t, Match("^a", "a"))
	assert.False(t, Match("^a", "0a"))
	assert.False(t, Match("a$", "a0"))
}
