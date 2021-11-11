package check

import (
	"fmt"

	"github.com/drykit-go/testx/internal/reflectutil"
)

// ValueCheckerProvider provides checks on type any.
type ValueCheckerProvider struct{ baseCheckerProvider }

// Custom checks the gotten value passes the given PassFunc[any].
// The description should give information about the expected value,
// as it outputs in format "exp <desc>" in case of failure.
func (p ValueCheckerProvider) Custom(desc string, f PassFunc[any]) Checker[any] {
	expl := func(label string, got any) string {
		return p.explain(label, desc, got)
	}
	return NewChecker(f, expl)
}

// Is checks the gotten value is equal to the target.
func (p ValueCheckerProvider) Is(tar any) Checker[any] {
	pass := func(got any) bool { return p.deq(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, tar, got)
	}
	return NewChecker(pass, expl)
}

// Not checks the gotten value is not equal to the target.
func (p ValueCheckerProvider) Not(values ...any) Checker[any] {
	var match any
	pass := func(got any) bool {
		for _, v := range values {
			if p.deq(got, v) {
				match = v
				return false
			}
		}
		return true
	}
	expl := func(label string, got any) string {
		return p.explainNot(label, match, got)
	}
	return NewChecker(pass, expl)
}

// IsZero checks the gotten value is a zero value, indicating it might not
// have been initialized.
func (p ValueCheckerProvider) IsZero() Checker[any] {
	expl := func(label string, got any) string {
		return p.explain(label, "to be a zero value", got)
	}
	return NewChecker(reflectutil.IsZero, expl)
}

// NotZero checks the gotten struct contains at least 1 non-zero value,
// meaning it has been initialized.
func (p ValueCheckerProvider) NotZero() Checker[any] {
	pass := func(got any) bool { return !reflectutil.IsZero(got) }
	expl := func(label string, got any) string {
		return p.explainNot(label, "to be a zero value", got)
	}
	return NewChecker(pass, expl)
}

// SameJSON checks the gotten value and the target value
// produce the same JSON, ignoring formatting and keys order.
// It panics if any error occurs in the marshaling process.
func (p ValueCheckerProvider) SameJSON(tar any) Checker[any] {
	var gotDec, tarDec any
	pass := func(got any) bool {
		return p.sameJSONProduced(got, tar, &gotDec, &tarDec)
	}
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("json data: %v", tarDec),
			fmt.Sprintf("json data: %v", gotDec),
		)
	}
	return NewChecker(pass, expl)
}
