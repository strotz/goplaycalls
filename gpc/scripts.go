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

func executeResponseHandler(source string, env *playEnvironment, response http.Response, out *results) (string, error) {
	printer := &scriptOutput{}

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
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()
		// TODO: it seems like httpclient tool detects json output and formats it. See httputil.DumpResponse(&response, true)
		r.Body = string(body)
	}

	err := vm.Set("response", r)
	if err != nil {
		return "", err
	}

	// TODO: verify proper client initialization, at least no errors
	_, err = vm.RunString(clientSource)
	if err != nil {
		return printer.Output(), err
	}
	out.value, err = vm.RunString(source)
	if err != nil {
		return printer.Output(), err
	}

	_, err = vm.RunString("client.runTests()")
	if err != nil {
		return printer.Output(), err
	}
	return printer.Output(), nil
}
