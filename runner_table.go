package testx

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/drykit-go/cond"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/internal/fmtexpl"
	"github.com/drykit-go/testx/internal/reflectutil"
)

var _ TableRunner = (*tableRunner)(nil)

// Case represents a Table test case. It must be provided values for
// Case.In, and Case.Exp or Case.Not or Case.Pass at least.
type Case struct {
	// Lab is the label of the current case to be printed if the current
	// case fails.
	Lab string

	// In is the input value injected in the tested func.
	In interface{}

	// Exp is the value expected to be returned when calling the tested func.
	// If Case.Exp == nil (zero value), no check is added. This is a necessary
	// behavior if one wants to use Case.Pass or Case.Not but not Case.Exp.
	// To specifically check for a nil value, use ExpNil.
	//
	// 	testx.Table(myFunc, nil).Cases([]testx.Case{
	// 		{In: 123, Pass: checkers},    // Exp == nil, no Exp check added besides checkers
	// 		{In: 123},                    // Exp == nil, no check added
	// 		{In: 123, Exp: nil},          // Exp == nil, no check added
	// 		{In: 123, Exp: testx.ExpNil}, // expect nil value
	// 	})
	Exp interface{}

	// Not is a slice of values expected not to be returned by the tested func.
	Not []interface{}

	// Pass is a slice of check.ValueChecker that the return values of the
	// tested func is expected to pass.
	Pass []check.ValueChecker
}

// TableConfig is an object of options allowing to configure a table runner.
// It allows to test functions having multiple input parameters or multiple
// return values.

type TableConfig struct {
	// InPos is the nth parameter in which In value is injected.
	// It is required if the tested func accepts multiple parameters.
	// Default is 0.
	InPos int

	// OutPos is the nth output value that is tested against Case.Exp.
	// It is required if the tested func returns multiple values.
	// Default is 0.
	OutPos int

	// FixedArgs is a slice of arguments to be injected into the tested func.
	// It is required if the tested func accepts multiple parameters.
	//
	// Let nparams the number of parameters of the tested func, len(FixedArgs)
	// must equal nparams or nparams - 1.
	//
	// The following examples produce the same result:
	//
	// 	testx.Table(myFunc).Config(
	// 		InPos: 1
	// 		FixedArgs: []interface{"myArg0", "myArg2"} // len(FixedArgs) == 2
	// 	)
	//
	// 	testx.Table(myFunc).Config(
	// 		InPos: 1
	// 		FixedArgs: []interface{0: "myArg0", 2: "myArg2"} // len(FixedArgs) == 3
	// 	)
	FixedArgs Args
}

type Args []interface{}

func (args Args) replaceAt(pos int, arg interface{}) Args {
	if pos >= len(args) {
		log.Panic("Args.replaceAt(i, v): i is out of range")
	}
	args[pos] = arg
	return args
}

func (args Args) String() string {
	argsStr := []string{}
	for _, arg := range args {
		argsStr = append(argsStr, fmt.Sprintf("%v", arg))
	}
	str := strings.Join(argsStr, ", ")
	return str
}

type tableRunner struct {
	baseRunner

	config TableConfig
	get    func(in interface{}) gottype

	rfunc *reflectutil.Func
	args  Args
}

func (r *tableRunner) Run(t *testing.T) {
	t.Helper()
	r.dryRun()
	r.run(t)
}

func (r *tableRunner) DryRun() TableResulter {
	r.dryRun()
	return tableResults{baseResults: r.baseResults()}
}

func (r *tableRunner) dryRun() {
	cond.PanicOnErr(r.validateConfig())

	args, err := r.makeFixedArgs(r.rfunc, r.config)
	cond.PanicOnErr(err)

	r.setGetFunc(args)
}

