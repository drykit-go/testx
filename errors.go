package testx

import (
	"errors"
	"fmt"
)

var (
	ErrTableRunnerConfig    = errors.New("invalid TableConfig")
	ErrTableRunnerFunc      = errors.New("invalid Table func")
	ErrTableRunnerFuncNumIn = fmt.Errorf(
		"%w: it must accept at least 1 parameter",
		ErrTableRunnerFunc,
	)
	ErrTableRunnerFuncNumOut = fmt.Errorf(
		"%w: it must return at least 1 value",
		ErrTableRunnerFunc,
	)
)

func errTableRunnerConfigInPos(funcName string, inPos, numIn int) error {
	return fmt.Errorf(
		"%w: InPos: exp 0 <= n < %d (number of parameters of %s), got %d",
		ErrTableRunnerConfig, numIn, funcName, inPos,
	)
}

func errTableRunnerConfigOutPos(funcName string, outPos, numOut int) error {
	return fmt.Errorf(
		"%w: OutPos: exp 0 <= n < %d (number of values returned by %s), got %d",
		ErrTableRunnerConfig, numOut, funcName, outPos,
	)
}

func errTableRunnerConfigFixedArgs(n int) error {
	return fmt.Errorf("%w: invalid FixedArgs number: %d", ErrTableRunnerConfig, n)
}
