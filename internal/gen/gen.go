package gen

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/drykit-go/strcase"
)

var tplFuncs = template.FuncMap{
	"camelcase": strcase.Camel,
}

type config struct {
	name, tpl, out string
	tplFuncs       template.FuncMap
	data           interface{}
}

// Types generates checkers declarations in packages check and checkconv
// for each type defined in var `checkertypes`. It should be run every time
// that list is modified.
//
// For instance, the following entry:
// 	{N: "Int", T: "int"},
// generates the following declarations:
//
// 	// file check/check.go
// 	type IntPassFunc func(got int) bool
// 	type IntPasser interface { Pass(got int) bool }
// 	type IntChecker interface { IntPasser; Explainer }
//
// 	// file check/checkers.go
// 	type intChecker struct { ... }
// 	func (c intChecker) Pass(got int) bool { ... }
// 	func NewIntChecker(passFunc IntPassFunc, explainFunc ExplainFunc) IntChecker { ... }
//
// 	// file checkconv/assert.go
// 	func FromInt(c check.IntChecker) check.ValueChecker { ... }
// 	// also adds a case in checkconv.Assert:
// 		case check.IntChecker:
// 			return FromInt(c)
func Types(tpl, out string) error {
	return generate(config{
		name:     "types",
		tpl:      tpl,
		out:      out,
		tplFuncs: tplFuncs,
		data:     checkertypes,
	})
}

// Interfaces generates the public interfaces for checker providers
// implementations in package check (provider_*.go files). It should be run
// every time their API is modified (method signature change, doc comment,
// new method, method removal, ...)
func Interfaces(tpl, out string) error {
	data, err := computeInterfaces()
	if err != nil {
		return err
	}
	return generate(config{
		name: "interfaces",
		tpl:  tpl,
		out:  out,
		data: data,
	})
}

func generate(cfg config) error {
	t, err := newTemplate(cfg.name, cfg.tpl, cfg.tplFuncs)
	if err != nil {
		return err
	}

	f, err := createOutFile(cfg.out)
	if err != nil {
		return err
	}

	if err := t.Execute(f, cfg.data); err != nil {
		return err
	}

	return runFormatter(cfg.out)
}

func newTemplate(name, src string, funcMap template.FuncMap) (*template.Template, error) {
	s, err := readTplFile(src)
	if err != nil {
		return nil, err
	}
	return template.New(name).Funcs(funcMap).Parse(s)
}

func readTplFile(filepath string) (string, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", filepath, err)
	}
	return string(b), nil
}

func createOutFile(filepath string) (*os.File, error) {
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	if _, err := f.WriteString(generateHeader()); err != nil {
		return nil, err
	}
	return f, nil
}

func generateHeader() string {
	s := strings.Builder{}
	s.WriteString("// Code generated by go generate ./...; DO NOT EDIT\n")
	s.WriteString("// Last generated on ")
	s.WriteString(time.Now().UTC().Format(time.RFC822))
	s.WriteString("\n\n")
	return s.String()
}

func runFormatter(filepath string) error {
	cmd := exec.Command("goimports", "-w", filepath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"goimports returned an error running on file %s, it is probably malformed",
			filepath,
		)
	}
	return nil
}
