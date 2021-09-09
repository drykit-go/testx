package astserializer

import (
	"go/ast"
	"log"
	"strings"
)

// stringBuilder is a wrapper of strings.Builder with additionnal methods
// to write serialized representation of go/ast types.
type stringBuilder struct{ strings.Builder }

// BuildFuncSignature builds a func signature given a name an *ast.FuncType
// and returns it as a string.
func BuildFuncSignature(name string, ftyp *ast.FuncType) string {
	b := stringBuilder{}
	b.writeFunc(name, ftyp)
	return b.String()
}

// writeExpr writes a serialized ast.Expr to the builder.
func (b *stringBuilder) writeExpr(expr ast.Expr) {
	switch t := expr.(type) {
	case *ast.Ident:
		b.WriteString(t.Name)
	case *ast.ArrayType:
		b.WriteString("[]")
		b.writeExpr(t.Elt)
	case *ast.Ellipsis:
		b.WriteString("...")
		b.writeExpr(t.Elt)
	case *ast.FuncType:
		b.writeFunc("func", t)
	case *ast.SelectorExpr:
		b.writeSelector(t)
	case *ast.StarExpr:
		b.WriteString("*")
		b.writeExpr(t.X)
	case *ast.ChanType:
		b.writeChan(t)
	case *ast.MapType:
		b.writeMap(t)
	case *ast.StructType:
		b.writeStruct(t)
	case *ast.InterfaceType:
		b.writeInterface(t)
	default:
		log.Panicf("❌ unhandled ast.Expr: %#v", t)
	}
}

// writeFunc writes a serialized func type to the builder.
func (b *stringBuilder) writeFunc(name string, t *ast.FuncType) {
	b.WriteString(name + "(")
	b.writeFuncParams(t.Params)
	b.WriteString(") ")
	b.writeFuncResults(t.Results)
}

// writeFuncParam writes a serialized func params list to the builder.
// It does nothing if params == nil
func (b *stringBuilder) writeFuncParams(params *ast.FieldList) {
	if params == nil {
		return
	}
	b.writeJoinedFields(params.List, ", ")
}

// writeFuncResults writes a serialized func results list to the builder.
// It does nothing if results == nil. If results contains several values,
// they are wrapped in parentheses.
func (b *stringBuilder) writeFuncResults(results *ast.FieldList) {
	if results == nil {
		return
	}
	if results.NumFields() > 1 {
		b.WriteString("(")
		defer b.WriteString(")")
	}
	b.writeJoinedFields(results.List, ", ")
}

// writeChan writes a serialized chan type to the builder.
func (b *stringBuilder) writeChan(t *ast.ChanType) {
	var chanstr string
	switch t.Dir {
	case ast.RECV:
		chanstr = "<-chan"
	case ast.SEND:
		chanstr = "chan<-"
	default:
		chanstr = "chan"
	}
	b.WriteString(chanstr)
	b.WriteString(" ")
	b.writeExpr(t.Value)
}

// writeStruct writes a serialized struct type to the builder.
func (b *stringBuilder) writeStruct(t *ast.StructType) {
	b.WriteString("struct{")
	b.writeJoinedFields(t.Fields.List, "; ")
	b.WriteString("}")
}

// writeStruct writes a serialized map type to the builder.
func (b *stringBuilder) writeMap(t *ast.MapType) {
	b.WriteString("map[")
	b.writeExpr(t.Key)
	b.WriteString("]")
	b.writeExpr(t.Value)
}

// writeStruct writes a serialized interface type to the builder.
func (b *stringBuilder) writeInterface(t *ast.InterfaceType) {
	b.WriteString("interface{")
	b.writeJoinedInterfaceFields(t.Methods.List, "; ")
	b.WriteString("}")
}

// writeStruct writes a serialized selector expression to the builder.
func (b *stringBuilder) writeSelector(t *ast.SelectorExpr) {
	b.writeExpr(t.X)
	b.WriteString(".")
	b.WriteString(t.Sel.Name)
}

// writeStruct writes joined identifiers separated by sep to the builder
func (b *stringBuilder) writeJoinedIdents(idents []*ast.Ident, sep string) {
	for i, ident := range idents {
		b.WriteString(ident.Name)
		if i != len(idents)-1 {
			b.WriteString(sep)
		}
	}
}

// writeJoinedFields writes joined fields separated by sep to the builder.
// It can be used to serialiaze a list of func params or results,
// as well as struct fields.
// For interface fields and methods, use writeJoinedInterfaceFields instead
func (b *stringBuilder) writeJoinedFields(fields []*ast.Field, sep string) {
	for i, f := range fields {
		b.writeJoinedIdents(f.Names, ", ")
		b.WriteString(" ")
		b.writeExpr(f.Type)
		if i != len(fields)-1 {
			b.WriteString(sep)
		}
	}
}

func (b *stringBuilder) writeJoinedInterfaceFields(fields []*ast.Field, sep string) {
	for i, f := range fields {
		switch t := f.Type.(type) {
		case *ast.Ident:
			b.WriteString(t.Name)
		case *ast.FuncType:
			b.WriteString(BuildFuncSignature(f.Names[0].Name, t))
		default:
			log.Panicf("joinInterfaceFields: unhandled type %#v", t)
		}
		if i != len(fields) {
			b.WriteString(sep)
		}
	}
}
