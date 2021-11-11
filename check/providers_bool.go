package check

// BoolCheckerProvider provides checks on type bool.
type BoolCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten bool is equal to the target.
func (p BoolCheckerProvider) Is(tar bool) Checker[bool] {
	pass := func(got bool) bool { return got == tar }
	expl := func(label string, got any) string {
		return p.explain(label, tar, got)
	}
	return NewChecker(pass, expl)
}
