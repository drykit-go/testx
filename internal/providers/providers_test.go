package providers_test

import (
	"fmt"
	"testing"

	check "github.com/drykit-go/testx/internal/checktypes"
)

func assertPassChecker[T any](
	t *testing.T,
	methodName string,
	c check.Checker[T],
	in T,
) {
	t.Helper()
	if !c.Pass(in) {
		failCheckerTest(t, true, methodName, c.Explain("value", in))
	}
}

func assertFailChecker[T any](
	t *testing.T,
	methodName string,
	c check.Checker[T],
	in T,
	expexpl string,
) {
	t.Helper()
	if c.Pass(in) {
		failCheckerTest(t, false, methodName, c.Explain("value", in))
	}
	assertGoodExplain(t, c, in, expexpl)
}

func assertGoodExplain[T any](
	t *testing.T,
	// FIXME: check.Explainer2[T] not working:
	// 	type check.Checker[bool] of c does not match check.Explainer2[T] (cannot infer T)compilerErrorCode(135)
	c check.Checker[T],
	in T,
	expexpl string,
) {
	t.Helper()
	gotexpl := c.Explain("label", in)
	expexpl = "label:\n" + expexpl
	if gotexpl != expexpl {
		t.Errorf("bad Explain output:\nexp \"%s\"\ngot \"%s\"", expexpl, gotexpl)
	}
}

func failCheckerTest(t *testing.T, expPass bool, name, expl string) {
	t.Helper()
	passOrFail := "PASS"
	if !expPass {
		passOrFail = "FAIL"
	}
	t.Errorf("\n❌ exp %s to %s, explain:\n%s", name, passOrFail, expl)
}

func makeExpl(expstr, gotstr string) string {
	return fmt.Sprintf("exp %s\ngot %s", expstr, gotstr)
}
