package testx

import (
	"fmt"
	"reflect"
	"testing"
)

// TODO: refactor funcs -> tableRunner methods, reorder

type tableRunner struct {
	t      *testing.T
	label  string
	get    func(interface{}) interface{}
	cases  []Case
	config *TableConfig

	result struct {
		npass int
		nfail int
	}
}

func (r *tableRunner) Run(t *testing.T) {
	deq := func(a, b interface{}) bool { return reflect.DeepEqual(a, b) }
	pass := func(eq, reverse bool) bool { return eq != reverse }

	r.t = t
	for _, c := range r.cases {
		if got := r.get(c.In); !pass(deq(got, c.Exp), c.Not) {
			r.fail(c.In, got, c.Exp, c.Not)
		}
	}
}

func (r *tableRunner) Cases(cases []Case) TableRunner {
	r.cases = append(r.cases, cases...)
	return r
}

// TODO:
func (r *tableRunner) Config(config *TableConfig) TableRunner {
	if config != nil {
		r.config = config
	}
	return r
}

// TODO:
func (r *tableRunner) Label(label string) TableRunner {
	r.label = label
	return r
}

func (r *tableRunner) fail(in, got, exp interface{}, not bool) {
	notStr := func(not bool) string {
		if not {
			return "not "
		}
		return ""
	}
	r.t.Errorf("%s(%v) -> expect %s%v, got %v",
		r.label, in, notStr(not), exp, got,
	)
}

type Case struct {
	// Lab is the label of the current case
	Lab string
	// In is the input value to the tested func
	In interface{}
	// Exp is the expected value expected to be returned by Get.
	Exp interface{}
	// Not reverses the test check for an equality
	Not bool
}

type TableRunner interface {
	Runner
	Cases(cases []Case) TableRunner
}

func Table(testedFunc interface{}, config *TableConfig) TableRunner {
	cfg, err := safeTableConfig(config)
	panicOnErr(err)

	f, err := newFuncReflection(testedFunc, cfg.Name)
	panicOnErr(err)

	panicOnErr(cfg.validateFuncCompat(f))

	args, err := makeArgs(f, &cfg)
	panicOnErr(err)

	return &tableRunner{
		label: f.name,
		get: func(in interface{}) interface{} {
			args[cfg.InPos] = reflect.ValueOf(in)
			out := f.rval.Call(args)
			got := out[cfg.OutPos]
			return got.Interface()
		},
		config: &cfg,
	}
}

func (cfg *TableConfig) validateFuncCompat(f *funcReflection) error {
	if outPos, numOut := cfg.OutPos, f.rtyp.NumOut(); outPos >= numOut {
		return fmt.Errorf(
			"invalid value for OutPos: must be < to the number of values returned by %s (%d > %d)",
			f.name, outPos, numOut,
		)
	}
	if inPos, numIn := cfg.InPos, f.rtyp.NumIn(); inPos >= numIn {
		return fmt.Errorf(
			"invalid value for InPos: must be < to the number of parameters of %s (%d > %d)",
			f.name, inPos, numIn,
		)
	}
	return nil
}

func safeTableConfig(cfg *TableConfig) (TableConfig, error) {
	if cfg == nil {
		return TableConfig{}, nil
	}
	if cfg.InPos < 0 {
		return TableConfig{}, fmt.Errorf(
			"invalid value for InPos: must be int >= 0, got %d",
			cfg.InPos,
		)
	}
	if cfg.OutPos < 0 {
		return TableConfig{}, fmt.Errorf(
			"invalid value for OutPos: must be int >= 0, got %d",
			cfg.OutPos,
		)
	}
	return *cfg, nil
}

type funcReflection struct {
	name string
	rval reflect.Value
	rtyp reflect.Type
}

func newFuncReflection(in interface{}, name string) (*funcReflection, error) {
	rval := reflect.ValueOf(in)
	rtyp := rval.Type()

	if kind := rtyp.Kind(); kind != reflect.Func {
		return nil, fmt.Errorf(
			"calling Table(f) with a non func argument (got %s %s)",
			rtyp.String(), rval.String(),
		)
	}

	if name == "" {
		name = getFuncName(in)
	}

	if rtyp.NumIn() == 0 {
		return nil, fmt.Errorf(
			"%s is not a valid func: it must accept at least 1 parameter(s)",
			name,
		)
	}

	if rtyp.NumOut() == 0 {
		return nil, fmt.Errorf(
			"%s is not a valid func: it must return at least 1 value(s)",
			name,
		)
	}

	return &funcReflection{
		name: name,
		rval: rval,
		rtyp: rtyp,
	}, nil
}

func makeArgs(f *funcReflection, cfg *TableConfig) ([]reflect.Value, error) {
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
		// TODO: + specific
		return nil, fmt.Errorf("invalid FixedArgs number: %d", nargs)
	}
}

func funcName(f interface{}, name string) string {
	if name != "" {
		return name
	}
	return getFuncName(f)
}

// TableConfig is an object of options allowing to configure
// the table runner. Its main purpose is to deal with tested functions
// that have multiple input parameters or outputs by setting
// the injection position of In value, the output position to read
// the obtained value from, and a slice of fixed arguments.
type TableConfig struct {
	// Name is a custom name for the given tested func.
	// Default name is computed in format mypackage.MyFunc.
	Name string

	// InPos is the nth parameter in which In value is injected.
	// It must be used when the tested function accepts several parameters
	// and you want to inject the cases input at nth position.
	// Default is 0.
	InPos int

	// OutPos is the nth output value that is tested against Exp.
	// It must be used when the tested function returns several values
	// and you want to test the nth returned value.
	// Default is 0.
	OutPos int

	// Args is a mapping of parameter position with fixed values.
	// These values will be used for all the cases.
	// It must be used with InPos to work expectedly.
	FixedArgs []interface{}
}
