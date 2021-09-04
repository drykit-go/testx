package check_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
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
		assertFailMapChecker(t, "Len", c, m)
	})

	t.Run("SameJSON pass", func(t *testing.T) {
		c := check.Map.SameJSON(struct {
			Name    string   `json:"name"`
			Age     int      `json:"age"`
			Friends []string `json:"friends"`
		}{
			Name:    "Marcel Patulacci",
			Age:     42,
			Friends: []string{"Robert Robichet", "Jean-Pierre Avidol"},
		})
		assertPassMapChecker(t, "SameJSON", c, m)
	})

	t.Run("SameJSON fail", func(t *testing.T) {
		c := check.Map.SameJSON(struct {
			Name    string   `json:"name"`
			Age     int      `json:"age"`
			Friends []string `json:"friends"`
		}{
			Name:    "Marcel Patulacci",
			Age:     42,
			Friends: []string{"Robert Robichet"},
		})
		assertFailMapChecker(t, "SameJSON", c, m)
	})

	t.Run("HasKeys pass", func(t *testing.T) {
		c := check.Map.HasKeys("name", "friends")
		assertPassMapChecker(t, "HasKeys", c, m)
	})

	t.Run("HasKeys fail", func(t *testing.T) {
		c := check.Map.HasKeys("name", "hello")
		assertFailMapChecker(t, "HasKeys", c, m)
	})

	t.Run("HasNotKeys pass", func(t *testing.T) {
		c := check.Map.HasNotKeys("hello", 42)
		assertPassMapChecker(t, "HasNotKeys", c, m)
	})

	t.Run("HasNotKeys fail", func(t *testing.T) {
		c := check.Map.HasNotKeys("name", "hello")
		assertFailMapChecker(t, "HasNotKeys", c, m)
	})

	t.Run("HasValues pass", func(t *testing.T) {
		c := check.Map.HasValues(42, []string{"Robert Robichet", "Jean-Pierre Avidol"})
		assertPassMapChecker(t, "HasValues", c, m)
	})

	t.Run("HasValues fail", func(t *testing.T) {
		c := check.Map.HasValues(42, "hello")
		assertFailMapChecker(t, "HasValues", c, m)
	})

	t.Run("HasNotValues pass", func(t *testing.T) {
		c := check.Map.HasNotValues("hello", -1)
		assertPassMapChecker(t, "HasNotValues", c, m)
	})

	t.Run("HasNotValues fail", func(t *testing.T) {
		c := check.Map.HasNotValues(-1, []string{"Robert Robichet", "Jean-Pierre Avidol"})
		assertFailMapChecker(t, "HasNotValues", c, m)
	})

	t.Run("CheckValues pass", func(t *testing.T) {
		c := check.Map.CheckValues(
			checkconv.FromInt(check.Int.InRange(41, 43)),
			[]interface{}{"age"},
		)
		assertPassMapChecker(t, "CheckValues", c, m)
	})

	t.Run("CheckValues fail", func(t *testing.T) {
		c := check.Map.CheckValues(
			checkconv.FromInt(check.Int.OutRange(41, 43)),
			[]interface{}{"age"},
		)
		assertFailMapChecker(t, "CheckValues", c, m)
	})
}

// Helpers

func assertPassMapChecker(t *testing.T, method string, c check.ValueChecker, gotm interface{}) {
	t.Helper()
	if !c.Pass(gotm) {
		failMapCheckerTest(t, true, method, gotm, c.Explain)
	}
}

func assertFailMapChecker(t *testing.T, method string, c check.ValueChecker, gotm interface{}) {
	t.Helper()
	if c.Pass(gotm) {
		failMapCheckerTest(t, false, method, gotm, c.Explain)
	}
}

func failMapCheckerTest(t *testing.T, expPass bool, method string, gotm interface{}, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Map."+method, explain("Map value", gotm))
}
