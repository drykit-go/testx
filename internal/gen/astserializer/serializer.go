package astserializer

import (
	"fmt"
	"go/ast"
	"io"
	"log"
	"strings"
)

func ComputeDocLines(src string, repl map[string]string) []string {
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

func BuildFuncSignature(name string, ftyp *ast.FuncType) string {
	b := strings.Builder{}
	b.WriteString(name + "(")
	writeFuncParamsString(&b, ftyp.Params)
	b.WriteString(") ")
	writeFuncResultsString(&b, ftyp.Results)
	return b.String()
}

func writeFuncParamsString(w io.StringWriter, params *ast.FieldList) {
	for i, p := range params.List {
		pname := joinIdentifiers(p.Names, ", ")
		ptype := serializeExpr(p.Type)
		_, err := w.WriteString(pname + " " + ptype)
		panicOnErr(err)
		if i < len(params.List)-1 {
			_, err := w.WriteString(", ")
			if err != nil {
				panicOnErr(err)
			}
		}
	}
}

func writeFuncResultsString(w io.StringWriter, results *ast.FieldList) {
	for i, r := range results.List {
		_, err := w.WriteString(fmt.Sprint(r.Type))
		panicOnErr(err)
		if i < results.NumFields()-1 {
			_, err := w.WriteString(", ")
			panicOnErr(err)
		}
	}
}

func serializeExpr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + serializeExpr(t.Elt)
	case *ast.Ellipsis:
		return "..." + serializeExpr(t.Elt)
	case *ast.FuncType:
		return BuildFuncSignature("func", t)
	case *ast.SelectorExpr:
		return serializeExpr(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + serializeExpr(t.X)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		log.Panicf("❌ unhandled ast.Expr: %#v", expr)
		return ""
	}
}

func joinIdentifiers(idents []*ast.Ident, sep string) (out string) {
	for i, ident := range idents {
		out += ident.Name
		if i < len(idents)-1 {
			out += sep
		}
	}
	return
}

func panicOnErr(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
