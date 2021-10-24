package check_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestBytesCheckerProvider(t *testing.T) {
	var (
		b      = []byte(`{"id":42,"name":"Marcel Patulacci"}`)
		sub    = []byte(`"id":42`)
		diff   = []byte(`{"id":43,"name":"Robert Robichet"}`)
		eqJSON = []byte("{\n\"id\":   42,\n\n\n  \"name\":\"Marcel Patulacci\" } ")
	)

	t.Run("Is pass", func(t *testing.T) {
		c := check.Bytes.Is(b)
		assertPassBytesChecker(t, "Is", c, b)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Bytes.Is(diff)
		assertFailBytesChecker(t, "Is", c, b)
	})

	t.Run("Not pass", func(t *testing.T) {
		c := check.Bytes.Not(diff, eqJSON)
		assertPassBytesChecker(t, "Not", c, b)
	})

	t.Run("Not fail", func(t *testing.T) {
		c := check.Bytes.Not(diff, eqJSON, b)
		assertFailBytesChecker(t, "Not", c, b)
	})

	t.Run("Len pass", func(t *testing.T) {
		c := check.Bytes.Len(check.Int.Is(len(b)))
		assertPassBytesChecker(t, "Len", c, b)
	})

	t.Run("Len fail", func(t *testing.T) {
		c := check.Bytes.Len(check.Int.Is(len(b) + 1))
		assertFailBytesChecker(t, "Len", c, b)
	})

	t.Run("SameJSON pass", func(t *testing.T) {
		c := check.Bytes.SameJSON(eqJSON)
		assertPassBytesChecker(t, "SameJSON", c, b)
		c = check.Bytes.SameJSON(b)
		assertPassBytesChecker(t, "SameJSON", c, b)
	})

	t.Run("SameJSON fail", func(t *testing.T) {
		c := check.Bytes.SameJSON(diff)
		assertFailBytesChecker(t, "SameJSON", c, b)
	})

	t.Run("Contains pass", func(t *testing.T) {
		c := check.Bytes.Contains(sub)
		assertPassBytesChecker(t, "Contains", c, b)
		c = check.Bytes.Contains(b)
		assertPassBytesChecker(t, "Contains", c, b)
	})

	t.Run("Contains fail", func(t *testing.T) {
		c := check.Bytes.Contains(diff)
		assertFailBytesChecker(t, "Contains", c, b)
		c = check.Bytes.Contains(eqJSON)
		assertFailBytesChecker(t, "Contains", c, b)
	})

	t.Run("NotContains pass", func(t *testing.T) {
		c := check.Bytes.NotContains(diff)
		assertPassBytesChecker(t, "NotContains", c, b)
		c = check.Bytes.NotContains(eqJSON)
		assertPassBytesChecker(t, "NotContains", c, b)
	})

	t.Run("NotContains fail", func(t *testing.T) {
		c := check.Bytes.NotContains(sub)
		assertFailBytesChecker(t, "NotContains", c, b)
		c = check.Bytes.NotContains(b)
		assertFailBytesChecker(t, "NotContains", c, b)
	})

	t.Run("AsMap pass", func(t *testing.T) {
		c := check.Bytes.AsMap(check.Map.HasKeys("id"))
		assertPassBytesChecker(t, "AsMap", c, b)
		assertPassBytesChecker(t, "AsMap", c, eqJSON)
	})

	t.Run("AsMap fail", func(t *testing.T) {
		c := check.Bytes.AsMap(check.Map.HasKeys("id", "nomatch"))
		assertFailBytesChecker(t, "AsMap", c, b)
		c = check.Bytes.AsMap(check.Map.HasKeys("id"))
		assertFailBytesChecker(t, "AsMap", c, sub)
	})

	t.Run("AsString pass", func(t *testing.T) {
		c := check.Bytes.AsString(check.String.Is(string(b)))
		assertPassBytesChecker(t, "AsString", c, b)
	})

	t.Run("AsString fail", func(t *testing.T) {
		c := check.Bytes.AsString(check.String.Is(string(diff)))
		assertFailBytesChecker(t, "AsString", c, b)
	})
}

// Helpers

func assertPassBytesChecker(t *testing.T, method string, c check.BytesChecker, b []byte) {
	t.Helper()
	if !c.Pass(b) {
		failBytesCheckerTest(t, true, method, b, c.Explain)
	}
}

func assertFailBytesChecker(t *testing.T, method string, c check.BytesChecker, b []byte) {
	t.Helper()
	if c.Pass(b) {
		failBytesCheckerTest(t, false, method, b, c.Explain)
	}
}

func failBytesCheckerTest(t *testing.T, expPass bool, method string, b []byte, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Bytes."+method, explain("Bytes value", b))
}
