package gpc

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type script struct {
	file    string
	content string
}

type step struct {
	name            string
	method          string
	url             string
	responseHandler *script
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
	currentStep := step{}
	var currentHandler *script = nil
	for _, item := range s.items {
		switch item.tok {
		case tokenRequestSeparator:
			currentStep.name = item.val
			if currentStep.valid() {
				res = append(res, currentStep)
				currentStep = step{}
				currentHandler = nil
			}
		case tokenVerb:
			if currentStep.method != "" {
				return nil, errors.New("request separator is missing (verb)")
			} else {
				currentStep.method = item.val
			}
		case tokenURL:
			if currentStep.method == "" {
				return nil, errors.New("method is missing")
			}
			if currentStep.url != "" {
				return nil, errors.New("request separator is missing (url)")
			} else {
				currentStep.url = item.val
			}
		case tokenResponseHandler:
			if !currentStep.valid() {
				return nil, errors.New("failed to declare response handler for invalid request")
			}
			currentStep.responseHandler = &script{}
			currentHandler = currentStep.responseHandler
		case tokenScriptFile:
			if currentHandler == nil {
				return nil, errors.New("missing handler context")
			}
			// TODO: is it a good time to check file?
			currentHandler.file = item.val
		case tokenEmbeddedScript:
			if currentHandler == nil {
				return nil, errors.New("missing handler context")
			}
			if !strings.HasPrefix(item.val, scriptStart) || !strings.HasSuffix(item.val, scriptEnd) {
				return nil, errors.New("invalid script")
			}
			currentHandler.content = strings.TrimSuffix(strings.TrimPrefix(item.val, scriptStart), scriptEnd)
		default:
			return nil, fmt.Errorf("unexpected token: %v - %v", item.tok, item.val)
		}
	}
	if currentStep.valid() {
		res = append(res, currentStep)
		currentStep = step{}
		currentHandler = nil
	}
	return res, nil
}
