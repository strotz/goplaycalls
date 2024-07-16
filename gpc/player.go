package gpc

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/strotz/goplaycalls/pipes"
)

type Player struct {
	steps  []step
	report Report
	Dialer pipes.DialerFunc
}

type execStep struct {
	step step // Definition of a step.
	req  *http.Request
	res  *http.Response

	rhResult executeResult
}

func (e execStep) ResponseHandlerOutput() string {
	return e.rhResult.console
}

func (e execStep) ResponseHandlerTestErrors() []string {
	return e.rhResult.failures
}

func (e execStep) Failed() bool {
	return len(e.rhResult.failures) > 0
}

type Report struct {
	steps []execStep
}

func (r Report) Steps() []execStep {
	return r.steps
}

func (r Report) TestFailed() bool {
	for _, step := range r.steps {
		if step.Failed() {
			return true
		}
	}
	return false
}

func (p *Player) Play() (Report, error) {
	report := Report{}
	cl := &http.Client{}
	if p.Dialer != nil {
		cl.Transport = &http.Transport{
			DialContext: p.Dialer,
		}
	}
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
			item.rhResult, err = executeResponseHandler(step.responseHandler.content, nil, *item.res, &r)
			if err != nil {
				// TODO: add item to report?
				return report, err
			}
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

func RunTests(filePath string, t *testing.T) Report {
	p, err := ParseFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	report, err := p.Play()
	assert.NoError(t, err)
	// TODO: output in one of the common formats
	// TODO:extract stack trace and make line:pos real
	if report.TestFailed() {
		for _, step := range report.Steps() {
			t.Logf("test console:\n%s", step.ResponseHandlerOutput())
			f := step.ResponseHandlerTestErrors()
			for _, failure := range f {
				t.Logf("failure:\n%s", failure)
			}
		}
		assert.Fail(t, "at least one test failed")
	}
	return report
}