func (r *tableRunner) Cases(cases []Case) TableRunner {
	for i, tc := range cases {
		i, tc := i, tc

		get := func() gottype { return r.get(tc.In) }
		getLabel := func() string {
			return fmtexpl.TableCaseLabel(r.rfunc.Name, i, tc.Lab, r.args)
		}

		caseCheck := func(c check.ValueChecker) baseCheck {
			return baseCheck{
				get:      get,
				getLabel: getLabel,
				label:    tc.Lab,
				checker:  c,
			}
		}

		// add Case.Exp check
		if tc.Exp != nil {
			exp := cond.Value(nil, tc.Exp, tc.Exp == ExpNil)
			r.addCheck(caseCheck(check.Value.Is(exp)))
		}

		// add Case.Not checks
		if len(tc.Not) != 0 {
			r.addCheck(caseCheck(check.Value.Not(tc.Not...)))
		}

		// add Case.Pass checks
		r.addChecks(tc.Lab, get, tc.Pass, true)
	}
	return r
}

func (r *tableRunner) Config(cfg TableConfig) TableRunner {
	r.config = cfg
	return r
}

func (r *tableRunner) setRfunc(in interface{}) error {
	rfunc, err := reflectutil.NewFunc(in)
	if err != nil {
		return fmt.Errorf("Table(func): %w", err)
	}
	ftype := rfunc.Value.Type()
	if ftype.NumIn() == 0 {
		return fmt.Errorf("Table(%s): %w", rfunc.Name, ErrTableRunnerFuncNumIn)
	}
	if ftype.NumOut() == 0 {
		return fmt.Errorf("Table(%s): %w", rfunc.Name, ErrTableRunnerFuncNumOut)
	}
	r.rfunc = rfunc
	return nil
}

func (r *tableRunner) setGetFunc(args Args) {
	r.get = func(in interface{}) gottype {
		pin, pout := r.config.InPos, r.config.OutPos
		r.args = args.replaceAt(pin, in)
		return r.rfunc.Call(r.args)[pout]
	}
}

func (r *tableRunner) validateConfig() error {
	validPos := func(pos, max int) bool { return pos >= 0 && pos < max }
	ftyp := r.rfunc.Value.Type()
	if pin, nin := r.config.InPos, ftyp.NumIn(); !validPos(pin, nin) {
		return errTableRunnerConfigInPos(r.rfunc.Name, pin, nin)
	}
	if pout, nout := r.config.OutPos, ftyp.NumOut(); !validPos(pout, nout) {
		return errTableRunnerConfigOutPos(r.rfunc.Name, pout, nout)
	}
	return nil
}

func (r *tableRunner) makeFixedArgs(rfunc *reflectutil.Func, cfg TableConfig) (Args, error) {
	nparams := rfunc.Value.Type().NumIn()
	nargs := len(cfg.FixedArgs)

	fillskip := func(at int) Args {
		args := make(Args, nparams)
		delta := 0
		for i := 0; i < nparams; i++ {
			if i == at {
				delta++
				continue
			}
			args[i] = cfg.FixedArgs[i-delta]
		}
		return args
	}

	switch d := nparams - nargs; d {
	case 0:
		return cfg.FixedArgs, nil
	case 1:
		return fillskip(cfg.InPos), nil
	default:
		return nil, errTableRunnerConfigFixedArgs(d)
	}
}

func newTableRunner(testedFunc interface{}) TableRunner {
	r := &tableRunner{}
	cond.PanicOnErr(r.setRfunc(testedFunc))
	return r
}

/*
	Results
*/

type tableResults struct {
	baseResults
}

func (res tableResults) PassedAt(i int) bool {
	if i >= len(res.checks) {
		panic(fmt.Sprintf("TableResults: index %d is out of range", i))
	}
	return res.checks[i].Passed
}

func (res tableResults) FailedAt(i int) bool {
	return !res.PassedAt(i)
}

func (res tableResults) PassedLabel(label string) bool {
	for _, c := range res.checks {
		if c.label == label {
			return c.Passed
		}
	}
	panic(fmt.Sprintf("TableResults: no test case with label %s", label))
}

func (res tableResults) FailedLabel(label string) bool {
	return !res.PassedLabel(label)
}

/*
	ExpNil
*/

type expNiler interface{ expNil() }

type expNilerImpl struct{}

func (expNilerImpl) expNil() {}

// ExpNil is a value indicating that nil is an expected value.
// It is meant to be used as a Case.Exp value in a TableRunner
// test.
var ExpNil expNiler = expNilerImpl{}
