package check_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestValueCheckerProvider(t *testing.T) {
	type Person struct {
		Name string
	}

	var (
		v      = Person{"hi"}
		vcopy  = Person{"hi"}
		vempty = Person{}
		badval = Person{"hello"}
		badtyp = struct{ Name string }{"hi"}

		emptyMap   map[int]bool
		emptySlice []float32

		zeros    = []interface{}{0, "", 0i + 0, vempty, emptyMap, emptySlice}
		nonZeros = []interface{}{1, "hi", 0i + 1, v, map[int]bool{}, []float32{}}
	)

	t.Run("Is pass", func(t *testing.T) {
		c := check.Value.Is(vcopy)
		assertPassValueChecker(t, "Is", c, v)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Value.Is(badval)
		assertFailValueChecker(t, "Is", c, v)
		c = check.Value.Is(badtyp)
		assertFailValueChecker(t, "Is", c, v)
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.Value.Not(badval, badtyp)
		assertPassValueChecker(t, "Not", c, v)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.Value.Not(badval, vcopy, badtyp)
		assertFailValueChecker(t, "Not", c, v)
	})

	t.Run("IsZero pass", func(t *testing.T) {
		c := check.Value.IsZero()
		for _, z := range zeros {
			assertPassValueChecker(t, "IsZero", c, z)
		}
	})

	t.Run("IsZero fail", func(t *testing.T) {
		c := check.Value.IsZero()
		for _, nz := range nonZeros {
			assertFailValueChecker(t, "IsZero", c, nz)
		}
	})

	t.Run("NotZero pass", func(t *testing.T) {
		c := check.Value.NotZero()
		for _, nz := range nonZeros {
			assertPassValueChecker(t, "NotZero", c, nz)
		}
	})

	t.Run("NotZero fail", func(t *testing.T) {
		c := check.Value.NotZero()
		for _, z := range zeros {
			assertFailValueChecker(t, "NotZero", c, z)
		}
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
