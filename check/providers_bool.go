package check

// boolCheckerProvider provides checks on type bool.
type boolCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten bool is equal to the target.
func (p boolCheckerProvider) Is(tar bool) BoolChecker {
	pass := func(got bool) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return p.explain(label, tar, got)
	}
	return NewBoolChecker(pass, expl)
}
