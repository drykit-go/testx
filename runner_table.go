package testx

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/drykit-go/cond"
	"github.com/drykit-go/slices"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/internal/fmtexpl"
	"github.com/drykit-go/testx/internal/reflectutil"
)

// Case represents a Table test case. It must be provided values for
// Case.In at least.
type Case[In, Exp any] struct {
	// Lab is the label of the current case to be printed if the current
	// case fails.
	Lab string

	// In is the input value injected in the tested func.
	In In

	// Exp is the value expected to be returned when calling the tested func.
	// It is ignored if Case.Not or Case.Pass is not empty.
	// Otherwise Case.Exp is always evaluated even if not set.
	Exp Exp

	// Not is a slice of values expected not to be returned by the tested func.
	// If set, Case.Exp is ignored.
	Not []Exp

	// Pass is a slice of checkers that the return value of the
	// tested func is expected to pass.
	// If set, Case.Exp is ignored.
	Pass []check.Checker[Exp]
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

type tableRunner[In, Exp any] struct {
	baseRunner

	config TableConfig
	get    func(in In) Exp

	rfunc *reflectutil.Func
	args  Args
}

func (r *tableRunner[In, Exp]) Run(t *testing.T) {
	t.Helper()
	cond.PanicOnErr(r.setGetFunc())
	r.run(t)
}

func (r *tableRunner[In, Exp]) DryRun() TableResulter {
	cond.PanicOnErr(r.setGetFunc())
	return tableResults{baseResults: r.dryRun()}
}

func (r *tableRunner[In, Exp]) setGetFunc() error {
	if err := r.validateConfig(); err != nil {
		return err
	}

	args, err := r.makeFixedArgs(r.rfunc, r.config)
	if err != nil {
		return err
	}

	r.get = func(in In) Exp {
		pin, pout := r.config.InPos, r.config.OutPos
		r.args = args.replaceAt(pin, in)
		out, ok := r.rfunc.Call(r.args)[pout].(Exp)
		if !ok {
			var vnil Exp
			return vnil
		}
		return out
	}
	return nil
}

func (r *tableRunner[In, Exp]) Cases(cases []Case[In, Exp]) TableRunner[In, Exp] {
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

		hasNotChecks := len(tc.Not) != 0
		hasPassChecks := len(tc.Pass) != 0
		hasExpCheck := !hasNotChecks && !hasPassChecks

		// add Case.Not checks
		if hasNotChecks {
			addCaseCheck(check.Value[any]().Not(slices.AsAny(tc.Not)...))
		}

		// add Case.Pass checks
		if hasPassChecks {
			r.addChecks(tc.Lab, get, check.WrapMany(tc.Pass...))
		}

		// add Case.Exp checks
		if hasExpCheck {
			addCaseCheck(check.Value[any]().Is(tc.Exp))
		}
	}
	return r
}

func (r *tableRunner[In, Exp]) Config(cfg TableConfig) TableRunner[In, Exp] {
	r.config = cfg
	return r
}

func (r *tableRunner[In, Exp]) setRfunc(in any) error {
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

func (r *tableRunner[In, Exp]) validateConfig() error {
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

func (r *tableRunner[In, Exp]) makeFixedArgs(rfunc *reflectutil.Func, cfg TableConfig) (Args, error) {
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

func newTableRunner[In, Exp any](testedFunc any) TableRunner[In, Exp] {
	r := &tableRunner[In, Exp]{}
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
