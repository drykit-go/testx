package check

import (
	"fmt"
	"regexp"
	"strings"
)

// stringCheckerProvider provides checks on type string.
type stringCheckerProvider struct{}

// Is checks the gotten string is equal to the target.
func (stringCheckerProvider) Is(tar string) StringChecker {
	pass := func(got string) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to equal %v, got %v",
			label, tar, got,
		)
	}
	return NewStringChecker(pass, expl)
}

// Len checks the gotten string's length passes the given IntChecker.
func (stringCheckerProvider) Len(c IntChecker) StringChecker {
	pass := func(got string) bool { return c.Pass(len(got)) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"unexpected %s length: %s",
			label, c.Explain(label, len(got.(string))),
		)
	}
	return NewStringChecker(pass, expl)
}

// Match checks the gotten string matches the given regexp.
func (stringCheckerProvider) Match(rgx *regexp.Regexp) StringChecker {
	pass := func(got string) bool { return rgx.MatchString(got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to match regexp %s, got %s",
			label, rgx.String(), got,
		)
	}
	return NewStringChecker(pass, expl)
}

// NotMatch checks the gotten string do not match the given regexp.
func (stringCheckerProvider) NotMatch(rgx *regexp.Regexp) StringChecker {
	pass := func(got string) bool { return !rgx.MatchString(got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to match regexp %s, got %s",
			label, rgx.String(), got,
		)
	}
	return NewStringChecker(pass, expl)
}

// Contains checks the gotten string contains the target substring.
func (stringCheckerProvider) Contains(sub string) StringChecker {
	pass := func(got string) bool { return strings.Contains(got, sub) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to contain substring %s, got %s",
			label, sub, got,
		)
	}
	return NewStringChecker(pass, expl)
}

// NotContains checks the gotten string do not contain the target
// substring.
func (stringCheckerProvider) NotContains(sub string) StringChecker {
	pass := func(got string) bool { return !strings.Contains(got, sub) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to contain substring %s, got %s",
			label, sub, got,
		)
	}
	return NewStringChecker(pass, expl)
}