package check

import (
	"fmt"
	"regexp"
	"strings"
)

// stringCheckerProvider provides checks on type string.
type stringCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten string is equal to the target.
func (p stringCheckerProvider) Is(tar string) Checker[string] {
	pass := func(got string) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return p.explain(label, tar, got)
	}
	return NewChecker(pass, expl)
}

// Not checks the gotten string is not equal to the target.
func (p stringCheckerProvider) Not(values ...string) Checker[string] {
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
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, match, got)
	}
	return NewChecker(pass, expl)
}

// Len checks the gotten string's length passes the given Checker[int].
func (p stringCheckerProvider) Len(c Checker[int]) Checker[string] {
	pass := func(got string) bool { return c.Pass(len(got)) }
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			"length to pass Checker[int]",
			c.Explain("length", len(got.(string))),
		)
	}
	return NewChecker(pass, expl)
}

// Match checks the gotten string matches the given regexp.
func (p stringCheckerProvider) Match(rgx *regexp.Regexp) Checker[string] {
	pass := func(got string) bool { return rgx.MatchString(got) }
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			fmt.Sprintf("to match regexp %s", rgx.String()),
			got,
		)
	}
	return NewChecker(pass, expl)
}

// NotMatch checks the gotten string do not match the given regexp.
func (p stringCheckerProvider) NotMatch(rgx *regexp.Regexp) Checker[string] {
	pass := func(got string) bool { return !rgx.MatchString(got) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label,
			fmt.Sprintf("to match regexp %s", rgx.String()),
			got,
		)
	}
	return NewChecker(pass, expl)
}

// Contains checks the gotten string contains the target substring.
func (p stringCheckerProvider) Contains(sub string) Checker[string] {
	pass := func(got string) bool { return strings.Contains(got, sub) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, "to contain substring "+sub, got)
	}
	return NewChecker(pass, expl)
}

// NotContains checks the gotten string do not contain the target
// substring.
func (p stringCheckerProvider) NotContains(sub string) Checker[string] {
	pass := func(got string) bool { return !strings.Contains(got, sub) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, "to contain substring "+sub, got)
	}
	return NewChecker(pass, expl)
}
