package testserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/strotz/goplaycalls/pipes"
)

// TestServer uses named pipes
type TestServer struct {
	name string
	l    net.Listener
	s    http.Server
}

// Start test server
func (t *TestServer) Start() {
	var err error
	t.l, err = pipes.CreateListener(t.name)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		err := t.s.Serve(t.l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	wg.Wait()
}

func (t *TestServer) Stop() {
	err := t.s.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

// NewTestServer creates a new http server on Unix pipes that serves only one service.
func NewTestServer(name string, verb string, url string, handler func(http.ResponseWriter, *http.Request)) *TestServer {
	res := &TestServer{
		name: name,
	}
	sm := http.NewServeMux()
	sm.HandleFunc(fmt.Sprintf("%s %s", verb, url), handler)
	res.s = http.Server{
		Handler: sm,
	}
	return res
}
