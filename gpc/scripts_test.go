package gpc

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	resp := http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
	}
	result := results{}
	output, err := executeResponseHandler("console.log('Hello World', response.status)", nil, resp, &result)
	require.NoError(t, err)
	require.Equal(t, "Hello World OK", output)
}
