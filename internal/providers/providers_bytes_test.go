package providers_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/internal/providers"
)

func TestBytesCheckerProvider(t *testing.T) {
	checkBytes := providers.BytesCheckerProvider{}

	var (
		b      = []byte(`{"id":42,"name":"Marcel Patulacci"}`)
		sub    = []byte(`"id":42`)
		diff   = []byte(`{"id":43,"name":"Robert Robichet"}`)
		eqJSON = []byte("{\n\"id\":   42,\n\n\n  \"name\":\"Marcel Patulacci\" } ")
	)

	mapof := func(b []byte) (m map[string]any) {
		json.Unmarshal(b, &m) //nolint:errcheck
		return
	}

	t.Run("Is pass", func(t *testing.T) {
		c := checkBytes.Is(b)
		assertPassChecker(t, "Bytes.Is", c, b)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := checkBytes.Is(diff)
		assertFailChecker(t, "Bytes.Is", c, b, makeExpl(
			fmt.Sprint(diff),
			fmt.Sprint(b),
		))
	})

	t.Run("Not pass", func(t *testing.T) {
		c := checkBytes.Not(diff, eqJSON)
		assertPassChecker(t, "Bytes.Not", c, b)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := checkBytes.Not(diff, eqJSON, b)
		assertFailChecker(t, "Bytes.Not", c, b, makeExpl(
			fmt.Sprintf("not %v", b),
			fmt.Sprint(b),
		))
	})

	t.Run("Len pass", func(t *testing.T) {
		c := checkBytes.Len(check.Int.Is(len(b)))
		assertPassChecker(t, "Bytes.Len", c, b)
	})

	t.Run("Len fail", func(t *testing.T) {
		gotlen := len(b)
		explen := gotlen + 1
		c := checkBytes.Len(check.Int.Is(explen))
		assertFailChecker(t, "Bytes.Len", c, b, makeExpl(
			"length to pass Checker[int]",
			"explanation: length:\n"+makeExpl(
				fmt.Sprint(explen),
				fmt.Sprint(gotlen),
			),
		))
	})

	t.Run("SameJSON pass", func(t *testing.T) {
		c := checkBytes.SameJSON(eqJSON)
		assertPassChecker(t, "Bytes.SameJSON", c, b)
		c = checkBytes.SameJSON(b)
		assertPassChecker(t, "Bytes.SameJSON", c, b)
	})

	t.Run("SameJSON fail", func(t *testing.T) {
		c := checkBytes.SameJSON(diff)
		assertFailChecker(t, "Bytes.SameJSON", c, b, makeExpl(
			fmt.Sprintf("json data: %v", mapof(diff)),
			fmt.Sprintf("json data: %v", mapof(b)),
		))
	})

	t.Run("Contains pass", func(t *testing.T) {
		c := checkBytes.Contains(sub)
		assertPassChecker(t, "Bytes.Contains", c, b)
		c = checkBytes.Contains(b)
		assertPassChecker(t, "Bytes.Contains", c, b)
	})

	t.Run("Contains fail", func(t *testing.T) {
		c := checkBytes.Contains(diff)
		assertFailChecker(t, "Bytes.Contains", c, b,
			makeExpl(
				fmt.Sprintf("to contain subslice %v", diff),
				fmt.Sprint(b),
			),
		)

		c = checkBytes.Contains(eqJSON)
		assertFailChecker(t, "Bytes.Contains", c, b,
			makeExpl(
				fmt.Sprintf("to contain subslice %v", eqJSON),
				fmt.Sprint(b),
			),
		)
	})

	t.Run("NotContains pass", func(t *testing.T) {
		c := checkBytes.NotContains(diff)
		assertPassChecker(t, "Bytes.NotContains", c, b)
		c = checkBytes.NotContains(eqJSON)
		assertPassChecker(t, "Bytes.NotContains", c, b)
	})

	t.Run("NotContains fail", func(t *testing.T) {
		c := checkBytes.NotContains(sub)
		assertFailChecker(t, "Bytes.NotContains", c, b, makeExpl(
			fmt.Sprintf("not to contain subslice %v", sub),
			fmt.Sprint(b),
		))

		c = checkBytes.NotContains(b)
		assertFailChecker(t, "Bytes.NotContains", c, b, makeExpl(
			fmt.Sprintf("not to contain subslice %v", b),
			fmt.Sprint(b),
		))
	})

	t.Run("AsMap pass", func(t *testing.T) {
		c := checkBytes.AsMap(check.Map[string, any]().HasKeys("id"))
		assertPassChecker(t, "Bytes.AsMap", c, b)
		assertPassChecker(t, "Bytes.AsMap", c, eqJSON)
	})

	t.Run("AsMap fail", func(t *testing.T) {
		c := checkBytes.AsMap(check.Map[string, any]().HasKeys("id", "nomatch"))
		assertFailChecker(t, "Bytes.AsMap", c, b, makeExpl(
			"to pass MapChecker",
			"explanation: json map:\n"+makeExpl(
				"to have keys [nomatch]",
				fmt.Sprint(mapof(b)),
			),
		))

		c = checkBytes.AsMap(check.Map[string, any]().HasKeys("id"))
		assertFailChecker(t, "Bytes.AsMap", c, sub, makeExpl(
			"to pass MapChecker",
			"error: json: cannot unmarshal string into Go value of type map[string]interface {}",
		))
	})

	t.Run("AsString pass", func(t *testing.T) {
		c := checkBytes.AsString(check.String.Is(string(b)))
		assertPassChecker(t, "Bytes.AsString", c, b)
	})

	t.Run("AsString fail", func(t *testing.T) {
		c := checkBytes.AsString(check.String.Is(string(diff)))
		assertFailChecker(t, "Bytes.AsString", c, b, makeExpl(
			"to pass Checker[string]",
			"explanation: converted bytes:\n"+makeExpl(
				string(diff),
				string(b),
			),
		))
	})
}
