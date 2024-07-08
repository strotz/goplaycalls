package gpc

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Player struct {
	steps  []step
	report Report
}

type execStep struct {
	step step // Definition of a step.
	req  *http.Request
	res  *http.Response
	rho  string // The output of response handler.
}

func (e execStep) ResponseHandlerOutput() string {
	return e.rho
}

type Report struct {
	steps []execStep
}

func (r Report) Steps() []execStep {
	return r.steps
}

func (r Report) Passed() bool {
	return true
}

func (r Report) LastError() error {
	return nil
}

func (p *Player) Play() (Report, error) {
	report := Report{}
	cl := &http.Client{}
	for _, step := range p.steps {
		item := execStep{
			step: step,
		}
		u, err := url.Parse(step.url)
		if err != nil {
			return report, err
		}
		item.req = &http.Request{
			Method: step.method,
			URL:    u,
		}
		item.res, err = cl.Do(item.req)
		if err != nil {
			return report, err
		}
		if step.responseHandler != nil {
			r := results{}
			output, err := executeResponseHandler(step.responseHandler.content, nil, *item.res, &r)
			if err != nil {
				// TODO: add item to report?
				return report, err
			}
			item.rho = output
		}
		report.steps = append(report.steps, item)
	}
	return report, nil
}

// ParseFile creates a new Player for http request file.
func ParseFile(filePath string) (*Player, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return newPlayer(bufio.NewReader(f))
}

func ParseString(data string) (*Player, error) {
	r := strings.NewReader(data)
	return newPlayer(r)
}

func newPlayer(r io.Reader) (*Player, error) {
	steps, err := makeRecipe(r)
	if err != nil {
		return nil, err
	}
	return &Player{
		steps: steps,
	}, nil
}
