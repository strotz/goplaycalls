package gpc

import (
	_ "embed"
	"io"
	"net/http"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

//go:embed client.js
var clientSource string

type playEnvironment struct{}

type results struct {
	value goja.Value
}

type scriptOutput struct {
	content strings.Builder
}

func (s *scriptOutput) Log(message string) {
	s.content.WriteString(message)
	s.content.WriteRune('\n')
}

func (s *scriptOutput) Warn(message string) {
	s.Log(message)
}

func (s *scriptOutput) Error(message string) {
	s.Log(message)
}

func (s *scriptOutput) Output() string {
	return s.content.String()
}

type ResponseAdapter struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

type executeResult struct {
	console  string
	failures []string
}

// executeResponseHandler executes response handler and returns the result of test execution along with console output.
func executeResponseHandler(source string, env *playEnvironment, response http.Response, out *results) (result executeResult, err error) {
	printer := &scriptOutput{}
	defer func() {
		result.console = printer.Output()
	}()

	registry := new(require.Registry) // this can be shared by multiple runtimes

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	_ = registry.Enable(vm)
	registry.RegisterNativeModule(console.ModuleName, console.RequireWithPrinter(printer))
	console.Enable(vm)

	r := ResponseAdapter{
		Status: response.StatusCode,
	}
	if response.Body != nil {
		var body []byte
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return
		}
		defer response.Body.Close()
		// TODO: it seems like httpclient tool detects json output and formats it. See httputil.DumpResponse(&response, true)
		r.Body = string(body)
	}
	err = vm.Set("response", r)
	if err != nil {
		return
	}

	// TODO: verify proper client initialization, at least no errors
	_, err = vm.RunString(clientSource)
	if err != nil {
		return
	}
	out.value, err = vm.RunString(source)
	if err != nil {
		return
	}

	// Run declared tests and process the results.
	testResult, err := vm.RunString("client.runTests()")
	ex := testResult.ToObject(vm)
	for _, key := range ex.Keys() {
		// TODO: it is slightly hacky, need to make strict abstraction
		result.failures = append(result.failures, ex.Get(key).String())
	}
	return
}
