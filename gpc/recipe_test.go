package gpc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeRecipe(t *testing.T) {
	t.Run("simple GET", func(t *testing.T) {
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
	})

	t.Run("GET with request handler", func(t *testing.T) {
		r := strings.NewReader(`### call example.com
GET example.com

> {%
console.log("Hello")
%}`)
		steps, err := makeRecipe(r)
		assert.NoError(t, err)
		require.Len(t, steps, 1)
		assert.Equal(t, step{
			name:   "call example.com",
			method: "GET",
			url:    "example.com",
			responseHandler: &script{
				content: `
console.log("Hello")
`,
			},
		}, steps[0])
	})
}
