package check

import (
	"fmt"
	"reflect"
	"strings"
)

// structCheckerProvider provides checks on kind struct.
type structCheckerProvider struct{ baseCheckerProvider }

// SameJSON checks the gotten struct and the target value
// produce the same JSON, ignoring the keys order.
// It panics if any error occurs in the marshaling process.
func (structCheckerProvider) SameJSON(tar interface{}) ValueChecker {
	return Value.SameJSON(tar)
}

// IsZero checks the gotten struct only contains zero values,
// meaning it has not been initialized.
func (structCheckerProvider) IsZero() ValueChecker {
	return Value.IsZero()
}

// NotZero checks the gotten struct contains at least 1 non-zero value,
// meaning it has been initialized.
func (p structCheckerProvider) NotZero() ValueChecker {
	return Value.NotZero()
}

// FieldsEqual checks all given fields equal the exp value.
// It panics if the fields do not exist or are not exported,
// or if the tested value is not a struct.
func (p structCheckerProvider) FieldsEqual(exp interface{}, fields []string) ValueChecker {
	var badFields []string
	pass := func(got interface{}) bool {
		gotv := reflect.ValueOf(got)
		for _, f := range fields {
			// panic hasard: field must exist and be exported
			gotf := gotv.FieldByName(f).Interface()
			if !deq(gotf, exp) {
				badFields = append(badFields, fmt.Sprintf(".%s=%v", f, gotf))
			}
		}
		return len(badFields) == 0
	}
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			fmt.Sprintf("fields to equal %v", exp),
			strings.Join(badFields, ","),
		)
	}
	return NewValueChecker(pass, expl)
}

// FieldsEqual checks all given fields pass the ValueChecker.
// It panics if the fields do not exist or are not exported,
// or if the tested value is not a struct.
func (p structCheckerProvider) CheckFields(c ValueChecker, fields []string) ValueChecker {
	var badFields []string
	pass := func(got interface{}) bool {
		gotv := reflect.ValueOf(got)
		for _, f := range fields {
			// panic hasard: field must exist and be exported
			gotf := gotv.FieldByName(f).Interface()
			if !c.Pass(gotf) {
				badFields = append(badFields, fmt.Sprintf(".%s", f))
			}
		}
		return len(badFields) == 0
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			fmt.Sprintf("fields [%s] to pass ValueChecker", strings.Join(badFields, ",")),
			c.Explain("field", got),
		)
	}
	return NewValueChecker(pass, expl)
}
