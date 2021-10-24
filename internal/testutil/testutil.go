package testutil

import "testing"

// AssertPanicMessage fails t if the calling func does not panic
// with the expected message. It must be called with defer:
//
// 	func TestFuncThatPanics(t *testing.T) {
// 		defer testutil.AssertPanicMessage(t, "oops, I panicked")
// 		FuncThatPanics()
// 	}
func AssertPanicMessage(t *testing.T, expMessage string) {
	t.Helper()
	r := recover()
	assertPanicked(t, r, expMessage, true)
}

// AssertPanic fails t if the calling func does not panic.
// It must be called with defer:
//
// 	func TestFuncThatPanics(t *testing.T) {
// 		defer testutil.AssertPanic(t)
// 		FuncThatPanics()
// 	}
func AssertPanic(t *testing.T) {
	t.Helper()
	r := recover()
	assertPanicked(t, r, "", false)
}

func assertPanicked(t *testing.T, rec interface{}, msg string, checkmsg bool) {
	t.Helper()
	if rec == nil {
		t.Error("expected to panic but did not")
	} else if checkmsg && rec != msg {
		t.Errorf("bad panic message:\nexp %s\ngot %s", msg, rec)
	}
}
