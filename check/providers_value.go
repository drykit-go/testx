package check

import (
	"fmt"
	"reflect"
)

// valueCheckerProvider provides checks on type interface{}.
type valueCheckerProvider struct{ baseCheckerProvider }

// Custom checks the gotten value passes the given ValuePassFunc.
// The description should give information about the expected value,
// as it outputs in format "exp <desc>" in case of failure.
func (p valueCheckerProvider) Custom(desc string, f ValuePassFunc) ValueChecker {
	expl := func(label string, got interface{}) string {
		return p.explain(label, desc, got)
	}
	return NewValueChecker(f, expl)
}

// Is checks the gotten value is equal to the target.
func (p valueCheckerProvider) Is(tar interface{}) ValueChecker {
	pass := func(got interface{}) bool { return deq(got, tar) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, tar, got)
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
		return p.explainNot(label, match, got)
	}
	return NewValueChecker(pass, expl)
}

// IsZero checks the gotten value is a zero value, indicating it might not
// have been initialized.
func (p valueCheckerProvider) IsZero() ValueChecker {
	pass := func(got interface{}) bool { return reflect.ValueOf(got).IsZero() }
	expl := func(label string, got interface{}) string {
		return p.explain(label, "to be a zero value", got)
	}
	return NewValueChecker(pass, expl)
}

// NotZero checks the gotten struct contains at least 1 non-zero value,
// meaning it has been initialized.
func (p valueCheckerProvider) NotZero() ValueChecker {
	pass := func(got interface{}) bool { return !p.IsZero().Pass(got) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, "to be a zero value", got)
	}
	return NewValueChecker(pass, expl)
}

// SameJSON checks the gotten value and the target value
// produce the same JSON, ignoring the keys order.
// It panics if any error occurs in the marshaling process.
func (p valueCheckerProvider) SameJSON(tar interface{}) ValueChecker {
	var gotDec, tarDec interface{}
	pass := func(got interface{}) bool {
		return sameJSONproduced(got, tar, &gotDec, &tarDec)
	}
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			fmt.Sprintf("json data: %v", tarDec),
			fmt.Sprintf("json data: %v", gotDec),
		)
	}
	return NewValueChecker(pass, expl)
}
