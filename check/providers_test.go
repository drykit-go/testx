package check_test

import "testing"

func failCheckerTest(t *testing.T, expPass bool, name, expl string) {
	t.Helper()
	passOrFail := "pass"
	if !expPass {
		passOrFail = "fail"
	}
	t.Errorf("should %s %s, explain:\n%s", passOrFail, name, expl)
}
