package gpc

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertFunc(t *testing.T, expected any, given any) bool {
	ep := reflect.ValueOf(expected).Pointer()
	p := reflect.ValueOf(given).Pointer()
	return assert.Equal(t, ep, p)
}

func TestRead(t *testing.T) {
	t.Run("read empty file", func(t *testing.T) {
		r := strings.NewReader(``)
		s := newScanner(r)
		ch := s.read()
		assert.Equal(t, eof, ch)
	})

	t.Run("peak and read of EOF", func(t *testing.T) {
		r := strings.NewReader(``)
		s := newScanner(r)
		ch := s.peak()
		assert.Equal(t, eof, ch)
		ch = s.read()
		assert.Equal(t, eof, ch)
	})

	t.Run("ignore whitespace", func(t *testing.T) {
		r := strings.NewReader(`


`)
		s := newScanner(r)
		s.ignoreWhiteSpaces()
		ch := s.read()
		assert.Equal(t, eof, ch)
	})

	t.Run("ignore whitespace lexer", func(t *testing.T) {
		r := strings.NewReader(`
x`)
		s := newScanner(r)
		fn := lexIgnore(s)
		assert.Equal(t, "", s.currentValue.String())
		assertFunc(t, lexDetectRequest, fn)
		expected := reflect.ValueOf(lexDetectRequest).Pointer()
		pointer := reflect.ValueOf(fn).Pointer()
		assert.Equal(t, expected, pointer)
	})

	t.Run("ignore whitespace lexer", func(t *testing.T) {
		r := strings.NewReader(``)
		s := newScanner(r)
		fn := lexIgnore(s)
		assert.Nil(t, fn)
		assert.Equal(t, "", s.currentValue.String())
	})

	t.Run("accept word", func(t *testing.T) {
		r := strings.NewReader(`### Make a request`)
		s := newScanner(r)
		s.acceptWord()
		assert.Equal(t, "###", s.currentValue.String())
	})

	t.Run("accept line", func(t *testing.T) {
		r := strings.NewReader(`pick this line
but not this one`)
		s := newScanner(r)
		s.acceptLine()
		assert.Equal(t, "pick this line", s.currentValue.String())
	})

	t.Run("detect request separator", func(t *testing.T) {
		r := strings.NewReader(`### Make a request`)
		s := newScanner(r)
		fn := lexDetectRequest(s)
		assert.Equal(t, "", s.currentValue.String())
		expected := reflect.ValueOf(lexRequestSeparator).Pointer()
		pointer := reflect.ValueOf(fn).Pointer()
		assert.Equal(t, expected, pointer)
	})

	t.Run("detect request verb", func(t *testing.T) {
		r := strings.NewReader(`GET example.com`)
		s := newScanner(r)
		fn := lexDetectRequest(s)
		assert.Equal(t, "", s.currentValue.String())
		expected := reflect.ValueOf(lexRequestUrl).Pointer()
		pointer := reflect.ValueOf(fn).Pointer()
		assert.Equal(t, expected, pointer)
		assert.Len(t, s.items, 1)
		assert.Equal(t, item{
			tok: tokenVerb,
			val: "GET",
		}, s.items[0])
	})

	t.Run("detect wrong verb", func(t *testing.T) {
		r := strings.NewReader(`BLAH example.com`)
		s := newScanner(r)
		fn := lexDetectRequest(s)
		assert.Equal(t, "", s.currentValue.String())
		assert.Nil(t, fn)
		assert.Len(t, s.items, 1)
		assert.Equal(t, item{
			tok: tokenError,
			val: "BLAH",
		}, s.items[0])
	})

	t.Run("detect empty verb", func(t *testing.T) {
		r := strings.NewReader(``)
		s := newScanner(r)
		fn := lexDetectRequest(s)
		assert.Equal(t, "", s.currentValue.String())
		assert.Nil(t, fn)
		assert.Len(t, s.items, 1)
		assert.Equal(t, item{
			tok: tokenError,
			val: "",
		}, s.items[0])
	})

	t.Run("detect request separator comment", func(t *testing.T) {
		r := strings.NewReader(`### Make a request`)
		s := newScanner(r)
		_ = lexDetectRequest(s)
		fn := lexRequestSeparator(s)
		assert.Equal(t, "", s.currentValue.String())
		expected := reflect.ValueOf(lexIgnore).Pointer()
		pointer := reflect.ValueOf(fn).Pointer()
		assert.Equal(t, expected, pointer)
		assert.Len(t, s.items, 1)
		assert.Equal(t, item{
			tok: tokenRequestSeparator,
			val: "Make a request",
		}, s.items[0])
	})

	t.Run("detect request handler embedded", func(t *testing.T) {
		r := strings.NewReader(`> {% console.log("hello") %}`)
		s := newScanner(r)
		fn := lexDetectRequest(s)
		assertFunc(t, lexScript, fn)
		require.Len(t, s.items, 1)
		assert.Equal(t, item{
			tok: tokenResponseHandler,
			val: "",
		}, s.items[0])

		fn = lexScript(s)
		assertFunc(t, lexIgnore, fn)
		require.Len(t, s.items, 2)
		assert.Equal(t, item{
			tok: tokenEmbeddedScript,
			val: "{% console.log(\"hello\") %}",
		}, s.items[1])
	})

	t.Run("extract embedded script", func(t *testing.T) {
		r := strings.NewReader(`> index.js  `)
		s := newScanner(r)
		fn := lexDetectRequest(s)
		assertFunc(t, lexScript, fn)
		require.Len(t, s.items, 1)
		assert.Equal(t, item{
			tok: tokenResponseHandler,
			val: "",
		}, s.items[0])

		fn = lexScript(s)
		assertFunc(t, lexIgnore, fn)
		require.Len(t, s.items, 2)
		assert.Equal(t, item{
			tok: tokenScriptFile,
			val: "index.js",
		}, s.items[1])
	})
}

