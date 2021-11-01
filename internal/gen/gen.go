package gen

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

type config struct {
	name, tpl, out string
	src            interface{}
	tplFuncs       template.FuncMap
}

// Interfaces generates the public interfaces for checker providers
// implementations in package check (provider_*.go files). It should be run
// every time their API is modified (method signature change, doc comment,
// new method, method removal, ...)
func Interfaces(tpl, out string) error {
	src, err := computeInterfaces()
	if err != nil {
		return err
	}
	return generate(config{
		name: "interfaces",
		tpl:  tpl,
		out:  out,
		src:  src,
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

	if err := t.Execute(f, cfg.src); err != nil {
		return err
	}

	// FIXME: make goimports work with generics syntax
	if err := runFormatter(cfg.out); err != nil {
		fmt.Printf("goimports returned the following error (likely due to unsupported generics syntax):\n%s\n", err)
	}
	return nil
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
