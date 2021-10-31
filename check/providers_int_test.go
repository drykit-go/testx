package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestIntCheckerProvider(t *testing.T) {
	const (
		n   = 42
		inf = n - 1
		sup = n + 1
	)
	var (
		nstr   = fmt.Sprint(n)
		infstr = fmt.Sprint(inf)
		supstr = fmt.Sprint(sup)
	)

	t.Run("Is pass", func(t *testing.T) {
		c := check.Int.Is(n)
		assertPassChecker(t, "Int.Is", c, n)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Int.Is(inf)
		assertFailChecker(t, "Int.Is", c, n, makeExpl(infstr, nstr))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.Int.Not(-1, 314, -n)
		assertPassChecker(t, "Int.Not", c, n)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.Int.Not(-1, 314, n, 1618)
		assertFailChecker(t, "Int.Not", c, n, makeExpl("not "+nstr, nstr))
	})

	t.Run("LT pass", func(t *testing.T) {
		c := check.Int.LT(sup)
		assertPassChecker(t, "Int.LT", c, n)
	})

	t.Run("LT fail", func(t *testing.T) {
		c := check.Int.LT(inf)
		assertFailChecker(t, "Int.LT", c, n, makeExpl("< "+infstr, nstr))
		c = check.Int.LT(n)
		assertFailChecker(t, "Int.LT", c, n, makeExpl("< "+nstr, nstr))
	})

	t.Run("LTE pass", func(t *testing.T) {
		c := check.Int.LTE(sup)
		assertPassChecker(t, "Int.LTE", c, n)
		c = check.Int.LTE(n)
		assertPassChecker(t, "Int.LTE", c, n)
	})

	t.Run("LTE fail", func(t *testing.T) {
		c := check.Int.LTE(inf)
		assertFailChecker(t, "Int.LTE", c, n, makeExpl("<= "+infstr, nstr))
	})

	t.Run("GT pass", func(t *testing.T) {
		c := check.Int.GT(inf)
		assertPassChecker(t, "Int.GT", c, n)
	})

	t.Run("GT fail", func(t *testing.T) {
		c := check.Int.GT(sup)
		assertFailChecker(t, "Int.GT", c, n, makeExpl("> "+supstr, nstr))
		c = check.Int.GT(n)
		assertFailChecker(t, "Int.GT", c, n, makeExpl("> "+nstr, nstr))
	})

	t.Run("GTE pass", func(t *testing.T) {
		c := check.Int.GTE(inf)
		assertPassChecker(t, "Int.GTE", c, n)
		c = check.Int.GTE(n)
		assertPassChecker(t, "Int.GTE", c, n)
	})

	t.Run("GTE fail", func(t *testing.T) {
		c := check.Int.GTE(sup)
		assertFailChecker(t, "Int.GTE", c, n, makeExpl(">= "+supstr, nstr))
	})

	t.Run("InRange pass", func(t *testing.T) {
		c := check.Int.InRange(inf, sup)
		assertPassChecker(t, "Int.InRange", c, n)

		c = check.Int.InRange(n, n)
		assertPassChecker(t, "Int.InRange", c, n)
	})

	t.Run("InRange fail", func(t *testing.T) {
		c := check.Int.InRange(sup, sup+1)
		assertFailChecker(t, "Int.InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, sup+1),
			nstr,
		))

		c = check.Int.InRange(sup, inf)
		assertFailChecker(t, "Int.InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, inf),
			nstr,
		))
	})

	t.Run("OutRange pass", func(t *testing.T) {
		c := check.Int.OutRange(sup, sup+1)
		assertPassChecker(t, "Int.OutRange", c, n)

		c = check.Int.OutRange(sup, inf)
		assertPassChecker(t, "Int.OutRange", c, n)
	})

	t.Run("OutRange fail", func(t *testing.T) {
		c := check.Int.OutRange(inf, sup)
		assertFailChecker(t, "Int.OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", inf, sup),
			nstr,
		))

		c = check.Int.OutRange(n, n)
		assertFailChecker(t, "Int.OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", n, n),
			nstr,
		))
	})
}
