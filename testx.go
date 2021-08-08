package testx

import (
	"testing"

	"github.com/drykit-go/testx/check"
)

type testCheck struct {
	label string
	get   func() gotType // TODO: func() (gotType, error)
	check check.UntypedChecker
}

type gotType interface{}

type getFunc func() gotType

func fail(t *testing.T, msg string) {
	t.Error(msg)
}
