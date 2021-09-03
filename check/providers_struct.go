package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// structCheckerProvider provides checks on type interface{}.
type structCheckerProvider struct{}

// SameJSON checks the gotten struct and the target value
// result in the same JSON.
// It panics if any error occurs in the marshaling process.
func (structCheckerProvider) SameJSON(tar interface{}) ValueChecker {
	var gotjson, tarjson []byte
	pass := func(got interface{}) bool {
		gotjson, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			log.Panic("cannot convert struct to json:", err)
		}
		tarjson, err := json.MarshalIndent(tar, "", "  ")
		if err != nil {
			log.Panic("cannot convert target to json:", err)
		}
		return bytes.Equal(gotjson, tarjson)
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"exp %s to match JSON:\n"+
				"%s\n"+
				"got:\n"+
				"%s",
			label, string(tarjson), string(gotjson),
		)
	}
	return NewValueChecker(pass, expl)
}

// IsZero checks the gotten struct only contains zero values,
// meaning it has not been initialized.
func (structCheckerProvider) IsZero() ValueChecker {
	pass := func(got interface{}) bool {
		gotv := reflect.ValueOf(got)
		return gotv.IsZero()
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"exp %s to contain only zero values\n"+
				"got %#v",
			label, got,
		)
	}
	return NewValueChecker(pass, expl)
}

// NotZero checks the gotten struct contains at least 1 non-zero value,
// meaning it has been initialized.
func (p structCheckerProvider) NotZero() ValueChecker {
	pass := func(got interface{}) bool {
		return !p.IsZero().Pass(got)
	}
	expl := func(label string, got interface{}) string {
		return p.IsZero().Explain(label+" not", got)
	}
	return NewValueChecker(pass, expl)
}

// FieldsEqual checks all given fields equal the exp value.
// It panics if the fields do not exist or are not exported,
// or if the tested value is not a struct.
func (p structCheckerProvider) FieldsEqual(exp interface{}, fields ...string) ValueChecker {
	var badFields []string
	pass := func(got interface{}) bool {
		gotv := reflect.ValueOf(got)
		for _, f := range fields {
			// panic hasard: field must exist and be exported
			gotf := gotv.FieldByName(f).Interface()
			if !p.eq(gotf, exp) {
				badFields = append(badFields, fmt.Sprintf(".%s=%v", f, gotf))
			}
		}
		return len(badFields) == 0
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"exp %s fields to equal %v\n"+
				"got %s",
			label, exp,
			strings.Join(badFields, ", "),
		)
	}
	return NewValueChecker(pass, expl)
}

// FieldsEqual checks all given fields pass the ValueChecker.
// It panics if the fields do not exist or are not exported,
// or if the tested value is not a struct.
func (structCheckerProvider) CheckFields(c ValueChecker, fields ...string) ValueChecker {
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
		return fmt.Sprintf(
			"exp %s fields to pass ValueChecker\n"+
				"got %s -> %s",
			label,
			strings.Join(badFields, ", "), c.Explain(label, got),
		)
	}
	return NewValueChecker(pass, expl)
}

func (p structCheckerProvider) eq(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
