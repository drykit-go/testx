package gen

import (
	"fmt"
	"go/ast"
	"go/doc"
	"log"
	"strings"
)

type serializer struct{}

func (s serializer) serializeExpr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + s.serializeExpr(t.Elt)
	case *ast.SelectorExpr:
		return s.serializeExpr(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + s.serializeExpr(t.X)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		log.Panicf("unhandled ast.Expr: %#v", expr)
		return ""
	}
}

func (s serializer) joinIdentifiers(idents []*ast.Ident, sep string) (out string) {
	for i, ident := range idents {
		out += ident.Name
		if i < len(idents)-1 {
			out += sep
		}
	}
	return
}

func (s serializer) buildSignature(f *doc.Func) string {
	b := strings.Builder{}
	b.WriteString(f.Name + "(")

	params := f.Decl.Type.Params.List
	for i, p := range params {
		pname := s.joinIdentifiers(p.Names, ", ")
		ptype := s.serializeExpr(p.Type)
		b.WriteString(pname + " " + ptype)
		if i < len(params)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString(") ")

	results := f.Decl.Type.Results
	for i, r := range results.List {
		b.WriteString(fmt.Sprint(r.Type))
		if i < results.NumFields()-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

func (s serializer) computeDocLines(src string, repl map[string]string) []string {
	lines := []string{}
	if src == "" {
		return lines
	}

	replArgs := []string{}
	for k, v := range repl {
		replArgs = append(replArgs, k, v)
	}
	replacer := strings.NewReplacer(replArgs...)

	for _, l := range strings.Split(src, "\n") {
		lines = append(lines, replacer.Replace(l))
	}

	return lines[:len(lines)-1]
}
