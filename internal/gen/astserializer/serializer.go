package astserializer

import (
	"go/ast"
	"io"
	"log"
	"strings"
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
	if params == nil {
		return
	}
	w.WriteString(joinFields(params.List, ", "))
}

func writeFuncResultsString(w io.StringWriter, results *ast.FieldList) {
	if results == nil {
		return
	}
	if results.NumFields() > 1 {
		w.WriteString("(")
		defer w.WriteString(")")
	}
	w.WriteString(joinFields(results.List, ", "))
}

func joinIdents(idents []*ast.Ident, sep string) string {
	b := strings.Builder{}
	for i, ident := range idents {
		b.WriteString(ident.Name)
		if i != len(idents)-1 {
			b.WriteString(sep)
		}
	}
	return b.String()
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
	case *ast.ChanType:
		return serializeChanType(t)
	case *ast.MapType:
		return "map[" + serializeExpr(t.Key) + "]" + serializeExpr(t.Value)
	case *ast.StructType:
		return serializeStructType(t)
	case *ast.InterfaceType:
		return serializeInterfaceType(t)
	default:
		log.Panicf("❌ unhandled ast.Expr: %#v", t)
		return ""
	}
}

func serializeChanType(t *ast.ChanType) string {
	var chanstr string
	switch t.Dir {
	case ast.RECV:
		chanstr = "<-chan"
	case ast.SEND:
		chanstr = "chan<-"
	default:
		chanstr = "chan"
	}
	return chanstr + " " + serializeExpr(t.Value)
}

func serializeStructType(t *ast.StructType) string {
	b := strings.Builder{}
	b.WriteString("struct{")
	b.WriteString(joinFields(t.Fields.List, "; "))
	b.WriteString("}")
	return b.String()
}

func serializeInterfaceType(t *ast.InterfaceType) string {
	b := strings.Builder{}
	b.WriteString("interface{")
	b.WriteString(joinInterfaceFields(t.Methods.List, "; "))
	b.WriteString("}")
	return b.String()
}

// joinFields returns a string representation of a slice of *ast.Field,
// separated by sep. It can be used to serialiaze func params or results,
// as well as struct fields.
func joinFields(fields []*ast.Field, sep string) string {
	b := strings.Builder{}
	for i, f := range fields {
		fname := joinIdents(f.Names, ", ")
		ftype := serializeExpr(f.Type)
		b.WriteString(fname)
		b.WriteString(" ")
		b.WriteString(ftype)
		if i != len(fields)-1 {
			b.WriteString(sep)
		}
	}
	return b.String()
}

func joinInterfaceFields(fields []*ast.Field, sep string) string {
	b := strings.Builder{}
	for i, field := range fields {
		switch t := field.Type.(type) {
		case *ast.Ident:
			b.WriteString(t.Name)
		case *ast.FuncType:
			b.WriteString(BuildFuncSignature(field.Names[0].Name, t))
		default:
			log.Panicf("joinInterfaceFields: unhandled type %#v", t)
		}
		if i != len(fields) {
			b.WriteString(sep)
		}
	}
	return b.String()
}
