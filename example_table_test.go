package testx_test

import (
	"errors"
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

func ExampleTable_monadic() {
	t := &testing.T{} // ignore: emulating a testing context

	// double is the func to be tested.
	double := func(x float64) float64 { return 2 * x }

	// func double has 1 parameter and 1 return value,
	// hence no config is needed
	testx.Table[float64, float64](double).Cases([]testx.Case[float64, float64]{
		{In: 0.0, Exp: 0.0},
		{In: -2.0, Pass: []check.Checker[float64]{check.Float64.InRange(-5, -3)}},
	}).Run(t)
}

func ExampleTable_dyadic() {
	t := &testing.T{} // ignore: emulating a testing context

	// divide is the func to be tested.
	// It returns x/y or a non-nil error if y == 0
	divide := func(x, y float64) (float64, error) {
		if y == 0 {
			return 0, errors.New("division by 0")
		}
		return x / y, nil
	}

	// func divide has several parameters and return values,
	// so we specify a config to determinate:
	// - at which param position Case.In is injected
	// - the values used for the other arguments (fixed for all cases)
	// - which return value we want to compare Case.Exp with
	//
	// In this example, we check the error value of divide (return value
	// at position 1).
	// We inject Case.In at position 1 (param y) and use a fixed value
	// of 42.0 at position 0 (param x) for all cases.
	testx.Table[float64, error](divide).Config(testx.TableConfig{
		// Positions start at 0
		InPos:     1,              // Case.In injected in param position 1 (y)
		OutPos:    1,              // Case.Exp compared to return value position 1 (error value)
		FixedArgs: []any{0: 42.0}, // param 0 (x) set to 42.0 for all cases
	}).Cases([]testx.Case[float64, error]{
		// FIXME: ExpNil typing issues
		// {In: 1.0, Exp: testx.ExpNil},                // divide(42.0, 1.0) -> (_, nil)
		{In: 0.0, Exp: errors.New("division by 0")}, // divide(42.0, 0.0) -> (_, err)
	}).Run(t)
}
