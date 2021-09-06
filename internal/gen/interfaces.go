package gen

import (
	"errors"
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
	Docs    []string
	Name    string
	Methods []MetaMethod
}

type MetaMethod struct {
	Sign string
	Docs []string
}

type MetaVar struct {
	Name, Type, Value string
}

func computeInterfaces() (ProvidersTemplateData, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "./", isProviderFile, parser.ParseComments)
	if err != nil {
		return ProvidersTemplateData{}, err
	}
	astp, ok := pkgs["check"] // TODO: package-agnostic
	if !ok {
		return ProvidersTemplateData{}, errors.New(
			"no files found for package \"check\", check path or filters",
		)
	}
	docp := doc.New(astp, "./", doc.AllDecls)

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
