package check

import (
	"fmt"
	"regexp"
	"strings"
)

// StringCheckerProvider provides checks on type string.
type StringCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten string is equal to the target.
func (p StringCheckerProvider) Is(tar string) Checker[string] {
	pass := func(got string) bool { return got == tar }
	expl := func(label string, got any) string {
		return p.explain(label, tar, got)
	}
	return NewChecker(pass, expl)
}

// Not checks the gotten string is not equal to the target.
func (p StringCheckerProvider) Not(values ...string) Checker[string] {
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
	return NewChecker(pass, expl)
}

// Len checks the gotten string's length passes the given Checker[int].
func (p StringCheckerProvider) Len(c Checker[int]) Checker[string] {
	pass := func(got string) bool { return c.Pass(len(got)) }
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			"length to pass Checker[int]",
			c.Explain("length", len(got.(string))),
		)
	}
	return NewChecker(pass, expl)
}

// Match checks the gotten string matches the given regexp.
func (p StringCheckerProvider) Match(rgx *regexp.Regexp) Checker[string] {
	pass := func(got string) bool { return rgx.MatchString(got) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("to match regexp %s", rgx.String()),
			got,
		)
	}
	return NewChecker(pass, expl)
}

// NotMatch checks the gotten string do not match the given regexp.
func (p StringCheckerProvider) NotMatch(rgx *regexp.Regexp) Checker[string] {
	pass := func(got string) bool { return !rgx.MatchString(got) }
	expl := func(label string, got any) string {
		return p.explainNot(label,
			fmt.Sprintf("to match regexp %s", rgx.String()),
			got,
		)
	}
	return NewChecker(pass, expl)
}

// Contains checks the gotten string contains the target substring.
func (p StringCheckerProvider) Contains(sub string) Checker[string] {
	pass := func(got string) bool { return strings.Contains(got, sub) }
	expl := func(label string, got any) string {
		return p.explain(label, "to contain substring "+sub, got)
	}
	return NewChecker(pass, expl)
}

// NotContains checks the gotten string do not contain the target
// substring.
func (p StringCheckerProvider) NotContains(sub string) Checker[string] {
	pass := func(got string) bool { return !strings.Contains(got, sub) }
	expl := func(label string, got any) string {
		return p.explainNot(label, "to contain substring "+sub, got)
	}
	return NewChecker(pass, expl)
}
