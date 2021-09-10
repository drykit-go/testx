package gen

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"strings"

	"github.com/drykit-go/testx/internal/gen/metatype"
	"github.com/drykit-go/testx/internal/gen/serialize"
)

const interfaceSuffix = "CheckerProvider"

// TODO: use strcase.Pascal, remove record
var caseMapping = map[string]string{
	"bool":       "Bool",
	"int":        "Int",
	"bytes":      "Bytes",
	"string":     "String",
	"duration":   "Duration",
	"httpHeader": "HTTPHeader",
	"value":      "Value",
	"struct":     "Struct",
	"map":        "Map",
	"slice":      "Slice",
}

func typeName(structName string) string {
	return strings.TrimSuffix(structName, interfaceSuffix)
}

func interfaceName(structName string) string {
	upper, ok := caseMapping[typeName(structName)]
	if !ok {
		log.Panicf("missing gen.CaseMapping type key for %s", structName)
	}
	return upper + interfaceSuffix
}

type ProvidersMetaData struct {
	Interfaces []metatype.Interface
	Vars       []metatype.Var
}

func computeInterfaces() (ProvidersMetaData, error) {
	docp, err := newDocPackage("check", isProviderFile) // TODO: package-agnostic
	if err != nil {
		return ProvidersMetaData{}, err
	}

	data := ProvidersMetaData{}
	for _, t := range docp.Types {
		data.Vars = append(data.Vars, computeMetaVar(t))
		data.Interfaces = append(data.Interfaces, computeMetaInterface(t))
	}

	return data, nil
}

// computeMetaInterface returns a MetaInterface after the given *doc.Type.
// It reads and attaches the type's name and docs and iterates over its methods.
// If a method is inherited, it embeds the computed interface name of the parent
// rather than adding it to the interface.
func computeMetaInterface(t *doc.Type) metatype.Interface {
	name := interfaceName(t.Name)
	mitf := metatype.Interface{
		Name: name,
		Docs: serialize.DocLines(t.Doc, map[string]string{
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
			mitf.EmbedInterface(interfaceName(orig))
			continue
		}
		mitf.AddFunc(metatype.Func{
			Sign: serialize.FuncSignature(m.Name, m.Decl.Type),
			Docs: serialize.DocLines(m.Doc, nil),
		})
	}
	return mitf
}

// computeMetaVar returns a MetaVar after the given *doc.Type.
func computeMetaVar(t *doc.Type) metatype.Var {
	return metatype.Var{ // var T TCheckerProvider = tCheckerProvider{}
		Name:  caseMapping[typeName(t.Name)],
		Type:  interfaceName(t.Name),
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
	astp, ok := pkgs[packageName] // TODO: package-agnostic
	if !ok {
		return nil, fmt.Errorf(
			"no files found for package %s, check path or filters",
			packageName,
		)
	}
	return doc.New(astp, "./", doc.AllDecls), nil
}

func isProviderFile(file fs.FileInfo) bool {
	return strings.HasPrefix(file.Name(), "providers_") &&
		isTestFile(file) &&
		isBaseFile(file)
}

func isTestFile(file fs.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_test.go")
}

func isBaseFile(file fs.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_base.go")
}
