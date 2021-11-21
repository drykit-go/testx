package providers_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/internal/providers"
)

func TestSliceCheckerProvider(t *testing.T) {
	checkSlice := providers.SliceCheckerProvider[any]{}

	s := []any{"hello", 42, "Marcel Patulacci", []float32{3.14}}

	t.Run("Len pass", func(t *testing.T) {
		c := checkSlice.Len(check.Int.Is(4))
		assertPassChecker(t, "Slice.Len", c, s)
	})

	t.Run("Len fail", func(t *testing.T) {
		c := checkSlice.Len(check.Int.Not(4))
		assertFailChecker(t, "Slice.Len", c, s, makeExpl(
			"length to pass Checker[int]",
			"explanation: length:\n"+makeExpl("not 4", "4"),
		))
	})

	t.Run("Cap pass", func(t *testing.T) {
		c := checkSlice.Cap(check.Int.Is(4))
		assertPassChecker(t, "Slice.Cap", c, s)
	})

	t.Run("Cap fail", func(t *testing.T) {
		c := checkSlice.Cap(check.Int.Not(4))
		assertFailChecker(t, "Slice.Cap", c, s, makeExpl(
			"capacity to pass Checker[int]",
			"explanation: capacity:\n"+makeExpl("not 4", "4"),
		))
	})

	t.Run("HasValues pass", func(t *testing.T) {
		c := checkSlice.HasValues("hello", 42, []float32{3.14})
		assertPassChecker(t, "Slice.HasValues", c, s)
	})

	t.Run("HasValues fail", func(t *testing.T) {
		c := checkSlice.HasValues([]float64{3.14})
		assertFailChecker(t, "Slice.HasValues", c, s, makeExpl(
			"to have values [[3.14]]",
			fmt.Sprint(s),
		))
	})

	t.Run("HasNotValues pass", func(t *testing.T) {
		c := checkSlice.HasNotValues("hi", -1, []float64{3.14})
		assertPassChecker(t, "Slice.HasNotValues", c, s)
	})

	t.Run("HasNotValues fail", func(t *testing.T) {
		c := checkSlice.HasNotValues("hi", -1, []float32{3.14})
		assertFailChecker(t, "Slice.HasNotValues", c, s, makeExpl(
			"not to have values [[3.14]]",
			fmt.Sprint(s),
		))
	})

	t.Run("CheckValues pass", func(t *testing.T) {
		c := checkSlice.CheckValues(
			check.Wrap[int](check.Int.InRange(41, 43)),
			func(_ int, v any) bool { _, ok := v.(int); return ok },
		)
		assertPassChecker(t, "Slice.CheckValues", c, s)
	})

	t.Run("CheckValues fail", func(t *testing.T) {
		c := checkSlice.CheckValues(
			check.Wrap[int](check.Int.OutRange(41, 43)),
			func(_ int, v any) bool { _, ok := v.(int); return ok },
		)
		assertFailChecker(t, "Slice.CheckValues", c, s, makeExpl(
			"values to pass Checker[Elem]",
			"explanation: values:\n"+makeExpl("not in range [41:43]", "[1:42]"),
		))
	})
}
