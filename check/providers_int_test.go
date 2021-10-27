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
		assertPassIntChecker(t, "Is", c, n)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Int.Is(inf)
		assertFailIntChecker(t, "Is", c, n, makeExpl(infstr, nstr))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.Int.Not(-1, 314, -n)
		assertPassIntChecker(t, "Not", c, n)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.Int.Not(-1, 314, n, 1618)
		assertFailIntChecker(t, "Not", c, n, makeExpl("not "+nstr, nstr))
	})

	t.Run("LT pass", func(t *testing.T) {
		c := check.Int.LT(sup)
		assertPassIntChecker(t, "LT", c, n)
	})

	t.Run("LT fail", func(t *testing.T) {
		c := check.Int.LT(inf)
		assertFailIntChecker(t, "LT", c, n, makeExpl("< "+infstr, nstr))
		c = check.Int.LT(n)
		assertFailIntChecker(t, "LT", c, n, makeExpl("< "+nstr, nstr))
	})

	t.Run("LTE pass", func(t *testing.T) {
		c := check.Int.LTE(sup)
		assertPassIntChecker(t, "LTE", c, n)
		c = check.Int.LTE(n)
		assertPassIntChecker(t, "LTE", c, n)
	})

	t.Run("LTE fail", func(t *testing.T) {
		c := check.Int.LTE(inf)
		assertFailIntChecker(t, "LTE", c, n, makeExpl("<= "+infstr, nstr))
	})

	t.Run("GT pass", func(t *testing.T) {
		c := check.Int.GT(inf)
		assertPassIntChecker(t, "GT", c, n)
	})

	t.Run("GT fail", func(t *testing.T) {
		c := check.Int.GT(sup)
		assertFailIntChecker(t, "GT", c, n, makeExpl("> "+supstr, nstr))
		c = check.Int.GT(n)
		assertFailIntChecker(t, "GT", c, n, makeExpl("> "+nstr, nstr))
	})

	t.Run("GTE pass", func(t *testing.T) {
		c := check.Int.GTE(inf)
		assertPassIntChecker(t, "GTE", c, n)
		c = check.Int.GTE(n)
		assertPassIntChecker(t, "GTE", c, n)
	})

	t.Run("GTE fail", func(t *testing.T) {
		c := check.Int.GTE(sup)
		assertFailIntChecker(t, "GTE", c, n, makeExpl(">= "+supstr, nstr))
	})

	t.Run("InRange pass", func(t *testing.T) {
		c := check.Int.InRange(inf, sup)
		assertPassIntChecker(t, "InRange", c, n)

		c = check.Int.InRange(n, n)
		assertPassIntChecker(t, "InRange", c, n)
	})

	t.Run("InRange fail", func(t *testing.T) {
		c := check.Int.InRange(sup, sup+1)
		assertFailIntChecker(t, "InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, sup+1),
			nstr,
		))

		c = check.Int.InRange(sup, inf)
		assertFailIntChecker(t, "InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, inf),
			nstr,
		))
	})

	t.Run("OutRange pass", func(t *testing.T) {
		c := check.Int.OutRange(sup, sup+1)
		assertPassIntChecker(t, "OutRange", c, n)

		c = check.Int.OutRange(sup, inf)
		assertPassIntChecker(t, "OutRange", c, n)
	})

	t.Run("OutRange fail", func(t *testing.T) {
		c := check.Int.OutRange(inf, sup)
		assertFailIntChecker(t, "OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", inf, sup),
			nstr,
		))

		c = check.Int.OutRange(n, n)
		assertFailIntChecker(t, "OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", n, n),
			nstr,
		))
	})
}

// Helpers

func assertPassIntChecker(t *testing.T, method string, c check.IntChecker, n int) {
	t.Helper()
	if !c.Pass(n) {
		failIntCheckerTest(t, true, method, n, c.Explain)
	}
}

func assertFailIntChecker(t *testing.T, method string, c check.IntChecker, n int, expexpl string) {
	t.Helper()
	if c.Pass(n) {
		failIntCheckerTest(t, false, method, n, c.Explain)
	}
	assertGoodExplain(t, c, n, expexpl)
}

func failIntCheckerTest(t *testing.T, expPass bool, method string, n int, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Int."+method, explain("Int value", n))
}
