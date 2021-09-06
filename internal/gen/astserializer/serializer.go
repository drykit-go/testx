package astserializer

import (
	"fmt"
	"go/ast"
	"io"
	"log"
	"strings"

	"github.com/drykit-go/cond"
)

// ComputeDocLines parses the raw doc (as returned by go/doc.Type.Doc),
// applies the given replacements and returns a slice of strings representing
// the resulting lines.
func ComputeDocLines(rawdoc string, repl map[string]string) []string {
	lines := []string{}
	if rawdoc == "" {
		return lines
	}

	// Set up strings replacer
	replArgs := []string{}
	for k, v := range repl {
		replArgs = append(replArgs, k, v)
	}
	replacer := strings.NewReplacer(replArgs...)

	// Parse and add new lines
	for _, l := range strings.Split(rawdoc, "\n") {
		lines = append(lines, replacer.Replace(l))
	}

	// Strip last line, always empty due to final '\n'
	return lines[:len(lines)-1]
}

// BuildFuncSignature builds a func signature given a name an *ast.FuncType
// and returns it as a string.
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
		cond.PanicOnErr(err)
		if i < len(params.List)-1 {
			_, err := w.WriteString(", ")
			cond.PanicOnErr(err)
		}
	}
}

func writeFuncResultsString(w io.StringWriter, results *ast.FieldList) {
	for i, r := range results.List {
		_, err := w.WriteString(fmt.Sprint(r.Type))
		cond.PanicOnErr(err)
		if i < results.NumFields()-1 {
			_, err := w.WriteString(", ")
			cond.PanicOnErr(err)
		}
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
