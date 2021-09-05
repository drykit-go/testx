package check_test

import "testing"

func failCheckerTest(t *testing.T, expPass bool, name, expl string) {
	t.Helper()
	passOrFail := "PASS"
	if !expPass {
		passOrFail = "FAIL"
	}
	t.Errorf("\n❌ exp %s to %s, explain:\n%s", name, passOrFail, expl)
}
