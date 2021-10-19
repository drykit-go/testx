package testutil_test

import (
	"testing"

	"github.com/drykit-go/testx/internal/testutil"
)

func TestAssertPanic(t *testing.T) {
	t.Run("no panic fails", func(t *testing.T) {
		tt := &testing.T{}
		testutil.AssertPanic(tt, "")
		if !tt.Failed() {
			t.Error("did not fail")
		}
	})

	t.Run("panic with uexpected message fails", func(t *testing.T) {
		tt := &testing.T{}
		func() {
			defer testutil.AssertPanic(tt, "some message")
			panic("bad message")
		}()
		if !tt.Failed() {
			t.Error("did not fail")
		}
	})

	t.Run("panic with expected message passes", func(t *testing.T) {
		tt := &testing.T{}
		func() {
			defer testutil.AssertPanic(tt, "ok")
			panic("ok")
		}()
		if tt.Failed() {
			t.Error("did fail")
		}
	})
}
