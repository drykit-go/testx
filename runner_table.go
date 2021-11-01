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
	In any

	// Exp is the value expected to be returned when calling the tested func.
	// If Case.Exp == nil (zero value), no check is added. This is a necessary
	// behavior if one wants to use Case.Pass or Case.Not but not Case.Exp.
	// To specifically check for a nil value, use ExpNil.
	//
	// 	testx.Table(myFunc, nil).Cases([]testx.Case{
	// 		{In: 123, Pass: checkers},    // Exp == nil, no Exp check added
	// 		{In: 123},                    // Exp == nil, no Exp check added
	// 		{In: 123, Exp: nil},          // Exp == nil, no Exp check added
	// 		{In: 123, Exp: testx.ExpNil}, // Exp == ExpNil, expect nil value
	// 	})
	Exp any

	// Not is a slice of values expected not to be returned by the tested func.
	Not []any

	// Pass is a slice of checkers that the return value of the
	// tested func is expected to pass.
	Pass []check.Checker[any]
}

// TableConfig is configuration object for TableRunner.
// It allows to test functions having multiple parameters or multiple
// return values.
// Its zero value is a valid config for functions of 1 parameter
// and 1 return value, so it can be omitted in that case.
type TableConfig struct {
	// InPos is the nth parameter in which Case.In value is injected,
	// starting at 0.
	// It is required if the tested func accepts multiple parameters.
	// Default is 0.
	InPos int

	// OutPos is the nth return value that is tested against Case.Exp,
	// starting at 0.
	// It is required if the tested func returns multiple values.
	// Default is 0.
	OutPos int

	// FixedArgs is a slice of arguments to be injected into the tested func.
	// Its values are fixed for all cases.
	// It is required if the tested func accepts multiple parameters.
	//
	// Let nparams the number of parameters of the tested func, len(FixedArgs)
	// must equal nparams or nparams - 1.
	//
	// The following configurations produce the same result:
	//
	// 	testx.Table(myFunc).Config(testx.TableConfig{
	// 		InPos: 1
	// 		FixedArgs: []interface{"myArg0", "myArg2"} // len(FixedArgs) == 2
	// 	})
	//
	// 	testx.Table(myFunc).Config(testx.TableConfig{
	// 		InPos: 1
	// 		FixedArgs: []interface{0: "myArg0", 2: "myArg2"} // len(FixedArgs) == 3
	// 	})
	FixedArgs Args
}

// Args is an alias to []any.
type Args []any

func (args Args) replaceAt(pos int, arg any) Args {
	if pos >= len(args) {
		log.Panic("Args.replaceAt(i, v): i is out of range")
	}
	args[pos] = arg
	return args
}

func (args Args) String() string {
	var b strings.Builder
	for i, arg := range args {
		b.WriteString(fmt.Sprint(arg))
		if i != len(args)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

type tableRunner struct {
	baseRunner

	config TableConfig
	get    func(in any) gottype

	rfunc *reflectutil.Func
	args  Args
}

func (r *tableRunner) Run(t *testing.T) {
	t.Helper()
	cond.PanicOnErr(r.setGetFunc())
	r.run(t)
}

func (r *tableRunner) DryRun() TableResulter {
	cond.PanicOnErr(r.setGetFunc())
	return tableResults{baseResults: r.dryRun()}
}

func (r *tableRunner) setGetFunc() error {
	if err := r.validateConfig(); err != nil {
		return err
	}

	args, err := r.makeFixedArgs(r.rfunc, r.config)
	if err != nil {
		return err
	}

	r.get = func(in any) gottype {
		pin, pout := r.config.InPos, r.config.OutPos
		r.args = args.replaceAt(pin, in)
		return r.rfunc.Call(r.args)[pout]
	}
	return nil
}

func (r *tableRunner) Cases(cases []Case) TableRunner {
	for i, tc := range cases {
		i, tc := i, tc

		get := func() gottype { return r.get(tc.In) }

		getLabel := func() string {
			return fmtexpl.TableCaseLabel(r.rfunc.Name, i, tc.Lab, r.args)
		}

		addCaseCheck := func(c check.Checker[any]) {
			r.addCheck(baseCheck{
				get:      get,
				getLabel: getLabel,
				label:    tc.Lab,
				checker:  c,
			})
		}

		// add Case.Exp check
		if tc.Exp != nil {
			exp := cond.Value(nil, tc.Exp, tc.Exp == ExpNil)
			addCaseCheck(check.Value.Is(exp))
		}

		// add Case.Not checks
		if len(tc.Not) != 0 {
			addCaseCheck(check.Value.Not(tc.Not...))
		}

		// add Case.Pass checks
		if len(tc.Pass) != 0 {
			r.addChecks(tc.Lab, get, tc.Pass)
		}
	}
	return r
}

func (r *tableRunner) Config(cfg TableConfig) TableRunner {
	r.config = cfg
	return r
}

func (r *tableRunner) setRfunc(in any) error {
	rfunc, err := reflectutil.NewFunc(in)
	if err != nil {
		return fmt.Errorf("Table(func): %w", err)
	}
	ftype := rfunc.Value.Type()
	if ftype.NumIn() == 0 {
		return fmt.Errorf("Table(%s): %w", rfunc.Name, errTableRunnerFuncNumIn)
	}
	if ftype.NumOut() == 0 {
		return fmt.Errorf("Table(%s): %w", rfunc.Name, errTableRunnerFuncNumOut)
	}
	r.rfunc = rfunc
	return nil
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

	switch d := nparams - nargs; d {
	case 0:
		return cfg.FixedArgs, nil
	case 1:
		args := make(Args, nparams)
		delta := 0
		for i := 0; i < nparams; i++ {
			if i == cfg.InPos {
				delta++
				continue
			}
			args[i] = cfg.FixedArgs[i-delta]
		}
		return args, nil
	default:
		return nil, errTableRunnerConfigFixedArgs(d)
	}
}

func newTableRunner(testedFunc any) TableRunner {
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

type expNil int

// ExpNil is a special value for Case.Exp that sets the expected value to nil.
const ExpNil expNil = 0
