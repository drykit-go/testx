package gen

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"strings"

	"github.com/drykit-go/slicex"

	"github.com/drykit-go/testx/internal/gen/docparser"
	"github.com/drykit-go/testx/internal/gen/metatype"
)

// ProvidersMetaData is a representation of the parsed doc for providers files
// in package check.
// It is meant to be used as a data source for template providers.gotmpl.
type ProvidersMetaData struct {
	Vars []metatype.Var
}

func computeProvidersMetaData() (ProvidersMetaData, error) {
	docp, err := newDocPackage("check", isProviderFile) // TODO: package-agnostic
	if err != nil {
		return ProvidersMetaData{}, err
	}

	data := ProvidersMetaData{Vars: numericMetaVars()}
	for _, t := range docp.Types {
		data.Vars = append(data.Vars, computeMetaVar(t))
	}

	return data, nil
}

// computeMetaVar returns a MetaVar after the given *doc.Type.
func computeMetaVar(t *doc.Type) metatype.Var {
	varname := structToVarName(t.Name)
	return metatype.Var{
		Name:  varname,
		Value: structInstanceString(t.Name),
		DocLines: docparser.ParseDocLines(t.Doc, map[string]string{
			t.Name: varname,
		}),
	}
}

// newDocPackage returns a *doc.Package matching packageName after applying
// the given filter, or the first non-nil error occurring in the process.
func newDocPackage(packageName string, filter func(fs.FileInfo) bool) (*doc.Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "./", filter, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	astp, ok := pkgs[packageName]
	if !ok {
		return nil, fmt.Errorf(
			"no files found for package %s, check path or filters",
			packageName,
		)
	}
	return doc.New(astp, "./", doc.AllDecls), nil
}

func numericMetaVars() []metatype.Var {
	type namedType struct{ N, T string }
	numerics := []namedType{
		{N: "Int", T: "int"},
		{N: "Int8", T: "int8"},
		{N: "Int16", T: "int16"},
		{N: "Int32", T: "int32"},
		{N: "Int64", T: "int64"},
		{N: "Uint", T: "uint"},
		{N: "Uint8", T: "uint8"},
		{N: "Uint16", T: "uint16"},
		{N: "Uint32", T: "uint32"},
		{N: "Uint64", T: "uint64"},
		{N: "Float32", T: "float32"},
		{N: "Float64", T: "float64"},
	}
	return slicex.Map(numerics, func(nt namedType) metatype.Var {
		return metatype.Var{
			Name:  nt.N,
			Value: "NumberCheckerProvider[" + nt.T + "]{}",
			DocLines: []string{
				fmt.Sprintf("%s provides checks on type %s.", nt.N, nt.T),
			},
		}
	})
}

// File filters

func isProviderFile(file fs.FileInfo) bool {
	return strings.HasPrefix(file.Name(), "providers_") &&
		!isTestFile(file) &&
		!isBaseFile(file)
}

func isTestFile(file fs.FileInfo) bool {
	return strings.HasSuffix(file.Name(), "_test.go")
}

func isBaseFile(file fs.FileInfo) bool {
	return strings.HasSuffix(file.Name(), "_base.go")
}

// String helpers

const providerStructSuffix = "CheckerProvider"

func structToVarName(structName string) string {
	return strings.TrimSuffix(structName, providerStructSuffix)
}

func structInstanceString(name string) string {
	return name + "{}"
}
