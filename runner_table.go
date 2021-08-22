package testx

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
)

var _ TableRunner = (*tableRunner)(nil)

type tableRunner struct {
	baseRunner

	label  string
	config TableConfig
	get    func(in interface{}) gotType
	cases  []Case
}

func (r *tableRunner) Run(t *testing.T) {
	for _, c := range r.cases {
		c := c
		r.addCheck(testCheck{
			label: defaultString(c.Lab, "value"), // TODO: check for unexpected formattings
			get:   func() gotType { return r.get(c.In) },
			check: r.makeChecker(c),
		})
	}
	r.run(t)
}

func (r *tableRunner) Cases(cases []Case) TableRunner {
	r.cases = append(r.cases, cases...)
	return r
}

func (r *tableRunner) Config(cfg *TableConfig) TableRunner {
	r.setConfig(cfg)
	return r
}

func (r *tableRunner) Label(label string) TableRunner {
	r.label = label
	return r
}

func (r *tableRunner) FixedArgs(args ...interface{}) TableRunner {
	r.config.FixedArgs = args
	return r
}

func (r *tableRunner) InPos(pos int) TableRunner {
	r.config.InPos = pos
	return r
}

func (r *tableRunner) OutPos(pos int) TableRunner {
	r.config.OutPos = pos
	return r
}

func Table(testedFunc interface{}, cfg *TableConfig) TableRunner {
	r := tableRunner{}
	r.setConfig(cfg)

	f, err := r.makeFuncReflection(testedFunc)
	panicOnErr(err)

	panicOnErr(r.validateConfig(f))

	args, err := r.makeArgs(f, r.config)
	panicOnErr(err)

	r.label = f.name
	r.setGetFunc(f, args)

	return &r
}

func (r *tableRunner) setConfig(cfg *TableConfig) {
	if cfg == nil {
		r.config = TableConfig{}
	} else {
		r.config = *cfg
	}
}

func (r *tableRunner) setGetFunc(f funcReflection, args []reflect.Value) {
	r.get = func(in interface{}) gotType {
		args[r.config.InPos] = reflect.ValueOf(in)
		out := f.rval.Call(args)
		got := out[r.config.OutPos]
		return got.Interface()
	}
}

func (r *tableRunner) validateConfig(f funcReflection) error {
	if r.config.InPos < 0 {
		return fmt.Errorf(
			"invalid value for InPos: must be >= 0, got %d",
			r.config.InPos,
		)
	}
	if r.config.OutPos < 0 {
		return fmt.Errorf(
			"invalid value for OutPos: must be >= 0, got %d",
			r.config.OutPos,
		)
	}
	if outPos, numOut := r.config.OutPos, f.rtyp.NumOut(); outPos >= numOut {
		return fmt.Errorf(
			"invalid value for OutPos: must be < to the number of values returned by %s (%d > %d)",
			f.name, outPos, numOut,
		)
	}
	if inPos, numIn := r.config.InPos, f.rtyp.NumIn(); inPos >= numIn {
		return fmt.Errorf(
			"invalid value for InPos: must be < to the number of parameters of %s (%d > %d)",
			f.name, inPos, numIn,
		)
	}
	return nil
}

func (r *tableRunner) makeChecker(c Case) check.UntypedChecker {
	var (
		pass check.UntypedPassFunc
		expl check.ExplainFunc
	)

	if checkconv.IsChecker(c.Exp) {
		pass = func(got interface{}) bool {
			gotv := reflect.ValueOf(got)
			expv := reflect.ValueOf(c.Exp)
			outv := expv.MethodByName("Pass").Call([]reflect.Value{gotv})
			return outv[0].Bool()
		}
		expl = func(label string, got interface{}) string {
			gotv := reflect.ValueOf(got)
			expv := reflect.ValueOf(c.Exp)
			labv := reflect.ValueOf(defaultString(label, "value"))
			expl := expv.MethodByName("Explain").Call([]reflect.Value{labv, gotv})[0]
			return fmt.Sprintf(
				"%s(%v) -> checker returned the following error:\n%s",
				r.label, c.In, expl.String(),
			)
		}
	} else {
		pass = func(got interface{}) bool {
			return xor(deq(got, c.Exp), c.Not)
		}
		expl = func(label string, got interface{}) string {
			return fmt.Sprintf(
				"%s(%v) -> expect %s%v, got %v",
				r.label, c.In, condString("not ", "", c.Not), c.Exp, got,
			)
		}
	}
	return check.NewUntypedCheck(pass, expl)
}

func (r *tableRunner) makeFuncReflection(in interface{}) (funcReflection, error) {
	rval := reflect.ValueOf(in)
	rtyp := rval.Type()

	if kind := rtyp.Kind(); kind != reflect.Func {
		return funcReflection{}, fmt.Errorf(
			"calling Table(f) with a non func argument (got %s %s)",
			rtyp.String(), rval.String(),
		)
	}

	name := getFuncName(in)

	if rtyp.NumIn() == 0 {
		return funcReflection{}, fmt.Errorf(
			"%s is not a valid func: it must accept at least 1 parameter",
			name,
		)
	}

	if rtyp.NumOut() == 0 {
		return funcReflection{}, fmt.Errorf(
			"%s is not a valid func: it must return at least 1 value",
			name,
		)
	}

	return funcReflection{
		name: name,
		rval: rval,
		rtyp: rtyp,
	}, nil
}

func (r *tableRunner) makeArgs(f funcReflection, cfg TableConfig) ([]reflect.Value, error) {
	nparams := f.rtyp.NumIn()
	nargs := len(cfg.FixedArgs)
	args := make([]reflect.Value, nparams)

	fillskip := func(at int) []reflect.Value {
		for i := 0; i < nparams; i++ {
			var v interface{}
			switch {
			case i == at:
				v = nil
			case i > at:
				v = cfg.FixedArgs[i-1]
			default:
				v = cfg.FixedArgs[i]
			}
			args[i] = reflect.ValueOf(v)
		}
		return args
	}

	fillall := func() []reflect.Value {
		for i, v := range cfg.FixedArgs {
			args[i] = reflect.ValueOf(v)
		}
		return args
	}

	switch d := nparams - nargs; d {
	case 0:
		return fillall(), nil
	case 1:
		return fillskip(cfg.InPos), nil
	default:
		return nil, fmt.Errorf("invalid FixedArgs number: %d", nargs)
	}
}

// Case represents a Table test case. It must be provided an In value
// and an Exp value at least.
type Case struct {
	// Lab is the label of the current case. If provided it will be printed
	// if a test case fails.
	Lab string

	// In is the input value to the tested func.
	In interface{}

	// Exp is the value expected to be returned when calling the tested func.
	Exp interface{}

	// Not reverses the test check for an equality.
	Not bool
}

type funcReflection struct {
	name string
	rval reflect.Value
	rtyp reflect.Type
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
