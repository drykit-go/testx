package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

func TestMapCheckerProvider(t *testing.T) {
	m := map[string]interface{}{
		"name":    "Marcel Patulacci",
		"age":     42,
		"friends": []string{"Robert Robichet", "Jean-Pierre Avidol"},
	}

	t.Run("Len pass", func(t *testing.T) {
		c := check.Map.Len(check.Int.Is(3))
		assertPassMapChecker(t, "Len", c, m)
	})

	t.Run("Len fail", func(t *testing.T) {
		c := check.Map.Len(check.Int.Not(3))
		assertFailMapChecker(t, "Len", c, m, makeExpl(
			"length to pass IntChecker",
			"explanation: length:\nexp not 3\ngot 3",
		))
	})

	t.Run("HasKeys pass", func(t *testing.T) {
		c := check.Map.HasKeys("name", "friends")
		assertPassMapChecker(t, "HasKeys", c, m)
	})

	t.Run("HasKeys fail", func(t *testing.T) {
		c := check.Map.HasKeys("name", "hello", "bad")
		assertFailMapChecker(t, "HasKeys", c, m, makeExpl(
			"to have keys [hello, bad]",
			fmt.Sprint(m),
		))
	})

	t.Run("HasNotKeys pass", func(t *testing.T) {
		c := check.Map.HasNotKeys("hello", 42)
		assertPassMapChecker(t, "HasNotKeys", c, m)
	})

	t.Run("HasNotKeys fail", func(t *testing.T) {
		c := check.Map.HasNotKeys("name", "hello", "age")
		assertFailMapChecker(t, "HasNotKeys", c, m, makeExpl(
			"not to have keys [name, age]",
			fmt.Sprint(m),
		))
	})

	t.Run("HasValues pass", func(t *testing.T) {
		c := check.Map.HasValues(42, []string{"Robert Robichet", "Jean-Pierre Avidol"})
		assertPassMapChecker(t, "HasValues", c, m)
	})

	t.Run("HasValues fail", func(t *testing.T) {
		c := check.Map.HasValues(42, "hello", true)
		assertFailMapChecker(t, "HasValues", c, m, makeExpl(
			"to have values [hello, true]",
			fmt.Sprint(m),
		))
	})

	t.Run("HasNotValues pass", func(t *testing.T) {
		c := check.Map.HasNotValues("hello", -1)
		assertPassMapChecker(t, "HasNotValues", c, m)
	})

	t.Run("HasNotValues fail", func(t *testing.T) {
		c := check.Map.HasNotValues(42, "hi", []string{"Robert Robichet", "Jean-Pierre Avidol"})
		assertFailMapChecker(t, "HasNotValues", c, m, makeExpl(
			"not to have values [42, [Robert Robichet Jean-Pierre Avidol]]",
			fmt.Sprint(m),
		))
	})

	t.Run("CheckValues pass", func(t *testing.T) {
		// keys subset
		c := check.Map.CheckValues(
			checkconv.FromInt(check.Int.InRange(41, 43)),
			"age",
		)
		assertPassMapChecker(t, "CheckValues", c, m)

		// all keys
		c = check.Map.CheckValues(check.Value.Not(0))
		assertPassMapChecker(t, "CheckValues", c, m)
	})

	t.Run("CheckValues fail", func(t *testing.T) {
		// keys subset
		c := check.Map.CheckValues(
			checkconv.FromInt(check.Int.OutRange(41, 43)),
			"age", "badkey",
		)
		assertFailMapChecker(t, "CheckValues", c, m, makeExpl(
			"values for keys [age badkey] to pass ValueChecker",
			"explanation: values:\nexp not in range [41:43]\ngot [age:42, badkey:<nil>]",
		))

		// all keys
		c = check.Map.CheckValues(check.Value.Is("Marcel Patulacci"))
		assertFailMapChecker(t, "CheckValues", c, m, makeExpl(
			"values for all keys to pass ValueChecker",
			"explanation: values:\nexp Marcel Patulacci\ngot "+
				"[age:42, friends:[Robert Robichet Jean-Pierre Avidol]]",
		))
	})
}

// Helpers

func assertPassMapChecker(t *testing.T, method string, c check.ValueChecker, gotm interface{}) {
	t.Helper()
	if !c.Pass(gotm) {
		failMapCheckerTest(t, true, method, gotm, c.Explain)
	}
}

func assertFailMapChecker(t *testing.T, method string, c check.ValueChecker, gotm interface{}, expexpl string) {
	t.Helper()
	if c.Pass(gotm) {
		failMapCheckerTest(t, false, method, gotm, c.Explain)
	}
	assertGoodExplain(t, c, gotm, expexpl)
}

func failMapCheckerTest(t *testing.T, expPass bool, method string, gotm interface{}, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Map."+method, explain("Map value", gotm))
}
