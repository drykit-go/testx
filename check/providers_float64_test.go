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
		assertPassChecker(t, "Float64.Is", c, n)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Float64.Is(inf)
		assertFailChecker(t, "Float64.Is", c, n, makeExpl(infstr, nstr))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.Float64.Not(-1, 314, -n)
		assertPassChecker(t, "Float64.Not", c, n)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.Float64.Not(-1, 314, n, 1618)
		assertFailChecker(t, "Float64.Not", c, n, makeExpl("not "+nstr, nstr))
	})

	t.Run("LT pass", func(t *testing.T) {
		c := check.Float64.LT(sup)
		assertPassChecker(t, "Float64.LT", c, n)
	})

	t.Run("LT fail", func(t *testing.T) {
		c := check.Float64.LT(inf)
		assertFailChecker(t, "Float64.LT", c, n, makeExpl("< "+infstr, nstr))
		c = check.Float64.LT(n)
		assertFailChecker(t, "Float64.LT", c, n, makeExpl("< "+nstr, nstr))
	})

	t.Run("LTE pass", func(t *testing.T) {
		c := check.Float64.LTE(sup)
		assertPassChecker(t, "Float64.LTE", c, n)
		c = check.Float64.LTE(n)
		assertPassChecker(t, "Float64.LTE", c, n)
	})

	t.Run("LTE fail", func(t *testing.T) {
		c := check.Float64.LTE(inf)
		assertFailChecker(t, "Float64.LTE", c, n, makeExpl("<= "+infstr, nstr))
	})

	t.Run("GT pass", func(t *testing.T) {
		c := check.Float64.GT(inf)
		assertPassChecker(t, "Float64.GT", c, n)
	})

	t.Run("GT fail", func(t *testing.T) {
		c := check.Float64.GT(sup)
		assertFailChecker(t, "Float64.GT", c, n, makeExpl("> "+supstr, nstr))
		c = check.Float64.GT(n)
		assertFailChecker(t, "Float64.GT", c, n, makeExpl("> "+nstr, nstr))
	})

	t.Run("GTE pass", func(t *testing.T) {
		c := check.Float64.GTE(inf)
		assertPassChecker(t, "Float64.GTE", c, n)
		c = check.Float64.GTE(n)
		assertPassChecker(t, "Float64.GTE", c, n)
	})

	t.Run("GTE fail", func(t *testing.T) {
		c := check.Float64.GTE(sup)
		assertFailChecker(t, "Float64.GTE", c, n, makeExpl(">= "+supstr, nstr))
	})

	t.Run("InRange pass", func(t *testing.T) {
		c := check.Float64.InRange(inf, sup)
		assertPassChecker(t, "Float64.InRange", c, n)

		c = check.Float64.InRange(n, n)
		assertPassChecker(t, "Float64.InRange", c, n)
	})

	t.Run("InRange fail", func(t *testing.T) {
		c := check.Float64.InRange(sup, sup+1)
		assertFailChecker(t, "Float64.InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, sup+1),
			nstr,
		))

		c = check.Float64.InRange(sup, inf)
		assertFailChecker(t, "Float64.InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, inf),
			nstr,
		))
	})

	t.Run("OutRange pass", func(t *testing.T) {
		c := check.Float64.OutRange(sup, sup+1)
		assertPassChecker(t, "Float64.OutRange", c, n)

		c = check.Float64.OutRange(sup, inf)
		assertPassChecker(t, "Float64.OutRange", c, n)
	})

	t.Run("OutRange fail", func(t *testing.T) {
		c := check.Float64.OutRange(inf, sup)
		assertFailChecker(t, "Float64.OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", inf, sup),
			nstr,
		))

		c = check.Float64.OutRange(n, n)
		assertFailChecker(t, "Float64.OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", n, n),
			nstr,
		))
	})
}
