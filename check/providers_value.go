package check

import (
	"fmt"

	"github.com/drykit-go/testx/internal/reflectutil"
)

// valueCheckerProvider provides checks on type interface{}.
type valueCheckerProvider struct{ baseCheckerProvider }

// Custom checks the gotten value passes the given PassFunc[interface{}].
// The description should give information about the expected value,
// as it outputs in format "exp <desc>" in case of failure.
func (p valueCheckerProvider) Custom(desc string, f PassFunc[interface{}]) Checker[interface{}] {
	expl := func(label string, got interface{}) string {
		return p.explain(label, desc, got)
	}
	return NewChecker(f, expl)
}

// Is checks the gotten value is equal to the target.
func (p valueCheckerProvider) Is(tar interface{}) Checker[interface{}] {
	pass := func(got interface{}) bool { return p.deq(got, tar) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, tar, got)
	}
	return NewChecker(pass, expl)
}

// Not checks the gotten value is not equal to the target.
func (p valueCheckerProvider) Not(values ...interface{}) Checker[interface{}] {
	var match interface{}
	pass := func(got interface{}) bool {
		for _, v := range values {
			if p.deq(got, v) {
				match = v
				return false
			}
		}
		return true
	}
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, match, got)
	}
	return NewChecker(pass, expl)
}

// IsZero checks the gotten value is a zero value, indicating it might not
// have been initialized.
func (p valueCheckerProvider) IsZero() Checker[interface{}] {
	expl := func(label string, got interface{}) string {
		return p.explain(label, "to be a zero value", got)
	}
	return NewChecker(reflectutil.IsZero, expl)
}

// NotZero checks the gotten struct contains at least 1 non-zero value,
// meaning it has been initialized.
func (p valueCheckerProvider) NotZero() Checker[interface{}] {
	pass := func(got interface{}) bool { return !reflectutil.IsZero(got) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, "to be a zero value", got)
	}
	return NewChecker(pass, expl)
}

// SameJSON checks the gotten value and the target value
// produce the same JSON, ignoring formatting and keys order.
// It panics if any error occurs in the marshaling process.
func (p valueCheckerProvider) SameJSON(tar interface{}) Checker[interface{}] {
	var gotDec, tarDec interface{}
	pass := func(got interface{}) bool {
		return p.sameJSONProduced(got, tar, &gotDec, &tarDec)
	}
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			fmt.Sprintf("json data: %v", tarDec),
			fmt.Sprintf("json data: %v", gotDec),
		)
	}
	return NewChecker(pass, expl)
}
