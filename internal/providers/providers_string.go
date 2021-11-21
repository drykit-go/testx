package providers

import (
	"fmt"
	"regexp"
	"strings"

	check "github.com/drykit-go/testx/internal/checktypes"
)

// StringCheckerProvider provides checks on type string.
type StringCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten string is equal to the target.
func (p StringCheckerProvider) Is(tar string) check.Checker[string] {
	pass := func(got string) bool { return got == tar }
	expl := func(label string, got any) string {
		return p.explain(label, tar, got)
	}
	return check.NewChecker(pass, expl)
}

// Not checks the gotten string is not equal to the target.
func (p StringCheckerProvider) Not(values ...string) check.Checker[string] {
	var match string
	pass := func(got string) bool {
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
	return check.NewChecker(pass, expl)
}

// Len checks the gotten string's length passes the given Checker[int].
func (p StringCheckerProvider) Len(c check.Checker[int]) check.Checker[string] {
	pass := func(got string) bool { return c.Pass(len(got)) }
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			"length to pass Checker[int]",
			c.Explain("length", len(got.(string))),
		)
	}
	return check.NewChecker(pass, expl)
}

// Match checks the gotten string matches the given regexp.
func (p StringCheckerProvider) Match(rgx *regexp.Regexp) check.Checker[string] {
	pass := func(got string) bool { return rgx.MatchString(got) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("to match regexp %s", rgx.String()),
			got,
		)
	}
	return check.NewChecker(pass, expl)
}

// NotMatch checks the gotten string do not match the given regexp.
func (p StringCheckerProvider) NotMatch(rgx *regexp.Regexp) check.Checker[string] {
	pass := func(got string) bool { return !rgx.MatchString(got) }
	expl := func(label string, got any) string {
		return p.explainNot(label,
			fmt.Sprintf("to match regexp %s", rgx.String()),
			got,
		)
	}
	return check.NewChecker(pass, expl)
}

// Contains checks the gotten string contains the target substring.
func (p StringCheckerProvider) Contains(sub string) check.Checker[string] {
	pass := func(got string) bool { return strings.Contains(got, sub) }
	expl := func(label string, got any) string {
		return p.explain(label, "to contain substring "+sub, got)
	}
	return check.NewChecker(pass, expl)
}

// NotContains checks the gotten string do not contain the target
// substring.
func (p StringCheckerProvider) NotContains(sub string) check.Checker[string] {
	pass := func(got string) bool { return !strings.Contains(got, sub) }
	expl := func(label string, got any) string {
		return p.explainNot(label, "to contain substring "+sub, got)
	}
	return check.NewChecker(pass, expl)
}
