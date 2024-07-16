package gpc

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/strotz/goplaycalls/pipes"
	"github.com/strotz/goplaycalls/testserver"
)

func echoHandler(response http.ResponseWriter, req *http.Request) {
	// TODO: it is echo service, need to return something :)
	log.Println("Method:", req.Method)
}

func TestCallGetRequest(t *testing.T) {
	ts := testserver.NewTestServer(t.Name(), http.MethodGet, "/a", echoHandler)
	ts.Start()
	t.Cleanup(ts.Stop)

	p, err := ParseString(`### Get operation
GET http://localhost:8080/a
`)
	p.Dialer = pipes.CreateDialer(t.Name())
	require.NoError(t, err)
	r, err := p.Play()
	assert.NoError(t, err)
	assert.False(t, r.TestFailed())
}

func TestCallPutRequest(t *testing.T) {
	ts := testserver.NewTestServer(t.Name(), http.MethodPut, "/b", echoHandler)
	ts.Start()
	t.Cleanup(ts.Stop)

	p, err := ParseString(`### Put operation
PUT http://localhost:8080/b
`)
	p.Dialer = pipes.CreateDialer(t.Name())
	require.NoError(t, err)
	r, err := p.Play()
	assert.NoError(t, err)
	assert.False(t, r.TestFailed())
}
