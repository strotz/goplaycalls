package gpc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMakeRecipe(t *testing.T) {
	r := strings.NewReader(`### call example.com
GET example.com`)
	steps, err := makeRecipe(r)
	assert.NoError(t, err)
	require.Len(t, steps, 1)
	assert.Equal(t, step{
		name:   "call example.com",
		method: "GET",
		url:    "example.com",
	}, steps[0])
}
