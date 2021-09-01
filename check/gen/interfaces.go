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
	"int":        "Int",
	"bytes":      "Bytes",
	"string":     "String",
	"duration":   "Duration",
	"httpHeader": "HTTPHeader",
	"value":      "Value",
}

func interfaceName(factoryName string) string {
	typ := strings.TrimSuffix(factoryName, interfaceSuffix)
	upper, ok := caseMapping[typ]
	if !ok {
		log.Panicf("missing gen.CaseMapping type key for %s", factoryName)
	}
	return upper + interfaceSuffix
}

type InterfaceData struct {
	Docs    []string
	Name    string
	Methods []MethodData
}

type MethodData struct {
	Sign string
	Docs []string
}

func computeInterfaces() ([]InterfaceData, error) {
	seri := serializer{}
	fset := token.NewFileSet()

	dir, err := parser.ParseDir(fset, "./", implFilesOnly, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	astp := dir["check"]
	docp := doc.New(astp, "./", doc.AllDecls)

	interfaces := []InterfaceData{}

	for _, t := range docp.Types {
		itf := InterfaceData{
			Name: interfaceName(t.Name),
			Docs: seri.computeDocLines(t.Doc, map[string]string{
				t.Name: interfaceName(t.Name),
			}),
		}

		for _, m := range t.Methods {
			// ignore private methods
			if !m.Decl.Name.IsExported() {
				continue
			}
			itf.Methods = append(itf.Methods, MethodData{
				Sign: seri.buildSignature(m),
				Docs: seri.computeDocLines(m.Doc, nil),
			})
		}
		interfaces = append(interfaces, itf)
	}

	return interfaces, nil
}

func implFilesOnly(file fs.FileInfo) bool {
	return strings.HasPrefix(file.Name(), "impl_") && excludeTestFiles(file)
}

func excludeTestFiles(file fs.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_test.go")
}
