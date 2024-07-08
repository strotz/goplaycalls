package gpc

import (
	"net/http"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

type playEnvironment struct{}

type results struct {
	value goja.Value
}

type scriptOutput struct {
	content strings.Builder
}

func (s *scriptOutput) Log(message string) {
	s.content.WriteString(message)
}

func (s *scriptOutput) Warn(message string) {
	s.content.WriteString(message)
}

func (s *scriptOutput) Error(message string) {
	s.content.WriteString(message)
}

func (s *scriptOutput) Output() string {
	return s.content.String()
}

type ResponseAdapter struct {
	Status string `json:"status"`
}

func executeResponseHandler(source string, env *playEnvironment, response http.Response, out *results) (string, error) {
	printer := &scriptOutput{}

	registry := new(require.Registry) // this can be shared by multiple runtimes

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	_ = registry.Enable(vm)
	registry.RegisterNativeModule(console.ModuleName, console.RequireWithPrinter(printer))
	console.Enable(vm)

	err := vm.Set("response", ResponseAdapter{
		Status: response.Status,
	})
	if err != nil {
		return "", err
	}

	out.value, err = vm.RunString(source)
	if err != nil {
		return printer.Output(), err
	}
	return printer.Output(), nil
}
