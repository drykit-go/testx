package testx_test

import (
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

func TestValueRunner(t *testing.T) {
	testx.Value(42).
		MustBe(42).
		MustNotBe(42, 42).
		MustPass(check.Int.InRange(43, 44)).
		Run(t)
}
