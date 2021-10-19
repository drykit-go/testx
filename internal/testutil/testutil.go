package testutil

import "testing"

func AssertPanic(t *testing.T, expMessage string) {
	t.Helper()
	r := recover()
	if r == nil {
		t.Error("expected to panic but did not")
	} else if r != expMessage {
		t.Errorf("bad panic message:\nexp %s\ngot %s", expMessage, r)
	}
}
