package providers_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/internal/providers"
)

func TestNumberCheckerProvider(t *testing.T) {
	checkFloat64 := providers.NumberCheckerProvider[float64]{}

	const (
		n   = 42.
		inf = n - 1
		sup = n + 1
	)
	var (
		nstr   = fmt.Sprint(n)
		infstr = fmt.Sprint(inf)
		supstr = fmt.Sprint(sup)
	)

	t.Run("Is pass", func(t *testing.T) {
		c := checkFloat64.Is(n)
		assertPassChecker(t, "Number.Is", c, n)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := checkFloat64.Is(inf)
		assertFailChecker(t, "Number.Is", c, n, makeExpl(infstr, nstr))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := checkFloat64.Not(-1, 314, -n)
		assertPassChecker(t, "Number.Not", c, n)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := checkFloat64.Not(-1, 314, n, 1618)
		assertFailChecker(t, "Number.Not", c, n, makeExpl("not "+nstr, nstr))
	})

	t.Run("LT pass", func(t *testing.T) {
		c := checkFloat64.LT(sup)
		assertPassChecker(t, "Number.LT", c, n)
	})

	t.Run("LT fail", func(t *testing.T) {
		c := checkFloat64.LT(inf)
		assertFailChecker(t, "Number.LT", c, n, makeExpl("< "+infstr, nstr))
		c = checkFloat64.LT(n)
		assertFailChecker(t, "Number.LT", c, n, makeExpl("< "+nstr, nstr))
	})

	t.Run("LTE pass", func(t *testing.T) {
		c := checkFloat64.LTE(sup)
		assertPassChecker(t, "Number.LTE", c, n)
		c = checkFloat64.LTE(n)
		assertPassChecker(t, "Number.LTE", c, n)
	})

	t.Run("LTE fail", func(t *testing.T) {
		c := checkFloat64.LTE(inf)
		assertFailChecker(t, "Number.LTE", c, n, makeExpl("<= "+infstr, nstr))
	})

	t.Run("GT pass", func(t *testing.T) {
		c := checkFloat64.GT(inf)
		assertPassChecker(t, "Number.GT", c, n)
	})

	t.Run("GT fail", func(t *testing.T) {
		c := checkFloat64.GT(sup)
		assertFailChecker(t, "Number.GT", c, n, makeExpl("> "+supstr, nstr))
		c = checkFloat64.GT(n)
		assertFailChecker(t, "Number.GT", c, n, makeExpl("> "+nstr, nstr))
	})

	t.Run("GTE pass", func(t *testing.T) {
		c := checkFloat64.GTE(inf)
		assertPassChecker(t, "Number.GTE", c, n)
		c = checkFloat64.GTE(n)
		assertPassChecker(t, "Number.GTE", c, n)
	})

	t.Run("GTE fail", func(t *testing.T) {
		c := checkFloat64.GTE(sup)
		assertFailChecker(t, "Number.GTE", c, n, makeExpl(">= "+supstr, nstr))
	})

	t.Run("InRange pass", func(t *testing.T) {
		c := checkFloat64.InRange(inf, sup)
		assertPassChecker(t, "Number.InRange", c, n)

		c = checkFloat64.InRange(n, n)
		assertPassChecker(t, "Number.InRange", c, n)
	})

	t.Run("InRange fail", func(t *testing.T) {
		c := checkFloat64.InRange(sup, sup+1)
		assertFailChecker(t, "Number.InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, sup+1),
			nstr,
		))

		c = checkFloat64.InRange(sup, inf)
		assertFailChecker(t, "Number.InRange", c, n, makeExpl(
			fmt.Sprintf("in range [%v:%v]", sup, inf),
			nstr,
		))
	})

	t.Run("OutRange pass", func(t *testing.T) {
		c := checkFloat64.OutRange(sup, sup+1)
		assertPassChecker(t, "Number.OutRange", c, n)

		c = checkFloat64.OutRange(sup, inf)
		assertPassChecker(t, "Number.OutRange", c, n)
	})

	t.Run("OutRange fail", func(t *testing.T) {
		c := checkFloat64.OutRange(inf, sup)
		assertFailChecker(t, "Number.OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", inf, sup),
			nstr,
		))

		c = checkFloat64.OutRange(n, n)
		assertFailChecker(t, "Number.OutRange", c, n, makeExpl(
			fmt.Sprintf("not in range [%v:%v]", n, n),
			nstr,
		))
	})
}
