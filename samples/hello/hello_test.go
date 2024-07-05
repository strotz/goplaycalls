package hello

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strotz/goplaycalls/gpc"
)

// Checks HTTP GET call to the server.
func TestHello(t *testing.T) {
	p, err := gpc.ParseFile("hello.http")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("call GET hello without server", func(t *testing.T) {
		r := p.Play()
		assert.NotNil(t, r)
		assert.False(t, r.Passed())
		assert.EqualError(t, r.LastError(), "Connection refused: localhost/[0:0:0:0:0:0:0:1]:8080")
	})

	t.Run("call GET hello", func(t *testing.T) {
		// Start test server
		// TODO: make it reusable
		sm := http.NewServeMux()
		sm.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			// TODO: this is empty function, add something interesting for hello
		})
		s := http.Server{
			Addr:    "localhost:8080",
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

		r := p.Play()
		assert.NoError(t, err)
		assert.True(t, r.Passed())
		assert.NoError(t, r.LastError())
		//"HTTP/1.1 200 OK\nDate: Thu, 04 Jul 2024 22:51:18 GMT\nContent-Length: 0\n\n<Response body is empty>\n"
	})
}
