package check

import (
	"fmt"
	"regexp"
	"strings"
)

// stringCheckerProvider provides checks on type string.
type stringCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten string is equal to the target.
func (p stringCheckerProvider) Is(tar string) StringChecker {
	pass := func(got string) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return p.explain(label, tar, got)
	}
	return NewStringChecker(pass, expl)
}

// Not checks the gotten string is not equal to the target.
func (p stringCheckerProvider) Not(values ...string) StringChecker {
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
	return NewStringChecker(pass, expl)
}

// Len checks the gotten string's length passes the given IntChecker.
func (p stringCheckerProvider) Len(c IntChecker) StringChecker {
	pass := func(got string) bool { return c.Pass(len(got)) }
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			"length to pass IntChecker",
			c.Explain("length", len(got.(string))),
		)
	}
	return NewStringChecker(pass, expl)
}

// Match checks the gotten string matches the given regexp.
func (p stringCheckerProvider) Match(rgx *regexp.Regexp) StringChecker {
	pass := func(got string) bool { return rgx.MatchString(got) }
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			fmt.Sprintf("to match regexp %s", rgx.String()),
			got,
		)
	}
	return NewStringChecker(pass, expl)
}

// NotMatch checks the gotten string do not match the given regexp.
func (p stringCheckerProvider) NotMatch(rgx *regexp.Regexp) StringChecker {
	pass := func(got string) bool { return !rgx.MatchString(got) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label,
			fmt.Sprintf("to match regexp %s", rgx.String()),
			got,
		)
	}
	return NewStringChecker(pass, expl)
}

// Contains checks the gotten string contains the target substring.
func (p stringCheckerProvider) Contains(sub string) StringChecker {
	pass := func(got string) bool { return strings.Contains(got, sub) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, "to contain substring "+sub, got)
	}
	return NewStringChecker(pass, expl)
}

// NotContains checks the gotten string do not contain the target
// substring.
func (p stringCheckerProvider) NotContains(sub string) StringChecker {
	pass := func(got string) bool { return !strings.Contains(got, sub) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, "to contain substring "+sub, got)
	}
	return NewStringChecker(pass, expl)
}
