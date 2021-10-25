package check_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestBoolCheckerProvider(t *testing.T) {
	const b = true

	t.Run("Is pass", func(t *testing.T) {
		c := check.Bool.Is(b)
		assertPassBoolChecker(t, "Is", c, b)
		c = check.Bool.Is(!b)
		assertPassBoolChecker(t, "Is", c, !b)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Bool.Is(!b)
		assertFailBoolChecker(t, "Is", c, b, makeExpl("false", "true"))
		c = check.Bool.Is(b)
		assertFailBoolChecker(t, "Is", c, !b, makeExpl("true", "false"))
	})
}

// Helpers

func assertPassBoolChecker(t *testing.T, method string, c check.BoolChecker, b bool) {
	t.Helper()
	if !c.Pass(b) {
		failBoolCheckerTest(t, true, method, b, c.Explain)
	}
}

func assertFailBoolChecker(t *testing.T, method string, c check.BoolChecker, b bool, expexpl string) {
	t.Helper()
	if c.Pass(b) {
		failBoolCheckerTest(t, false, method, b, c.Explain)
	}
	if expexpl != "" {
		assertGoodExplain(t, c, b, expexpl)
	}
}

func failBoolCheckerTest(t *testing.T, expPass bool, method string, b bool, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Bool."+method, explain("Bool value", b))
}
