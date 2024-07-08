package gpc

import (
	"errors"
	"fmt"
	"io"
)

type step struct {
	name   string
	method string
	url    string
}

func (s step) valid() bool {
	return s.method != "" || s.url != ""
}

func makeRecipe(reader io.Reader) ([]step, error) {
	s := newScanner(reader)
	s.scan()
	for _, item := range s.items {
		if item.tok == tokenError {
			return nil, fmt.Errorf("failed to make steps: %v", item.val)
		}
	}
	res := []step{}
	current := step{}
	for _, item := range s.items {
		switch item.tok {
		case tokenRequestSeparator:
			current.name = item.val
			if current.valid() {
				res = append(res, current)
				current = step{}
			}
		case tokenVerb:
			if current.method != "" {
				return nil, errors.New("request separator is missing (verb)")
			} else {
				current.method = item.val
			}
		case tokenURL:
			if current.method == "" {
				return nil, errors.New("method is missing")
			}
			if current.url != "" {
				return nil, errors.New("request separator is missing (url)")
			} else {
				current.url = item.val
			}
		default:
			return nil, fmt.Errorf("unexpected token: %v - %v", item.tok, item.val)
		}
	}
	if current.valid() {
		res = append(res, current)
	}
	return res, nil
}
