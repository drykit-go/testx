package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestValueCheckerProvider(t *testing.T) {
	type Thing struct {
		Name string
	}

	var (
		vorig  = Thing{"hi"}
		vsame  = Thing{"hi"}
		vzero  = Thing{}
		badval = Thing{"hello"}
		badtyp = struct{ Name string }{"hi"}

		emptyMap   map[int]bool
		emptySlice []float32

		zeros   = []any{0, "", 0i + 0, vzero, emptyMap, emptySlice}
		nozeros = []any{1, "hi", 0i + 1, vorig, map[int]bool{}, []float32{}}
	)

	itf := func(v any) any {
		return v
	}

	t.Run("Is pass", func(t *testing.T) {
		c := check.Value.Is(vsame)
		assertPassChecker(t, "Value.Is", c, itf(vorig))
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Value.Is(badval)
		assertFailChecker(t, "Value.Is", c, itf(vorig), makeExpl("{hello}", "{hi}"))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.Value.Not(badval, badtyp)
		assertPassChecker(t, "Value.Not", c, itf(vorig))
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.Value.Not(badval, vsame, badtyp)
		assertFailChecker(t, "Value.Not", c, itf(vorig), makeExpl("not {hi}", "{hi}"))
	})

	t.Run("IsZero pass", func(t *testing.T) {
		c := check.Value.IsZero()
		for _, z := range zeros {
			assertPassChecker(t, "Value.IsZero", c, z)
		}
	})

	t.Run("IsZero fail", func(t *testing.T) {
		c := check.Value.IsZero()
		for _, nz := range nozeros {
			assertFailChecker(t, "Value.IsZero", c, nz, makeExpl(
				"to be a zero value",
				fmt.Sprint(nz),
			))
		}
	})

	t.Run("NotZero pass", func(t *testing.T) {
		c := check.Value.NotZero()
		for _, nz := range nozeros {
			assertPassChecker(t, "Value.NotZero", c, nz)
		}
	})

	t.Run("NotZero fail", func(t *testing.T) {
		c := check.Value.NotZero()
		for _, z := range zeros {
			assertFailChecker(t, "Value.NotZero", c, z, makeExpl(
				"not to be a zero value",
				fmt.Sprint(z),
			))
		}
	})

	isEvenInt := func(n any) bool {
		nint, ok := n.(int)
		return ok && nint&1 == 0
	}

	t.Run("Custom pass", func(t *testing.T) {
		c := check.Value.Custom("", isEvenInt)
		assertPassChecker(t, "Value.Custom", c, 42)
	})

	t.Run("Custom fail", func(t *testing.T) {
		c := check.Value.Custom("even int", isEvenInt)
		assertFailChecker(t, "Value.Custom", c, -1, makeExpl("even int", "-1"))
	})

	t.Run("SameJSON pass", func(t *testing.T) {
		mapequiv := map[string]any{
			"Name": "hi",
		}
		c := check.Value.SameJSON(mapequiv)
		assertPassChecker(t, "Value.SameJSON", c, itf(vorig))
	})

	t.Run("SameJSON fail", func(t *testing.T) {
		mapdiff := map[string]any{
			"Name": "bad",
		}
		c := check.Value.SameJSON(mapdiff)
		assertFailChecker(t, "Value.SameJSON", c, itf(vorig), makeExpl(
			"json data: map[Name:bad]",
			"json data: map[Name:hi]",
		))
	})
}
