package testutil

import "testing"

// AssertPanic fails t if the calling func does not panic
// with the expected message. It must be called with defer:
//
// 	func TestFuncThatPanics(t *testing.T) {
// 		defer testutil.AssertPanic(t, "oops, I panicked")
// 		FuncThatPanics()
// 	}
func AssertPanic(t *testing.T, expMessage string) {
	t.Helper()
	r := recover()
	if r == nil {
		t.Error("expected to panic but did not")
	} else if r != expMessage {
		t.Errorf("bad panic message:\nexp %s\ngot %s", expMessage, r)
	}
}
