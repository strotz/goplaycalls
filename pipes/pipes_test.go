package pipes

import (
	"log"
	"net/http"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerClose(t *testing.T) {
	pipeName := path.Join(t.TempDir(), t.Name())
	l, err := CreateListener(pipeName)
	assert.NoError(t, err)
	s := http.Server{}

	var wg sync.WaitGroup

	// Listen for request
	wg.Add(1)
	go func() {
		wg.Done()
		err := s.Serve(l)
		require.ErrorIs(t, err, http.ErrServerClosed)
	}()

	wg.Wait()
	assert.NoError(t, s.Close())
}

type handler struct {
}

func (h handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("Serve request", request.URL)
}

func TestServerRequest(t *testing.T) {
	pipeName := path.Join(t.TempDir(), t.Name())

	l, err := CreateListener(pipeName)
	assert.NoError(t, err)
	d := CreateDialer(pipeName)

	s := http.Server{
		Handler: &handler{},
	}

	c := http.Client{
		Transport: &http.Transport{
			DialContext: d,
		},
	}

	var wg sync.WaitGroup

	// Listen for request
	wg.Add(1)
	go func() {
		wg.Done()
		err := s.Serve(l)
		require.ErrorIs(t, err, http.ErrServerClosed)
	}()

	wg.Wait()

	for i := 0; i < 10; i++ {
		log.Println("GET ", i)
		_, err := c.Get("http://pipe/")
		assert.NoError(t, err)
	}
	assert.NoError(t, s.Close())
}
