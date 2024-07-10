package hello

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/strotz/goplaycalls/gpc"
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
		// Start test server
		// TODO: make it reusable
		sm := http.NewServeMux()
		sm.HandleFunc("/hello", Handler)
		s := http.Server{
			Addr:    ":8080",
			Handler: sm,
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Done()
			err := s.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				require.NoError(t, err)
			}
		}()
		wg.Wait()
		defer func() {
			err := s.Shutdown(context.Background())
			require.NoError(t, err)
		}()

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
