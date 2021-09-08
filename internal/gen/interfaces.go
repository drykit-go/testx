package gen

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"strings"

	"github.com/drykit-go/testx/internal/gen/astserializer"
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

type ProvidersTemplateData struct {
	Interfaces []MetaInterface
	Vars       []MetaVar
}

type MetaInterface struct {
	Docs     []string
	Name     string
	Embedded []string
	Methods  []MetaMethod
}

func (mi *MetaInterface) embedInterface(name string) {
	for _, itf := range mi.Embedded {
		if itf == name {
			return
		}
	}
	mi.Embedded = append(mi.Embedded, name)
}

type MetaMethod struct {
	Sign string
	Docs []string
}

type MetaVar struct {
	Name, Type, Value string
}

func computeInterfaces() (ProvidersTemplateData, error) {
	docp, err := newDocPackage("check", isProviderFile) // TODO: package-agnostic
	if err != nil {
		return ProvidersTemplateData{}, err
	}

	data := ProvidersTemplateData{}
	for _, t := range docp.Types {
		mitf := MetaInterface{ // type TCheckerProvider interface{ ... }
			Name: interfaceName(t.Name),
			Docs: astserializer.ComputeDocLines(t.Doc, map[string]string{
				t.Name: interfaceName(t.Name),
			}),
		}
		mvar := MetaVar{ // var T TCheckerProvider = tCheckerProvider{}
			Name:  caseMapping[typeName(t.Name)],
			Type:  mitf.Name,
			Value: t.Name + "{}",
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
				mitf.embedInterface(interfaceName(orig))
				continue
			}
			mitf.Methods = append(mitf.Methods, MetaMethod{
				Sign: astserializer.BuildFuncSignature(m.Name, m.Decl.Type),
				Docs: astserializer.ComputeDocLines(m.Doc, nil),
			})
		}
		data.Interfaces = append(data.Interfaces, mitf)
		data.Vars = append(data.Vars, mvar)
	}

	return data, nil
}

// newDocPackage returns a *doc.Package matching packageName after applying
// the given filter, or the first non-nil error occuring in the process.
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
