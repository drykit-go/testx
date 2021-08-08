package testix

import (
	"testing"

	"github.com/drykit-go/testix/check"
)

type testCheck struct {
	label string
	get   func() gotType // TODO: func() (gotType, error)
	check check.UntypedChecker
}

type gotType interface{}

type getFunc func() gotType

// func failVal(t *testing.T, label string, exp, got interface{}) {
// 	t.Errorf("expected %s %v, got %v", label, exp, got)
// }

// func failErr(t *testing.T, label string, exp interface{}, got error) {
// 	t.Errorf("expected %s %v, got error %s", label, exp, got)
// }

func fail(t *testing.T, msg string) {
	t.Error(msg)
}
