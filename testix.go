package testix

import (
	"testing"

	"github.com/drykit-go/testix/check"
)

// TODO: package name
// - testx -> testx.HandlerFunc().ResponseCode()
// - testix -> testix.HandlerFunc().ResponseCode()
// - checkthat -> checkthat.HandlerFunc().HasResponseCode()
// - checkmy -> checkmy.HandlerFunc().MatchesResponseCode()
// - testmy -> testmy.HandlerFunc().ResponseCode()
// - testreq -> testreq.HandlerFunc().ResponseCode()
// - expect
// - xpect

// TODO: test cases helpers

// TODO: response headers checks

// TODO: finish moving test from intenal/handlers_test.go to testix/handlers_test.go

// TODO: decide between:
// - testix.HandlerFunc(h, r).Run(t)
// - testix.HandlerFunc(t, h, r).Run()

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
