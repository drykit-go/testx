package check

import (
	"fmt"
	"reflect"
)

// valueCheckerProvider provides checks on type interface{}.
type valueCheckerProvider struct{}

// Custom checks the gotten value passes the given ValuePassFunc.
// The description should typically begin with keywords like "expect"
// or "should" for intelligible output.
// For instance, "expect odd number" would output:
// 	> "expect odd number, got 42"
func (valueCheckerProvider) Custom(desc string, f ValuePassFunc) ValueChecker {
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"%s: %s, got %v",
			label, desc, got,
		)
	}
	return NewValueChecker(f, expl)
}

// Is checks the gotten value is equal to the target.
func (p valueCheckerProvider) Is(tar interface{}) ValueChecker {
	pass := func(got interface{}) bool { return deq(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to equal %v, got %v",
			label, tar, got,
		)
	}
	return NewValueChecker(pass, expl)
}

// Not checks the gotten value is not equal to the target.
func (p valueCheckerProvider) Not(values ...interface{}) ValueChecker {
	var match interface{}
	pass := func(got interface{}) bool {
		for _, v := range values {
			if deq(got, v) {
				match = v
				return false
			}
		}
		return true
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to equal %v, got %v",
			label, match, got,
		)
	}
	return NewValueChecker(pass, expl)
}

// IsZero checks the gotten value is or only contains zero values,
// meaning it has not been initialized.
func (valueCheckerProvider) IsZero() ValueChecker {
	pass := func(got interface{}) bool {
		return reflect.ValueOf(got).IsZero()
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"exp %s to be or contain only zero values\n"+
				"got %#v",
			label, got,
		)
	}
	return NewValueChecker(pass, expl)
}

// NotZero checks the gotten struct contains at least 1 non-zero value,
// meaning it has been initialized.
func (p valueCheckerProvider) NotZero() ValueChecker {
	pass := func(got interface{}) bool {
		return !p.IsZero().Pass(got)
	}
	expl := func(label string, got interface{}) string {
		return p.IsZero().Explain(label+" not", got)
	}
	return NewValueChecker(pass, expl)
}

// SameJSON checks the gotten value and the target value
// produce the same JSON, ignoring the keys order.
// It panics if any error occurs in the marshaling process.
func (valueCheckerProvider) SameJSON(tar interface{}) ValueChecker {
	var gotDec, tarDec interface{}
	pass := func(got interface{}) bool {
		return sameJSONproduced(got, tar, &gotDec, &tarDec)
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"exp %s to match JSON:\n"+
				"%#v\n"+
				"got:\n"+
				"%#v",
			label, tarDec, gotDec,
		)
	}
	return NewValueChecker(pass, expl)
}
