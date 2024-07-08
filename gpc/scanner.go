package gpc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Implement parsing of .http files

// Ignore
// RequestSeparator Comment
// Verb URL
// <empty line?>
// > {% .... %}
// > file.js

const spaceChars = " \t\r\n"
const lineEnds = "\r\n"
const scriptStart = "{%"
const scriptEnd = "%}"

type token int

const (
	tokenError token = iota

	tokenRequestSeparator
	tokenVerb
	tokenURL

	tokenResponseHandler

	tokenEmbeddedScript
	tokenScriptFile
)

const requestSeparator = "###"
const responseHandlerStart = ">"

type item struct {
	tok token
	val string
}

func (i item) String() string {
	switch i.tok {
	case tokenError:
		return fmt.Sprintf("error: %v", i.val)
	default:
		return i.val
	}
}

var eof = rune(0)

type scanner struct {
	reader       *bufio.Reader
	items        []item
	currentValue strings.Builder
}

// stateFn scans input while tracking the lexer state and returns state function that tracks the next state.
type stateFn func(scanner *scanner) stateFn

// newScanner returns a new scanner
func newScanner(reader io.Reader) *scanner {
	return &scanner{
		reader: bufio.NewReader(reader),
	}
}

func (s *scanner) emitItem(it item) {
	s.items = append(s.items, it)
}

// read returns the next rune from the input. It returns rune eof when input is over
// or readError when other error.
func (s *scanner) read() rune {
	r, _, err := s.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return eof
		}
		log.Fatalln("failed to read:", err)
	}
	return r
}

func (s *scanner) unread() {
	err := s.reader.UnreadRune()
	if err != nil {
		log.Fatalln("failed to unread:", err)
	}
}

func (s *scanner) peak() rune {
	r := s.read()
	if r != eof {
		s.unread()
	}
	return r
}

type acceptFn func(r rune) bool

func (s *scanner) accept(fn acceptFn) bool {
	ch := s.read()
	if ch == eof {
		return false
	}
	if fn(ch) {
		s.currentValue.WriteRune(ch)
		return true
	}
	s.unread()
	return false
}

// acceptWhiteSpaces accumulates spaces and EOL
func (s *scanner) acceptWhiteSpaces() {
	for {
		// Collect all white space characters.
		if !s.accept(func(r rune) bool {
			return strings.ContainsRune(spaceChars, r)
		}) {
			break
		}
	}
}

// ignoreWhiteSpaces accumulates spaces and EOL
func (s *scanner) ignoreWhiteSpaces() {
	s.acceptWhiteSpaces()
	s.currentValue.Reset()
}

// acceptWord collects one word of characters in currentValue
func (s *scanner) acceptWord() {
	for {
		// Collect everything, except white spaces.
		if !s.accept(func(r rune) bool {
			return !strings.ContainsRune(spaceChars, r)
		}) {
			break
		}
	}
}

func (s *scanner) acceptLine() {
	for {
		// Collect everything until EOL, including spaces.
		if !s.accept(func(r rune) bool {
			return !strings.ContainsRune(lineEnds, r)
		}) {
			break
		}
	}
}

func (s *scanner) emitError() {
	s.emitItem(item{
		tok: tokenError,
		val: s.currentValue.String(),
	})
	s.currentValue.Reset()
}

// lexIgnore skips whitespaces in the beginning or end of the file. Does not emit.
func lexIgnore(s *scanner) stateFn {
	s.ignoreWhiteSpaces()
	if s.peak() == eof {
		return nil
	}
	return lexDetectRequest
}

// lexDetectRequest detects next token and returns proper stateFn. Could emit an error.
func lexDetectRequest(s *scanner) stateFn {
	// looking for request separator or verb
	s.acceptWord()
	switch s.currentValue.String() {
	case requestSeparator:
		s.currentValue.Reset()
		return lexRequestSeparator
	case http.MethodGet:
		s.emitItem(item{
			tok: tokenVerb,
			val: http.MethodGet,
		})
		s.currentValue.Reset()
		return lexRequestUrl
	case responseHandlerStart:
		s.emitItem(item{
			tok: tokenResponseHandler,
			val: "",
		})
		s.currentValue.Reset()
		return lexScript
		// TODO: it seems like empty line after request has a certain meaning
	}
	s.emitError()
	return nil
}

// lexRequestSeparator detects the name of the request
func lexRequestSeparator(s *scanner) stateFn {
	// collect everything until the end of the line
	s.acceptLine()
	val := strings.TrimSpace(s.currentValue.String())
	if len(val) > 0 {
		s.emitItem(item{
			tok: tokenRequestSeparator,
			val: val,
		})
		s.currentValue.Reset()
	}
	return lexIgnore
}

// lexRequestUrl emits the URL
func lexRequestUrl(s *scanner) stateFn {
	s.ignoreWhiteSpaces()
	s.acceptWord()
	if s.currentValue.Len() > 0 {
		s.emitItem(item{
			tok: tokenURL,
			val: s.currentValue.String(),
		})
		s.currentValue.Reset()
	}
	return lexIgnore
}

// lexScript detects either embedded script or external file
func lexScript(s *scanner) stateFn {
	s.ignoreWhiteSpaces()
	s.acceptWord()
	if !strings.HasPrefix(s.currentValue.String(), scriptStart) {
		// Does not look like embedded script
		s.emitItem(item{
			tok: tokenScriptFile,
			val: s.currentValue.String(),
		})
		s.currentValue.Reset()
		return lexIgnore
	}
	for {
		s.acceptWhiteSpaces()
		s.acceptWord()
		if strings.HasSuffix(s.currentValue.String(), scriptEnd) {
			s.emitItem(item{
				tok: tokenEmbeddedScript,
				val: s.currentValue.String(),
			})
			return lexIgnore
		}
		if s.peak() == eof {
			return nil
		}
	}
}

func (s *scanner) scan() {
	for state := lexIgnore; state != nil; {
		state = state(s)
	}
}
