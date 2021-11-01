package gen

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"go/types"
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
	Interfaces []metatype.Interface
	Vars       []metatype.Var
}

func computeInterfaces() (ProvidersMetaData, error) {
	docp, err := newDocPackage("check", isProviderFile) // TODO: package-agnostic
	if err != nil {
		return ProvidersMetaData{}, err
	}

	data := ProvidersMetaData{Vars: numericMetaVars()}
	for _, t := range docp.Types {
		data.Vars = append(data.Vars, computeMetaVar(t))

		// FIXME: t.Method is empty for numberCheckerProvider,
		// probably a bug due to it being type-parameterized.
		// Thus we hardcode its interface in the template file.
		if t.Name != "numberCheckerProvider" {
			data.Interfaces = append(data.Interfaces, computeMetaInterface(t))
		}
	}

	return data, nil
}

// computeMetaInterface returns a MetaInterface after the given *doc.Type.
// It reads and attaches the type's name and docs and iterates over its methods.
// If a method is inherited, it embeds the computed interface name of the parent
// rather than adding it to the interface.
func computeMetaInterface(t *doc.Type) metatype.Interface {
	name := structToInterfaceName(t.Name)
	mitf := metatype.Interface{
		Name: name,
		DocLines: docparser.ParseDocLines(t.Doc, map[string]string{
			t.Name: name,
		}),
	}

	for _, m := range t.Methods {
		// ignore private methods
		if !m.Decl.Name.IsExported() {
			continue
		}
		// If a method is inherited, embed the interface of the parent
		// rather than including all its methods
		if m.Level != 0 {
			orig := strings.TrimPrefix(m.Orig, "*") // m.Orig might have a leading "*"
			mitf.EmbedInterface(structToInterfaceName(orig))
			continue
		}
		mitf.AddFunc(metatype.Func{
			Sign:     m.Name + strings.TrimPrefix(types.ExprString(m.Decl.Type), "func"),
			DocLines: docparser.ParseDocLines(m.Doc, nil),
		})
	}
	return mitf
}

// computeMetaVar returns a MetaVar after the given *doc.Type.
func computeMetaVar(t *doc.Type) metatype.Var {
	return metatype.Var{ // var T TCheckerProvider = tCheckerProvider{}
		Name:  structToVarName(t.Name),
		Type:  structToInterfaceName(t.Name),
		Value: t.Name + "{}",
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
			Type:  "NumberCheckerProvider[" + nt.T + "]",
			Value: "numberCheckerProvider[" + nt.T + "]{}",
		}
	})
}

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
