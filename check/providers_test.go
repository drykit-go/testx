package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
)

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

func assertGoodExplain(
	t *testing.T,
	c check.Explainer,
	gotval interface{},
	expexpl string,
) {
	t.Helper()
	gotexpl := c.Explain("label", gotval)
	expexpl = "label:\n" + expexpl
	if gotexpl != expexpl {
		t.Errorf("bad Explain output:\nexp \"%s\"\ngot \"%s\"", expexpl, gotexpl)
	}
}
