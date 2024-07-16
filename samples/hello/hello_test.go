package hello

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/strotz/goplaycalls/gpc"
	"github.com/strotz/goplaycalls/pipes"
	"github.com/strotz/goplaycalls/testserver"
)

// Checks HTTP GET call to the server.
func TestHello(t *testing.T) {
	p, err := gpc.ParseFile("hello.http")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("call GET hello without server", func(t *testing.T) {
		_, err := p.Play()
		require.ErrorContains(t, err,
			"connect: connection refused")
	})

	t.Run("call GET hello", func(t *testing.T) {
		ts := testserver.NewTestServer(t.Name(), http.MethodGet, "/hello", Handler)
		ts.Start()
		t.Cleanup(ts.Stop)

		p.Dialer = pipes.CreateDialer(t.Name())
		r, err := p.Play()
		assert.NoError(t, err)
		assert.True(t, r.TestFailed())

		steps := r.Steps()
		require.Len(t, steps, 1)

		assert.Equal(t, `hello status: 200
name: {"name":"Double Belomor"}
RUN: Request executed successfully
PASS: Request executed successfully
RUN: Failed test
FAILED: Failed test
Error: Name has to be Hello, but got Double Belomor
`, steps[0].ResponseHandlerOutput())

		assert.True(t, steps[0].Failed())
		assert.Equal(t, []string{
			"Error: Name has to be Hello, but got Double Belomor",
		}, steps[0].ResponseHandlerTestErrors())
	})
}
