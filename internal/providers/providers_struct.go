package providers

import (
	"fmt"
	"reflect"
	"strings"

	check "github.com/drykit-go/testx/internal/checktypes"
	"github.com/drykit-go/testx/internal/reflectutil"
)

// StructCheckerProvider provides checks on kind struct.
type StructCheckerProvider struct{ ValueCheckerProvider[any] }

// FieldsEqual checks all given fields equal the exp value.
// It panics if the fields do not exist or are not exported,
// or if the tested value is not a struct.
func (p StructCheckerProvider) FieldsEqual(exp any, fields []string) check.Checker[any] {
	var bads []string
	pass := func(got any) bool {
		reflectutil.MustBeOfKind(got, reflect.Struct)
		bads = p.badFields(got, fields, func(k string, v any) bool {
			return p.deq(v, exp)
		})
		return len(bads) == 0
	}
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("fields [%s] to equal %v", p.formatFields(fields), exp),
			strings.Join(bads, ", "),
		)
	}
	return check.NewChecker(pass, expl)
}

// CheckFields checks all given fields pass the Checker[any].
// It panics if the fields do not exist or are not exported,
// or if the tested value is not a struct.
func (p StructCheckerProvider) CheckFields(c check.Checker[any], fields []string) check.Checker[any] {
	var bads []string
	pass := func(got any) bool {
		reflectutil.MustBeOfKind(got, reflect.Struct)
		bads = p.badFields(got, fields, func(k string, v any) bool {
			return c.Pass(v)
		})
		return len(bads) == 0
	}
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			fmt.Sprintf("fields [%s] to pass Checker[any]", p.formatFields(fields)),
			c.Explain("fields", strings.Join(bads, ", ")),
		)
	}
	return check.NewChecker(pass, expl)
}

func (p StructCheckerProvider) badFields(
	gotstruct any,
	fields []string,
	pass func(k string, v any) bool,
) (bads []string) {
	vstruct := reflect.ValueOf(gotstruct)
	for _, k := range fields {
		// panic hazard: field must exist and be exported
		v := vstruct.FieldByName(k).Interface()
		if !pass(k, v) {
			bads = append(bads, fmt.Sprintf(".%s=%v", k, v))
		}
	}
	return
}

func (StructCheckerProvider) formatFields(fields []string) string {
	n := len(fields)
	var b strings.Builder
	for i, f := range fields {
		b.WriteByte('.')
		b.WriteString(f)
		if i != n-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}