func TestScan(t *testing.T) {
	t.Run("scan empty file", func(t *testing.T) {
		r := strings.NewReader(``)
		s := newScanner(r)
		s.scan()
		assert.Len(t, s.items, 0)
	})

	t.Run("scan typical get", func(t *testing.T) {
		r := strings.NewReader(`### Get operation
GET https://example.com
		`)
		s := newScanner(r)
		s.scan()
		expected := []item{
			{
				tok: tokenRequestSeparator,
				val: "Get operation",
			},
			{
				tok: tokenVerb,
				val: "GET",
			},
			{
				tok: tokenURL,
				val: "https://example.com",
			},
		}
		assert.EqualValues(t, expected, s.items)
	})

	t.Run("scan typical PUT", func(t *testing.T) {
		r := strings.NewReader(`### Put operation
PUT https://example.com
		`)
		s := newScanner(r)
		s.scan()
		expected := []item{
			{
				tok: tokenRequestSeparator,
				val: "Put operation",
			},
			{
				tok: tokenVerb,
				val: "PUT",
			},
			{
				tok: tokenURL,
				val: "https://example.com",
			},
		}
		assert.EqualValues(t, expected, s.items)
	})

	t.Run("scan typical delete", func(t *testing.T) {
		r := strings.NewReader(`### Delete operation
DELETE https://example.com
		`)
		s := newScanner(r)
		s.scan()
		expected := []item{
			{
				tok: tokenRequestSeparator,
				val: "Delete operation",
			},
			{
				tok: tokenVerb,
				val: "DELETE",
			},
			{
				tok: tokenURL,
				val: "https://example.com",
			},
		}
		assert.EqualValues(t, expected, s.items)
	})

	t.Run("scan get with request handler", func(t *testing.T) {
		r := strings.NewReader(`### Get operation
GET https://example.com

> index.js

		`)
		s := newScanner(r)
		s.scan()
		expected := []item{
			{
				tok: tokenRequestSeparator,
				val: "Get operation",
			},
			{
				tok: tokenVerb,
				val: "GET",
			},
			{
				tok: tokenURL,
				val: "https://example.com",
			},
			{
				tok: tokenResponseHandler,
				val: "",
			},
			{
				tok: tokenScriptFile,
				val: "index.js",
			},
		}
		assert.EqualValues(t, expected, s.items)
	})
}
