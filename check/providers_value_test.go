package check_test

import (
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

		zeros   = []interface{}{0, "", 0i + 0, vzero, emptyMap, emptySlice}
		nozeros = []interface{}{1, "hi", 0i + 1, vorig, map[int]bool{}, []float32{}}

		isInt = func(got interface{}) bool {
			_, ok := got.(int)
			return ok
		}
		intValue    = 42
		stringValue = "hi"
	)

	t.Run("Is pass", func(t *testing.T) {
		c := check.Value.Is(vsame)
		assertPassValueChecker(t, "Is", c, vorig)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Value.Is(badval)
		assertFailValueChecker(t, "Is", c, vorig)
		c = check.Value.Is(badtyp)
		assertFailValueChecker(t, "Is", c, vorig)
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.Value.Not(badval, badtyp)
		assertPassValueChecker(t, "Not", c, vorig)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.Value.Not(badval, vsame, badtyp)
		assertFailValueChecker(t, "Not", c, vorig)
	})

	t.Run("IsZero pass", func(t *testing.T) {
		c := check.Value.IsZero()
		for _, z := range zeros {
			assertPassValueChecker(t, "IsZero", c, z)
		}
	})

	t.Run("IsZero fail", func(t *testing.T) {
		c := check.Value.IsZero()
		for _, nz := range nozeros {
			assertFailValueChecker(t, "IsZero", c, nz)
		}
	})

	t.Run("NotZero pass", func(t *testing.T) {
		c := check.Value.NotZero()
		for _, nz := range nozeros {
			assertPassValueChecker(t, "NotZero", c, nz)
		}
	})

	t.Run("NotZero fail", func(t *testing.T) {
		c := check.Value.NotZero()
		for _, z := range zeros {
			assertFailValueChecker(t, "NotZero", c, z)
		}
	})

	t.Run("Custom pass", func(t *testing.T) {
		c := check.Value.Custom("", isInt)
		assertPassValueChecker(t, "Custom", c, intValue)
	})

	t.Run("Custom fail", func(t *testing.T) {
		c := check.Value.Custom("", isInt)
		assertFailValueChecker(t, "Custom", c, stringValue)
	})

	t.Run("SameJSON pass", func(t *testing.T) {
		mapequiv := map[string]interface{}{
			"Name": "hi",
		}
		c := check.Value.SameJSON(mapequiv)
		assertPassValueChecker(t, "SameJSON", c, vorig)
	})

	t.Run("SameJSON fail", func(t *testing.T) {
		mapdiff := map[string]interface{}{
			"Name": "bad",
		}
		c := check.Value.SameJSON(mapdiff)
		assertFailValueChecker(t, "SameJSON", c, vorig)
	})
}

// Helpers

func assertPassValueChecker(t *testing.T, method string, c check.ValueChecker, v interface{}) {
	t.Helper()
	if !c.Pass(v) {
		failValueCheckerTest(t, true, method, v, c.Explain)
	}
}

func assertFailValueChecker(t *testing.T, method string, c check.ValueChecker, v interface{}) {
	t.Helper()
	if c.Pass(v) {
		failValueCheckerTest(t, false, method, v, c.Explain)
	}
}

func failValueCheckerTest(t *testing.T, expPass bool, method string, v interface{}, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Value."+method, explain("Value value", v))
}
