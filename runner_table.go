package testx

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/drykit-go/cond"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
	"github.com/drykit-go/testx/internal/fmtexpl"
	"github.com/drykit-go/testx/internal/reflectutil"
)

var _ TableRunner = (*tableRunner)(nil)

// Case represents a Table test case. It must be provided an In value
// and an Exp value at least.
type Case struct {
	// Lab is the label of the current case. If provided it will be printed
	// if a test case fails.
	Lab string

	// In is the input value to the tested func.xz
	In interface{}

	// Exp is the value expected to be returned when calling the tested func.
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

	// OutPos is the nth output value that is tested against Exp.
	// It is required if the tested func returns multiple values.
	// Default is 0.
	OutPos int

	// FixedArgs is a slice of arguments to be injected into the tested func.
	// It is required if the tested func accepts multiple parameters.
	// len(FixedArgs) must equal nparams or nparams - 1, nparams being
	// the number of parameters of the tested func.
	// Depending on the value of len(FixedArgs) - nparams, Case.In will either
	// replace (0) or be inserted (-1) at index InPos.
	FixedArgs []interface{}
}

type tableRunner struct {
	baseRunner

	label  string
	config TableConfig
	get    func(in interface{}) gottype

	refFunc *reflectutil.Func
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

func (r *tableRunner) setGetFunc(args []reflect.Value) {
	r.get = func(in interface{}) gottype {
		args[r.config.InPos] = reflect.ValueOf(in)
		out := r.refFunc.Rval.Call(args)
		got := out[r.config.OutPos]
		return got.Interface()
	}
}

func (r *tableRunner) validateConfig() error {
	if r.config.InPos < 0 {
		return fmt.Errorf(
			"invalid value for TableConfig.InPos: exp >= 0, got %d",
			r.config.InPos,
		)
	}
	if r.config.OutPos < 0 {
		return fmt.Errorf(
			"invalid value for TableConfig.OutPos: exp >= 0, got %d",
			r.config.OutPos,
		)
	}
	if outPos, numOut := r.config.OutPos, r.refFunc.Rtyp.NumOut(); outPos >= numOut {
		return fmt.Errorf(
			"invalid value for OutPos: must be < to the number of values returned by %s (%d > %d)",
			r.refFunc.Name, outPos, numOut,
		)
	}
	if inPos, numIn := r.config.InPos, r.refFunc.Rtyp.NumIn(); inPos >= numIn {
		return fmt.Errorf(
			"invalid value for InPos: must be < to the number of parameters of %s (%d > %d)",
			r.refFunc.Name, inPos, numIn,
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

func (r *tableRunner) setRefFunc(in interface{}) error {
	rf, err := reflectutil.NewFunc(in)
	if err != nil {
		return fmt.Errorf("Table(func): %w", err)
	}
	if rf.Rtyp.NumIn() == 0 {
		return fmt.Errorf(
			"Table(%s): invalid func: it must accept at least 1 parameter",
			rf.Name,
		)
	}
	if rf.Rtyp.NumOut() == 0 {
		return fmt.Errorf(
			"Table(%s): invalid func: it must return at least 1 value",
			rf.Name,
		)
	}
	r.refFunc = rf
	return nil
}

func (r *tableRunner) makeFixedArgs() ([]reflect.Value, error) {
	nparams := r.refFunc.Rtyp.NumIn()
	nargs := len(r.config.FixedArgs)

	fillskip := func(at int) []reflect.Value {
		args := make([]reflect.Value, nparams)
		delta := 0
		for i := 0; i < nparams; i++ {
			if i == at {
				delta++
				continue
			}
			args[i] = reflect.ValueOf(r.config.FixedArgs[i-delta])
		}
		return args
	}

	switch d := nparams - nargs; d {
	case 0:
		return reflectutil.WrapValues(r.config.FixedArgs), nil
	case 1:
		return fillskip(r.config.InPos), nil
	default:
		return nil, fmt.Errorf("invalid FixedArgs number: %d", nargs)
	}
}

// Table returns a TableRunner that runs test cases provided via
// Cases() method on the given testedFunc. A TableConfig may be
// required if the testedFunc accepts several parameters or returns
// several values.
func newTableRunner(testedFunc interface{}, cfg *TableConfig) TableRunner {
	r := tableRunner{}
	r.setConfig(cfg)

	cond.PanicOnErr(
		r.setRefFunc(testedFunc),
		r.validateConfig(),
	)

	args, err := r.makeFixedArgs()
	cond.PanicOnErr(err)

	r.label = r.refFunc.Name
	r.setGetFunc(args)

	return &r
}

type tableResults struct {
	baseResults
}

func (res tableResults) PassedAt(i int) bool {
	if i >= len(res.checks) {
		log.Panicf("Provided index %d is out of range", i)
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
	log.Panicf("No test case with label %s", label)
	return false
}

func (res tableResults) FailedLabel(label string) bool {
	return !res.PassedLabel(label)
}
