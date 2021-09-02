package check

import "fmt"

// intCheckerProvider provides checks on type int.
type intCheckerProvider struct{}

// Is checks the gotten int is equal to the target.
func (intCheckerProvider) Is(tar int) IntChecker {
	pass := func(got int) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s == %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// Not checks the gotten int is not equal to the target.
func (intCheckerProvider) Not(tar int) IntChecker {
	pass := func(got int) bool { return got != tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s != %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// InRange checks the gotten int is in the closed interval [lo:hi].
func (f intCheckerProvider) InRange(lo, hi int) IntChecker {
	pass := func(got int) bool { return f.inrange(got, lo, hi) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s in range [%d:%d], got %d",
			label, lo, hi, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// OutRange checks the gotten int is not in the closed interval [lo:hi].
func (f intCheckerProvider) OutRange(lo, hi int) IntChecker {
	pass := func(got int) bool { return !f.inrange(got, lo, hi) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not in range [%d:%d], got %d",
			label, lo, hi, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// GT checks the gotten int is greater than the target.
func (f intCheckerProvider) GT(tar int) IntChecker {
	pass := func(got int) bool { return !f.lte(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s > %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// GTE checks the gotten int is greater or equal to the target.
func (f intCheckerProvider) GTE(tar int) IntChecker {
	pass := func(got int) bool { return !f.lt(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s >= to %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// LT checks the gotten int is lesser than the target.
func (f intCheckerProvider) LT(tar int) IntChecker {
	pass := func(got int) bool { return f.lt(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s < %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// LTE checks the gotten int is lesser or equal to the target.
func (f intCheckerProvider) LTE(tar int) IntChecker {
	pass := func(got int) bool { return f.lte(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s <= to %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// Helpers

func (f intCheckerProvider) lt(a, b int) bool  { return a < b }
func (f intCheckerProvider) lte(a, b int) bool { return a <= b }
func (f intCheckerProvider) inrange(n, lo, hi int) bool {
	return !f.lt(n, lo) && f.lte(n, hi)
}
