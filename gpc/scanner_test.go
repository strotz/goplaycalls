package gpc

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
`)
		s := newScanner(r)
		fn := lexIgnore(s)
		assert.Equal(t, eof, s.read())
		assert.Equal(t, "", s.currentValue)
		expected := reflect.ValueOf(lexDetectRequest).Pointer()
		pointer := reflect.ValueOf(fn).Pointer()
		assert.Equal(t, expected, pointer)
	})

	t.Run("detect request", func(t *testing.T) {
		r := strings.NewReader(`### Make a request`)
		s := newScanner(r)
		fn := lexDetectRequest(s)
		assert.Equal(t, "", s.currentValue)
		expected := reflect.ValueOf(lexRequestSeparator).Pointer()
		pointer := reflect.ValueOf(fn).Pointer()
		assert.Equal(t, expected, pointer)
	})
}

//
//func TestScan(t *testing.T) {
//	t.Run("scan empty file", func(t *testing.T) {
//		r := strings.NewReader(``)
//		s := newScanner(r)
//		it, err := s.scan()
//		assert.Equal(t, tokenEOF, it.tok)
//		assert.Equal(t, "", it.val)
//		assert.NoError(t, err)
//	})
//
//	t.Run("scan typical get", func(t *testing.T) {
//		r := strings.NewReader(`### Get operation
//GET https://example.com
//		`)
//		s := newScanner(r)
//		seq, err := s.scanAll()
//		assert.NoError(t, err)
//		expected := []item{
//			{
//				tok: tokenComment,
//				val: "Get operation",
//			},
//			{
//				tok: tokenVerb,
//				val: "GET",
//			},
//			{
//				tok: tokenURL,
//				val: "https://example.com",
//			},
//			{
//				tok: tokenEOF,
//				val: "",
//			},
//		}
//		assert.EqualValues(t, expected, seq)
//	})
//}
