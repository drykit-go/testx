package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestMapCheckerProvider(t *testing.T) {
	m := map[string]any{
		"name":    "Marcel Patulacci",
		"age":     42,
		"friends": []string{"Robert Robichet", "Jean-Pierre Avidol"},
	}

	t.Run("Len pass", func(t *testing.T) {
		c := check.Map[string, any]().Len(check.Int.Is(3))
		assertPassChecker(t, "Map.Len", c, m)
	})

	t.Run("Len fail", func(t *testing.T) {
		c := check.Map[string, any]().Len(check.Int.Not(3))
		assertFailChecker(t, "Map.Len", c, m, makeExpl(
			"length to pass Checker[int]",
			"explanation: length:\n"+makeExpl("not 3", "3"),
		))
	})

	t.Run("HasKeys pass", func(t *testing.T) {
		c := check.Map[string, any]().HasKeys("name", "friends")
		assertPassChecker(t, "Map.HasKeys", c, m)
	})

	t.Run("HasKeys fail", func(t *testing.T) {
		c := check.Map[string, any]().HasKeys("name", "hello", "bad")
		assertFailChecker(t, "Map.HasKeys", c, m, makeExpl(
			"to have keys [hello, bad]",
			fmt.Sprint(m),
		))
	})

	t.Run("HasNotKeys pass", func(t *testing.T) {
		c := check.Map[string, any]().HasNotKeys("hello", "42")
		assertPassChecker(t, "Map.HasNotKeys", c, m)
	})

	t.Run("HasNotKeys fail", func(t *testing.T) {
		c := check.Map[string, any]().HasNotKeys("name", "hello", "age")
		assertFailChecker(t, "Map.HasNotKeys", c, m, makeExpl(
			"not to have keys [name, age]",
			fmt.Sprint(m),
		))
	})

	t.Run("HasValues pass", func(t *testing.T) {
		c := check.Map[string, any]().HasValues(42, []string{"Robert Robichet", "Jean-Pierre Avidol"})
		assertPassChecker(t, "Map.HasValues", c, m)
	})

	t.Run("HasValues fail", func(t *testing.T) {
		c := check.Map[string, any]().HasValues(42, "hello", true)
		assertFailChecker(t, "Map.HasValues", c, m, makeExpl(
			"to have values [hello, true]",
			fmt.Sprint(m),
		))
	})

	t.Run("HasNotValues pass", func(t *testing.T) {
		c := check.Map[string, any]().HasNotValues("hello", -1)
		assertPassChecker(t, "Map.HasNotValues", c, m)
	})

	t.Run("HasNotValues fail", func(t *testing.T) {
		c := check.Map[string, any]().HasNotValues(42, "hi", []string{"Robert Robichet", "Jean-Pierre Avidol"})
		assertFailChecker(t, "Map.HasNotValues", c, m, makeExpl(
			"not to have values [42, [Robert Robichet Jean-Pierre Avidol]]",
			fmt.Sprint(m),
		))
	})

	t.Run("CheckValues pass", func(t *testing.T) {
		// keys subset
		c := check.Map[string, any]().CheckValues(
			check.Wrap(check.Int.InRange(41, 43)),
			"age",
		)
		assertPassChecker(t, "Map.CheckValues", c, m)

		// all keys
		c = check.Map[string, any]().CheckValues(check.Value[any]().Not(0))
		assertPassChecker(t, "Map.CheckValues", c, m)
	})

	t.Run("CheckValues fail", func(t *testing.T) {
		// keys subset
		c := check.Map[string, any]().CheckValues(
			check.Wrap(check.Int.OutRange(41, 43)),
			"age", "badkey",
		)
		assertFailChecker(t, "Map.CheckValues", c, m, makeExpl(
			"values for keys [age badkey] to pass Checker[V]",
			"explanation: values:\n"+makeExpl(
				"not in range [41:43]",
				"[age:42, badkey:<nil>]",
			),
		))

		// all keys
		c = check.Map[string, any]().CheckValues(check.Value[any]().Is("Marcel Patulacci"))
		assertFailChecker(t, "Map.CheckValues", c, m, makeExpl(
			"values for all keys to pass Checker[V]",
			"explanation: values:\n"+makeExpl(
				"Marcel Patulacci",
				"[age:42, friends:[Robert Robichet Jean-Pierre Avidol]]",
			),
		))
	})
}
