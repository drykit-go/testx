package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestFloat64CheckerProvider(t *testing.T) {
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
		c := check.Float64.Is(n)
		assertPassFloat64Checker(t, "Is", c, n)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Float64.Is(inf)
		assertFailFloat64Checker(t, "Is", c, n, makeExpl(infstr, nstr))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.Float64.Not(-1, 314, -n)
		assertPassFloat64Checker(t, "Not", c, n)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.Float64.Not(-1, 314, n, 1618)
		assertFailFloat64Checker(t, "Not", c, n, makeExpl("not "+nstr, nstr))
	})

	t.Run("LT pass", func(t *testing.T) {
		c := check.Float64.LT(sup)
		assertPassFloat64Checker(t, "LT", c, n)
	})

	t.Run("LT fail", func(t *testing.T) {
		c := check.Float64.LT(inf)
		assertFailFloat64Checker(t, "LT", c, n, makeExpl("< "+infstr, nstr))
		c = check.Float64.LT(n)
		assertFailFloat64Checker(t, "LT", c, n, makeExpl("< "+nstr, nstr))
	})

	t.Run("LTE pass", func(t *testing.T) {
		c := check.Float64.LTE(sup)
		assertPassFloat64Checker(t, "LTE", c, n)
		c = check.Float64.LTE(n)
		assertPassFloat64Checker(t, "LTE", c, n)
	})

	t.Run("LTE fail", func(t *testing.T) {
		c := check.Float64.LTE(inf)
		assertFailFloat64Checker(t, "LTE", c, n, makeExpl("<= "+infstr, nstr))
	})

	t.Run("GT pass", func(t *testing.T) {
		c := check.Float64.GT(inf)
		assertPassFloat64Checker(t, "GT", c, n)
	})

	t.Run("GT fail", func(t *testing.T) {
		c := check.Float64.GT(sup)
		assertFailFloat64Checker(t, "GT", c, n, makeExpl("> "+supstr, nstr))
		c = check.Float64.GT(n)
		assertFailFloat64Checker(t, "GT", c, n, makeExpl("> "+nstr, nstr))
	})

	t.Run("GTE pass", func(t *testing.T) {
		c := check.Float64.GTE(inf)
		assertPassFloat64Checker(t, "GTE", c, n)
		c = check.Float64.GTE(n)
		assertPassFloat64Checker(t, "GTE", c, n)
	})

	t.Run("GTE fail", func(t *testing.T) {
		c := check.Float64.GTE(sup)
		assertFailFloat64Checker(t, "GTE", c, n, makeExpl(">= "+supstr, nstr))
	})

	t.Run("InRange pass", func(t *testing.T) {
		c := check.Float64.InRange(inf, sup)
		assertPassFloat64Checker(t, "InRange", c, n)

		c = check.Float64.InRange(n, n)
		assertPassFloat64Checker(t, "InRange", c, n)
	})

	t.Run("InRange fail", func(t *testing.T) {
		c := check.Float64.InRange(sup, sup+1)
		assertFailFloat64Checker(t, "InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, sup+1),
			nstr,
		))

		c = check.Float64.InRange(sup, inf)
		assertFailFloat64Checker(t, "InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, inf),
			nstr,
		))
	})

	t.Run("OutRange pass", func(t *testing.T) {
		c := check.Float64.OutRange(sup, sup+1)
		assertPassFloat64Checker(t, "OutRange", c, n)

		c = check.Float64.OutRange(sup, inf)
		assertPassFloat64Checker(t, "OutRange", c, n)
	})

	t.Run("OutRange fail", func(t *testing.T) {
		c := check.Float64.OutRange(inf, sup)
		assertFailFloat64Checker(t, "OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", inf, sup),
			nstr,
		))

		c = check.Float64.OutRange(n, n)
		assertFailFloat64Checker(t, "OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", n, n),
			nstr,
		))
	})
}

// Helpers

func assertPassFloat64Checker(t *testing.T, method string, c check.Float64Checker, n float64) {
	t.Helper()
	if !c.Pass(n) {
		failFloat64CheckerTest(t, true, method, n, c.Explain)
	}
}

func assertFailFloat64Checker(t *testing.T, method string, c check.Float64Checker, n float64, expexpl string) {
	t.Helper()
	if c.Pass(n) {
		failFloat64CheckerTest(t, false, method, n, c.Explain)
	}
	assertGoodExplain(t, c, n, expexpl)
}

func failFloat64CheckerTest(t *testing.T, expPass bool, method string, n float64, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Float64."+method, explain("Float64 value", n))
}
