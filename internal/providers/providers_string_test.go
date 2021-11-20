package providers_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/internal/providers"
)

func TestStringCheckerProvider(t *testing.T) {
	checkString := providers.StringCheckerProvider{}

	const (
		s   = "sator arepo tenet opera rotas"
		sub = "tenet"
		exp = s + "."
	)

	t.Run("Is pass", func(t *testing.T) {
		c := checkString.Is(s)
		assertPassChecker(t, "String.Is", c, s)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := checkString.Is(exp)
		assertFailChecker(t, "String.Is", c, s, makeExpl(exp, s))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := checkString.Not("hello", sub, exp)
		assertPassChecker(t, "String.Not", c, s)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := checkString.Not("hello", sub, s, exp)
		assertFailChecker(t, "String.Not", c, s, makeExpl("not "+s, s))
	})

	t.Run("Len pass", func(t *testing.T) {
		c := checkString.Len(check.Int.Is(len(s)))
		assertPassChecker(t, "String.Len", c, s)
	})

	t.Run("Len fail", func(t *testing.T) {
		gotlen := len(s)
		explen := gotlen + 1
		c := checkString.Len(check.Int.Is(explen))
		assertFailChecker(t, "String.Len", c, s, makeExpl(
			"length to pass Checker[int]",
			"explanation: length:\n"+makeExpl(
				fmt.Sprint(explen),
				fmt.Sprint(gotlen),
			),
		))
	})

	t.Run("Match pass", func(t *testing.T) {
		c := checkString.Match(regexp.MustCompile(`(?i)\sTENET\s`))
		assertPassChecker(t, "String.Match", c, s)
	})

	t.Run("Match fail", func(t *testing.T) {
		r := regexp.MustCompile(`\sTENET\s`)
		c := checkString.Match(r)
		assertFailChecker(t, "String.Match", c, s,
			makeExpl("to match regexp "+r.String(), s),
		)
	})

	t.Run("NotMatch pass", func(t *testing.T) {
		c := checkString.NotMatch(regexp.MustCompile(`\sTENET\s`))
		assertPassChecker(t, "String.NotMatch", c, s)
	})

	t.Run("NotMatch fail", func(t *testing.T) {
		r := regexp.MustCompile(`(?i)\sTENET\s`)
		c := checkString.NotMatch(r)
		assertFailChecker(t, "String.NotMatch", c, s, makeExpl(
			"not to match regexp "+r.String(),
			s,
		))
	})

	t.Run("Contains pass", func(t *testing.T) {
		c := checkString.Contains(sub)
		assertPassChecker(t, "String.Contains", c, s)
		c = checkString.Contains(s)
		assertPassChecker(t, "String.Contains", c, s)
	})

	t.Run("Contains fail", func(t *testing.T) {
		c := checkString.Contains(exp)
		assertFailChecker(t, "String.Contains", c, s, makeExpl(
			"to contain substring "+exp,
			s,
		))
	})

	t.Run("NotContains pass", func(t *testing.T) {
		c := checkString.NotContains(exp)
		assertPassChecker(t, "String.NotContains", c, s)
	})

	t.Run("NotContains fail", func(t *testing.T) {
		c := checkString.NotContains(sub)
		assertFailChecker(t, "String.NotContains", c, s, makeExpl(
			"not to contain substring "+sub,
			s,
		))
		c = checkString.NotContains(s)
		assertFailChecker(t, "String.NotContains", c, s, makeExpl(
			"not to contain substring "+s,
			s,
		))
	})
}
