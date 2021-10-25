package check_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

func TestSliceCheckerProvider(t *testing.T) {
	s := []interface{}{"hello", 42, "Marcel Patulacci", []float32{3.14}}

	t.Run("Len pass", func(t *testing.T) {
		c := check.Slice.Len(check.Int.Is(4))
		assertPassSliceChecker(t, "Len", c, s)
	})

	t.Run("Len fail", func(t *testing.T) {
		c := check.Slice.Len(check.Int.Not(4))
		assertFailSliceChecker(t, "Len", c, s)
	})

	t.Run("Cap pass", func(t *testing.T) {
		c := check.Slice.Cap(check.Int.Is(4))
		assertPassSliceChecker(t, "Cap", c, s)
	})

	t.Run("Cap fail", func(t *testing.T) {
		c := check.Slice.Cap(check.Int.Not(4))
		assertFailSliceChecker(t, "Cap", c, s)
	})

	t.Run("HasValues pass", func(t *testing.T) {
		c := check.Slice.HasValues("hello", 42, []float32{3.14})
		assertPassSliceChecker(t, "HasValues", c, s)
	})

	t.Run("HasValues fail", func(t *testing.T) {
		c := check.Slice.HasValues([]float64{3.14})
		assertFailSliceChecker(t, "HasValues", c, s)
	})

	t.Run("HasNotValues pass", func(t *testing.T) {
		c := check.Slice.HasNotValues("hi", -1, []float64{3.14})
		assertPassSliceChecker(t, "HasNotValues", c, s)
	})

	t.Run("HasNotValues fail", func(t *testing.T) {
		c := check.Slice.HasNotValues("hi", -1, []float32{3.14})
		assertFailSliceChecker(t, "HasNotValues", c, s)
	})

	t.Run("CheckValues pass", func(t *testing.T) {
		c := check.Slice.CheckValues(
			checkconv.FromInt(check.Int.InRange(41, 43)),
			func(_ int, v interface{}) bool { _, ok := v.(int); return ok },
		)
		assertPassSliceChecker(t, "CheckValues", c, s)
	})

	t.Run("CheckValues fail", func(t *testing.T) {
		c := check.Slice.CheckValues(
			checkconv.FromInt(check.Int.OutRange(41, 43)),
			func(_ int, v interface{}) bool { _, ok := v.(int); return ok },
		)
		assertFailSliceChecker(t, "CheckValues", c, s)
	})
}

// Helpers

func assertPassSliceChecker(t *testing.T, method string, c check.ValueChecker, slc interface{}) {
	t.Helper()
	if !c.Pass(slc) {
		failSliceCheckerTest(t, true, method, slc, c.Explain)
	}
}

func assertFailSliceChecker(t *testing.T, method string, c check.ValueChecker, slc interface{}) {
	t.Helper()
	if c.Pass(slc) {
		failSliceCheckerTest(t, false, method, slc, c.Explain)
	}
}

func failSliceCheckerTest(t *testing.T, expPass bool, method string, slc interface{}, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Slice."+method, explain("Slice value", slc))
}
