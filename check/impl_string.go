package check

import (
	"fmt"
	"regexp"
	"strings"
)

type stringCheckerProvider struct{}

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
