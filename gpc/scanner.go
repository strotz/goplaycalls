package gpc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// Implement parsing of .http files
// It is somewhat inspired by

type token int

// Ignore
// RequestSeparator Comment
// Verb URL
// EOF

const (
	tokenError token = iota
	tokenEOF

	tokenVerb
	tokenComment

	tokenURL
)

const requestSeparator = "###"

type item struct {
	tok token
	val string
}

func (i item) String() string {
	switch i.tok {
	case tokenError:
		return fmt.Sprintf("error: %v", i.val)
	case tokenEOF:
		return "EOF"
	default:
		return i.val
	}
}

//// lexer holds the state of the scanner
//type lexer struct {
//	name  string    // used only for error reports.
//	input string    // the string being scanned.
//	start int       // start position of the item.
//	pos   int       // current position in the input.
//	width int       // width of the last rune read.
//	items chan item // channel of scanned items.
//}
//
//func lex(name, input string) *lexer {
//	l := &lexer{
//		name:  name,
//		input: input,
//		items: make(chan item),
//	}
//	go l.run()
//	return l // l.items
//}
//
//func (l *lexer) run() {
//	for state := lexText; state != nil; {
//		state = state(l)
//	}
//	close(l.items) // no more tokens
//}
//
//func (l *lexer) emit(t token) {
//	l.items <- item{t, l.input[l.start:l.pos]}
//	l.start = l.pos
//}
//
//
//func lexText(l *lexer) stateFn {
//	for {
//		if strings.HasPrefix(l.input[l.pos:], leftMeta) {
//			if l.pos > l.start {
//				l.emit(tokenText)
//			}
//			return lexLeftMeta // next state
//		}
//		if l.next() == eof {
//			break
//		}
//	}
//	// reach eof
//	if l.pos >= len(l.input) {
//		l.emit(tokenText)
//	}
//	l.emit(tokenEOF)
//	return nil
//}
//
//func lexLeftMeta(l *lexer) stateFn {
//	l.pos += len(leftMeta)
//	l.emit(tokenLeftMeta)
//	return lexInsideAction // inside {{}}
//}

var eof = rune(0)
var eol = rune('\n')

type scanner struct {
	reader       *bufio.Reader
	items        chan item
	currentValue string
}

// stateFn scans input while tracking the lexer state and returns state function that tracks the next state.
type stateFn func(scanner *scanner) stateFn

// newScanner returns a new scanner
func newScanner(reader io.Reader) *scanner {
	return &scanner{
		reader:       bufio.NewReader(reader),
		items:        make(chan item),
		currentValue: "",
	}
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

func (s *scanner) accept(valid string) bool {
	ch := s.read()
	if strings.ContainsRune(valid, ch) {
		s.currentValue += string(ch)
		return true
	}
	if ch != eof {
		s.unread()
	}
	return false
}

// ignoreWhiteSpaces accumulates spaces and EOL
func (s *scanner) ignoreWhiteSpaces() {
	for {
		if !s.accept(" \t\r\n") {
			break
		}
	}
	s.currentValue = ""
}

// lexIgnore skips whitespaces in the beginning or end of the file. Does not emit.
func lexIgnore(s *scanner) stateFn {
	s.ignoreWhiteSpaces()
	return lexDetectRequest
}

// lexDetectRequest detects next token and returns proper stateFn. Does not emit.
func lexDetectRequest(s *scanner) stateFn {
	// looking for request separator or verb
	return nil
}

func lexRequestSeparator(s *scanner) stateFn {
	return nil
}

//func (s *scanner) scan() (item, error) {
//	ch := s.read()
//	log.Println("Char:", string(ch))
//	if ch == readError {
//		return item{
//			tok: tokenError,
//			val: "",
//		}, err
//	}
//	switch ch {
//	case eol:
//		return item{tok: tokenEOL, val: ""}, nil
//	case eof:
//		return item{tok: tokenEOF, val: ""}, nil
//	default:
//		return item{tok: tokenError, val: string(ch)}, fmt.Errorf("unknown token: %v, error: %w", ch, err)
//	}
//}
//
//func (s *scanner) scanAll() ([]item, error) {
//	var res []item
//	for {
//		it, err := s.scan()
//		res = append(res, it)
//		if err != nil || it.tok == tokenError || it.tok == tokenEOF {
//			return res, err
//		}
//	}
//}
