package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

func TestSliceCheckerProvider(t *testing.T) {
	s := []any{"hello", 42, "Marcel Patulacci", []float32{3.14}}
	// FIXME: remove forced conversion
	itf := func(v []any) any {
		return v
	}

	t.Run("Len pass", func(t *testing.T) {
		c := check.Slice.Len(check.Int.Is(4))
		assertPassChecker(t, "Slice.Len", c, itf(s))
	})

	t.Run("Len fail", func(t *testing.T) {
		c := check.Slice.Len(check.Int.Not(4))
		assertFailChecker(t, "Slice.Len", c, itf(s), makeExpl(
			"length to pass Checker[int]",
			"explanation: length:\n"+makeExpl("not 4", "4"),
		))
	})

	t.Run("Cap pass", func(t *testing.T) {
		c := check.Slice.Cap(check.Int.Is(4))
		assertPassChecker(t, "Slice.Cap", c, itf(s))
	})

	t.Run("Cap fail", func(t *testing.T) {
		c := check.Slice.Cap(check.Int.Not(4))
		assertFailChecker(t, "Slice.Cap", c, itf(s), makeExpl(
			"capacity to pass Checker[int]",
			"explanation: capacity:\n"+makeExpl("not 4", "4"),
		))
	})

	t.Run("HasValues pass", func(t *testing.T) {
		c := check.Slice.HasValues("hello", 42, []float32{3.14})
		assertPassChecker(t, "Slice.HasValues", c, itf(s))
	})

	t.Run("HasValues fail", func(t *testing.T) {
		c := check.Slice.HasValues([]float64{3.14})
		assertFailChecker(t, "Slice.HasValues", c, itf(s), makeExpl(
			"to have values [[3.14]]",
			fmt.Sprint(s),
		))
	})

	t.Run("HasNotValues pass", func(t *testing.T) {
		c := check.Slice.HasNotValues("hi", -1, []float64{3.14})
		assertPassChecker(t, "Slice.HasNotValues", c, itf(s))
	})

	t.Run("HasNotValues fail", func(t *testing.T) {
		c := check.Slice.HasNotValues("hi", -1, []float32{3.14})
		assertFailChecker(t, "Slice.HasNotValues", c, itf(s), makeExpl(
			"not to have values [[3.14]]",
			fmt.Sprint(s),
		))
	})

	// t.Run("CheckValues pass", func(t *testing.T) {
	// 	c := check.Slice.CheckValues(
	// 		checkconv.FromInt(check.Int.InRange(41, 43)),
	// 		func(_ int, v any) bool { _, ok := v.(int); return ok },
	// 	)
	// 	assertPassChecker(t, "Slice.CheckValues", c, itf(s))
	// })

	t.Run("CheckValues fail", func(t *testing.T) {
		c := check.Slice.CheckValues(
			checkconv.FromInt(check.Int.OutRange(41, 43)),
			func(_ int, v any) bool { _, ok := v.(int); return ok },
		)
		assertFailChecker(t, "Slice.CheckValues", c, itf(s), makeExpl(
			"values to pass Checker[any]",
			"explanation: values:\n"+makeExpl("not in range [41:43]", "[1:42]"),
		))
	})
}
