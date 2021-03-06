package check

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/drykit-go/testx/internal/reflectutil"
)

// structCheckerProvider provides checks on kind struct.
type structCheckerProvider struct{ valueCheckerProvider }

// FieldsEqual checks all given fields equal the exp value.
// It panics if the fields do not exist or are not exported,
// or if the tested value is not a struct.
func (p structCheckerProvider) FieldsEqual(exp interface{}, fields []string) ValueChecker {
	var bads []string
	pass := func(got interface{}) bool {
		reflectutil.MustBeOfKind(got, reflect.Struct)
		bads = p.badFields(got, fields, func(k string, v interface{}) bool {
			return p.deq(v, exp)
		})
		return len(bads) == 0
	}
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			fmt.Sprintf("fields [%s] to equal %v", p.formatFields(fields), exp),
			strings.Join(bads, ", "),
		)
	}
	return NewValueChecker(pass, expl)
}

// CheckFields checks all given fields pass the ValueChecker.
// It panics if the fields do not exist or are not exported,
// or if the tested value is not a struct.
func (p structCheckerProvider) CheckFields(c ValueChecker, fields []string) ValueChecker {
	var bads []string
	pass := func(got interface{}) bool {
		reflectutil.MustBeOfKind(got, reflect.Struct)
		bads = p.badFields(got, fields, func(k string, v interface{}) bool {
			return c.Pass(v)
		})
		return len(bads) == 0
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			fmt.Sprintf("fields [%s] to pass ValueChecker", p.formatFields(fields)),
			c.Explain("fields", strings.Join(bads, ", ")),
		)
	}
	return NewValueChecker(pass, expl)
}

func (p structCheckerProvider) badFields(
	gotstruct interface{},
	fields []string,
	pass func(k string, v interface{}) bool,
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

func (structCheckerProvider) formatFields(fields []string) string {
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
