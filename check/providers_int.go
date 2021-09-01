package check

import "fmt"

type intCheckerProvider struct{}

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
