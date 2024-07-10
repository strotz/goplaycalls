package gpc

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	resp := http.Response{
		StatusCode: http.StatusOK,
	}
	result := results{}
	output, err := executeResponseHandler("console.log('Hello World', response.status)", nil, resp, &result)
	require.NoError(t, err)
	require.Equal(t, "Hello World 200\n", output.console)
}

func TestClient(t *testing.T) {
	t.Run("content", func(t *testing.T) {
		b, err := os.ReadFile("client.js")
		assert.NoError(t, err)
		assert.Equal(t, string(b), clientSource)
	})

	t.Run("script basic client", func(t *testing.T) {
		resp := http.Response{
			StatusCode: http.StatusOK,
		}
		result := results{}
		output, err := executeResponseHandler("client.log(`Hello ${client.name}`)", nil, resp, &result)
		require.NoError(t, err)
		require.Equal(t, "Hello HTTP Client\n", output.console)
	})

	t.Run("add and run passing test", func(t *testing.T) {
		resp := http.Response{
			StatusCode: http.StatusOK,
		}
		result := results{}
		output, err := executeResponseHandler("client.test('first', function() {client.assert(response.status === 200, \"Response status is not 200\");})", nil, resp, &result)
		require.NoError(t, err)
		require.Equal(t, `RUN: first
PASS: first
`, output.console)
	})

	t.Run("add and run failing test", func(t *testing.T) {
		resp := http.Response{
			StatusCode: http.StatusNotFound,
		}
		result := results{}
		output, err := executeResponseHandler("client.test('second', function() {client.assert(response.status === 200, \"Response status is not 200\");})", nil, resp, &result)
		require.NoError(t, err)
		require.Equal(t, `RUN: second
FAILED: second
Error: Response status is not 200
`, output.console)
		require.Len(t, output.failures, 1)
		assert.Equal(t, output.failures[0], "Error: Response status is not 200")
	})
}
