package testx

import (
	"errors"
	"fmt"
)

var (
	// ErrTableRunnerConfig is returned when TableRunner is provided
	// a TableConfig that is invalid or incompatible with the tested func.
	ErrTableRunnerConfig = errors.New("invalid TableConfig")
	// ErrTableRunnerFunc is returned when TableRunner is initialized
	// with an incompatible function (most likely it doesn't accept
	// parameters or doesn't return any values).
	ErrTableRunnerFunc = errors.New("invalid Table func")
	// ErrTableRunnerFuncNumIn is returned when TableRunner is initialized
	// with a function that doesn't accept parameters.
	ErrTableRunnerFuncNumIn = fmt.Errorf(
		"%w: it must accept at least 1 parameter",
		ErrTableRunnerFunc,
	)
	// ErrTableRunnerFuncNumOut is returned when TableRunner is initialized
	// with a function that doesn't return any value.
	ErrTableRunnerFuncNumOut = fmt.Errorf(
		"%w: it must return at least 1 value",
		ErrTableRunnerFunc,
	)
)

// errTableRunnerConfigInPos returns an error reporting an invalid value
// for TableConfig.InPos.
func errTableRunnerConfigInPos(funcName string, inPos, numIn int) error {
	return fmt.Errorf(
		"%w: InPos: exp 0 <= n < %d (number of parameters of %s), got %d",
		ErrTableRunnerConfig, numIn, funcName, inPos,
	)
}

// errTableRunnerConfigOutPos returns an error reporting an invalid value
// for TableConfig.OutPos.
func errTableRunnerConfigOutPos(funcName string, outPos, numOut int) error {
	return fmt.Errorf(
		"%w: OutPos: exp 0 <= n < %d (number of values returned by %s), got %d",
		ErrTableRunnerConfig, numOut, funcName, outPos,
	)
}

// errTableRunnerConfigOutPos returns an error reporting an invalid value
// for TableConfig.FixedArgs.
func errTableRunnerConfigFixedArgs(n int) error {
	return fmt.Errorf("%w: invalid FixedArgs number: %d", ErrTableRunnerConfig, n)
}
