package check_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestStringCheckerProvider(t *testing.T) {
	const (
		s   = "sator arepo tenet opera rotas"
		sub = "tenet"
		exp = s + "."
	)

	t.Run("Is pass", func(t *testing.T) {
		c := check.String.Is(s)
		assertPassStringChecker(t, "Is", c, s)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.String.Is(exp)
		assertFailStringChecker(t, "Is", c, s, makeExpl(exp, s))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.String.Not("hello", sub, exp)
		assertPassStringChecker(t, "Not", c, s)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.String.Not("hello", sub, s, exp)
		assertFailStringChecker(t, "Not", c, s, makeExpl("not "+s, s))
	})

	t.Run("Len pass", func(t *testing.T) {
		c := check.String.Len(check.Int.Is(len(s)))
		assertPassStringChecker(t, "Len", c, s)
	})

	t.Run("Len fail", func(t *testing.T) {
		gotlen := len(s)
		explen := gotlen + 1
		c := check.String.Len(check.Int.Is(explen))
		assertFailStringChecker(t, "Len", c, s, makeExpl(
			"length to pass IntChecker",
			"explanation: length:\n"+makeExpl(
				fmt.Sprint(explen),
				fmt.Sprint(gotlen),
			),
		))
	})

	t.Run("Match pass", func(t *testing.T) {
		c := check.String.Match(regexp.MustCompile(`(?i)\sTENET\s`))
		assertPassStringChecker(t, "Match", c, s)
	})

	t.Run("Match fail", func(t *testing.T) {
		r := regexp.MustCompile(`\sTENET\s`)
		c := check.String.Match(r)
		assertFailStringChecker(t, "Match", c, s,
			makeExpl("to match regexp "+r.String(), s),
		)
	})

	t.Run("NotMatch pass", func(t *testing.T) {
		c := check.String.NotMatch(regexp.MustCompile(`\sTENET\s`))
		assertPassStringChecker(t, "NotMatch", c, s)
	})

	t.Run("NotMatch fail", func(t *testing.T) {
		r := regexp.MustCompile(`(?i)\sTENET\s`)
		c := check.String.NotMatch(r)
		assertFailStringChecker(t, "NotMatch", c, s, makeExpl(
			"not to match regexp "+r.String(),
			s,
		))
	})

	t.Run("Contains pass", func(t *testing.T) {
		c := check.String.Contains(sub)
		assertPassStringChecker(t, "Contains", c, s)
		c = check.String.Contains(s)
		assertPassStringChecker(t, "Contains", c, s)
	})

	t.Run("Contains fail", func(t *testing.T) {
		c := check.String.Contains(exp)
		assertFailStringChecker(t, "Contains", c, s, makeExpl(
			"to contain substring "+exp,
			s,
		))
	})

	t.Run("NotContains pass", func(t *testing.T) {
		c := check.String.NotContains(exp)
		assertPassStringChecker(t, "NotContains", c, s)
	})

	t.Run("NotContains fail", func(t *testing.T) {
		c := check.String.NotContains(sub)
		assertFailStringChecker(t, "NotContains", c, s, makeExpl(
			"not to contain substring "+sub,
			s,
		))
		c = check.String.NotContains(s)
		assertFailStringChecker(t, "NotContains", c, s, makeExpl(
			"not to contain substring "+s,
			s,
		))
	})
}

// Helpers

func assertPassStringChecker(t *testing.T, method string, c check.StringChecker, in string) {
	t.Helper()
	if !c.Pass(in) {
		failStringCheckerTest(t, true, method, in, c.Explain)
	}
}

func assertFailStringChecker(t *testing.T, method string, c check.StringChecker, in, expexpl string) {
	t.Helper()
	if c.Pass(in) {
		failStringCheckerTest(t, false, method, in, c.Explain)
	}
	assertGoodExplain(t, c, in, expexpl)
}

func failStringCheckerTest(t *testing.T, expPass bool, method, in string, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "String."+method, explain("String value", in))
}
