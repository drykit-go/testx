package testx

import (
	"fmt"
	"log"
	"testing"

	"github.com/drykit-go/cond"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
	"github.com/drykit-go/testx/internal/fmtexpl"
	"github.com/drykit-go/testx/internal/reflectutil"
)

var _ TableRunner = (*tableRunner)(nil)

// Case represents a Table test case. It must be provided values for
// Case.In and Case.Exp at least.
type Case struct {
	// Lab is the label of the current case to be printed if the current
	// case fails.
	Lab string

	// In is the input value injected in the tested func.
	In interface{}

	// Exp is the value expected to be returned when calling the tested func.
	// If Exp is a checker, the checker is run instead.
	Exp interface{}

	// Not reverses the test check for an equality.
	Not bool
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
	// 		FixedArgs: []interface{0: "myArg0", 2: "myArg2"} // len(FixedArgs == 3)
	// 	)
	FixedArgs Args
}

type Args []interface{}

func (args Args) replaceAt(pos int, arg interface{}) Args {
	if pos >= len(args) {
		log.Panic("Args.set(i, v): i is out of range")
	}
	args[pos] = arg
	return args
}

type tableRunner struct {
	baseRunner

	label  string
	config TableConfig
	get    func(in interface{}) gottype

	rfunc *reflectutil.Func
}

func (r *tableRunner) Run(t *testing.T) {
	t.Helper()
	r.run(t)
}

func (r *tableRunner) DryRun() TableResulter {
	return tableResults{baseResults: r.baseResults()}
}

func (r *tableRunner) Cases(cases []Case) TableRunner {
	for _, c := range cases {
		c := c
		r.addCheck(baseCheck{
			label:   c.Lab,
			get:     func() gottype { return r.get(c.In) },
			checker: r.makeChecker(c),
		})
	}
	return r
}

func (r *tableRunner) Config(cfg *TableConfig) TableRunner {
	r.setConfig(cfg)
	return r
}

func (r *tableRunner) setConfig(cfg *TableConfig) {
	if cfg == nil {
		r.config = TableConfig{}
	} else {
		r.config = *cfg
	}
}

func (r *tableRunner) setRfunc(in interface{}) error {
	rfunc, err := reflectutil.NewFunc(in)
	if err != nil {
		return fmt.Errorf("Table(func): %w", err)
	}
	ftype := rfunc.Value.Type()
	if ftype.NumIn() == 0 {
		return fmt.Errorf(
			"Table(%s): invalid func: it must accept at least 1 parameter",
			rfunc.Name,
		)
	}
	if ftype.NumOut() == 0 {
		return fmt.Errorf(
			"Table(%s): invalid func: it must return at least 1 value",
			rfunc.Name,
		)
	}
	r.rfunc = rfunc
	return nil
}

func (r *tableRunner) setGetFunc(args Args) {
	r.get = func(in interface{}) gottype {
		pin, pout := r.config.InPos, r.config.OutPos
		return r.rfunc.Call(args.replaceAt(pin, in))[pout]
	}
}

func (r *tableRunner) validateConfig() error {
	validPos := func(pos, max int) bool { return pos >= 0 && pos < max }
	ftyp := r.rfunc.Value.Type()
	if pout, nout := r.config.OutPos, ftyp.NumOut(); !validPos(pout, nout) {
		return fmt.Errorf(
			"%w: OutPos: exp 0 <= n < %d (number of values returned by %s), got %d",
			ErrInvalidTableConfig, nout, r.rfunc.Name, pout,
		)
	}
	if pin, nin := r.config.InPos, ftyp.NumIn(); !validPos(pin, nin) {
		return fmt.Errorf(
			"%w: InPos: exp 0 <= n < %d (number of parameters of %s), got %d",
			ErrInvalidTableConfig, nin, r.rfunc.Name, pin,
		)
	}
	return nil
}

func (r *tableRunner) makeChecker(c Case) check.ValueChecker {
	checker, ok := checkconv.Cast(c.Exp)
	if ok {
		return checker
	}

	pass := func(got interface{}) bool { return xor(deq(got, c.Exp), c.Not) }
	expl := func(_ string, got interface{}) string {
		expStr := fmt.Sprintf("%s%v", cond.String("not ", "", c.Not), c.Exp)
		return fmtexpl.FuncResult(r.label, c.Lab, c.In, expStr, got)
	}
	return check.NewValueChecker(pass, expl)
}

func (r *tableRunner) makeFixedArgs() (Args, error) {
	nparams := r.rfunc.Value.Type().NumIn()
	nargs := len(r.config.FixedArgs)

	fillskip := func(at int) Args {
		args := make(Args, nparams)
		delta := 0
		for i := 0; i < nparams; i++ {
			if i == at {
				delta++
				continue
			}
			args[i] = r.config.FixedArgs[i-delta]
		}
		return args
	}

	switch d := nparams - nargs; d {
	case 0:
		return r.config.FixedArgs, nil
	case 1:
		return fillskip(r.config.InPos), nil
	default:
		return nil, fmt.Errorf("invalid FixedArgs number: %d", nargs)
	}
}

func newTableRunner(testedFunc interface{}, cfg *TableConfig) TableRunner {
	r := tableRunner{}
	r.setConfig(cfg)

	cond.PanicOnErr(r.setRfunc(testedFunc))
	cond.PanicOnErr(r.validateConfig())

	args, err := r.makeFixedArgs()
	cond.PanicOnErr(err)

	r.label = r.rfunc.Name
	r.setGetFunc(args)

	return &r
}

/*
	Results
*/

type tableResults struct {
	baseResults
}

func (res tableResults) PassedAt(i int) bool {
	if i >= len(res.checks) {
		log.Panicf("TableResults: index %d is out of range", i)
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
	log.Panicf("TableResults: no test case with label %s", label)
	return false
}

func (res tableResults) FailedLabel(label string) bool {
	return !res.PassedLabel(label)
}
