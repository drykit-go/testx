package check

import (
	"fmt"
)

// numberCheckerProvider provides checks on numeric types.
type numberCheckerProvider[T Numeric] struct{ baseCheckerProvider }

// Is checks the gotten Number is equal to the target.
func (p numberCheckerProvider[T]) Is(tar T) Checker[T] {
	pass := func(got T) bool { return got == tar }
	expl := func(label string, got any) string {
		return p.explain(label, tar, got)
	}
	return NewChecker(pass, expl)
}

// Not checks the gotten Number is not equal to the target.
func (p numberCheckerProvider[T]) Not(values ...T) Checker[T] {
	var match T
	pass := func(got T) bool {
		for _, v := range values {
			if got == v {
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

// InRange checks the gotten Number is in the closed interval [lo:hi].
func (p numberCheckerProvider[T]) InRange(lo, hi T) Checker[T] {
	pass := func(got T) bool { return p.inrange(got, lo, hi) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf("in range [%v:%v]", lo, hi), got)
	}
	return NewChecker(pass, expl)
}

// OutRange checks the gotten Number is not in the closed interval [lo:hi].
func (p numberCheckerProvider[T]) OutRange(lo, hi T) Checker[T] {
	pass := func(got T) bool { return !p.inrange(got, lo, hi) }
	expl := func(label string, got any) string {
		return p.explainNot(label, fmt.Sprintf("in range [%v:%v]", lo, hi), got)
	}
	return NewChecker(pass, expl)
}

// GT checks the gotten Number is greater than the target.
func (p numberCheckerProvider[T]) GT(tar T) Checker[T] {
	pass := func(got T) bool { return !p.lte(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf("> %v", tar), got)
	}
	return NewChecker(pass, expl)
}

// GTE checks the gotten Number is greater or equal to the target.
func (p numberCheckerProvider[T]) GTE(tar T) Checker[T] {
	pass := func(got T) bool { return !p.lt(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf(">= %v", tar), got)
	}
	return NewChecker(pass, expl)
}

// LT checks the gotten Number is lesser than the target.
func (p numberCheckerProvider[T]) LT(tar T) Checker[T] {
	pass := func(got T) bool { return p.lt(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf("< %v", tar), got)
	}
	return NewChecker(pass, expl)
}

// LTE checks the gotten Number is lesser or equal to the target.
func (p numberCheckerProvider[T]) LTE(tar T) Checker[T] {
	pass := func(got T) bool { return p.lte(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf("<= %v", tar), got)
	}
	return NewChecker(pass, expl)
}

// Helpers

func (numberCheckerProvider[T]) lt(a, b T) bool  { return a < b }
func (numberCheckerProvider[T]) lte(a, b T) bool { return a <= b }
func (p numberCheckerProvider[T]) inrange(n, lo, hi T) bool {
	return !p.lt(n, lo) && p.lte(n, hi)
}
