package gen

import (
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"strings"
)

const interfaceSuffix = "CheckerProvider"

var caseMapping = map[string]string{
	"bool":       "Bool",
	"int":        "Int",
	"bytes":      "Bytes",
	"string":     "String",
	"duration":   "Duration",
	"httpHeader": "HTTPHeader",
	"value":      "Value",
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
	seri := serializer{}
	fset := token.NewFileSet()

	dir, err := parser.ParseDir(fset, "./", providersFilesOnly, parser.ParseComments)
	if err != nil {
		return ProvidersTemplateData{}, err
	}
	astp := dir["check"]
	docp := doc.New(astp, "./", doc.AllDecls)

	data := ProvidersTemplateData{}
	for _, t := range docp.Types {
		itf := MetaInterface{
			Name: interfaceName(t.Name),
			Docs: seri.computeDocLines(t.Doc, map[string]string{
				t.Name: interfaceName(t.Name),
			}),
		}

		vvar := MetaVar{
			Name:  caseMapping[typeName(t.Name)],
			Type:  itf.Name,
			Value: t.Name + "{}",
		}

		for _, m := range t.Methods {
			// ignore private methods
			if !m.Decl.Name.IsExported() {
				continue
			}
			itf.Methods = append(itf.Methods, MetaMethod{
				Sign: seri.buildSignature(m),
				Docs: seri.computeDocLines(m.Doc, nil),
			})
		}
		data.Interfaces = append(data.Interfaces, itf)
		data.Vars = append(data.Vars, vvar)
	}

	return data, nil
}

func providersFilesOnly(file fs.FileInfo) bool {
	return strings.HasPrefix(file.Name(), "providers_") && excludeTestFiles(file)
}

func excludeTestFiles(file fs.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_test.go")
}
