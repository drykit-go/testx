package providers

import check "github.com/drykit-go/testx/internal/checktypes"

// BoolCheckerProvider provides checks on type bool.
type BoolCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten bool is equal to the target.
func (p BoolCheckerProvider) Is(tar bool) check.Checker[bool] {
	pass := func(got bool) bool { return got == tar }
	expl := func(label string, got any) string {
		return p.explain(label, tar, got)
	}
	return check.NewChecker(pass, expl)
}
