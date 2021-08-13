package testx_test

import (
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

func TestValueRunner(t *testing.T) {
	testx.Value(42).
		MustBe(42).
		MustNotBe(3, "hello").
		MustPass(check.Int.InRange(41, 43)).
		Run(t)
}
