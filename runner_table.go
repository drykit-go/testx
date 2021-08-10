package testx

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/drykit-go/testx/check/checkconv"
)

// TODO: refactor funcs -> tableRunner methods, reorder

type TableRunner interface {
	Runner
	Cases(cases []Case) TableRunner
}

type tableRunner struct {
	t      *testing.T
	label  string
	get    func(interface{}) interface{}
	cases  []Case
	config *TableConfig
}

func (r *tableRunner) Run(t *testing.T) {
	r.t = t
	for _, c := range r.cases {
		got := r.get(c.In)
		if pass, expl := r.pass(c, got); !pass {
			r.fail(expl)
		}
	}
}

// TODO: once results data is stored into tableRunner,
// change back the return value to bool only
// TODO: method of Case instead of tableRunner?
// -> implies to add testName to Case struct for default label
func (r *tableRunner) pass(c Case, got interface{}) (bool, string) {
	switch {
	case checkconv.IsChecker(c.Exp):
		return r.passChecker(c, got)
	default:
		return r.passValue(c, got)
	}
}

func (r *tableRunner) passValue(c Case, got interface{}) (bool, string) {
	pass, expl := xor(deq(got, c.Exp), c.Not), ""
	if !pass {
		expl = r.explValue(c.In, got, c.Exp, c.Not)
	}
	return pass, expl
}

func (r *tableRunner) explValue(in, got, exp interface{}, not bool) string {
	return fmt.Sprintf(
		"%s(%v) -> expect %s%v, got %v",
		r.label, in, condString("not ", "", not), exp, got,
	)
}

// NOTE: passChecker does not support c.Not, because c.Not is to be removed.
func (r *tableRunner) passChecker(c Case, got interface{}) (bool, string) {
	gotv := reflect.ValueOf(got)
	expv := reflect.ValueOf(c.Exp)
	outv := expv.MethodByName("Pass").Call([]reflect.Value{gotv})
	pass, expl := outv[0].Bool(), ""
	if !pass {
		expl = r.explChecker(c.Lab, c.In, gotv, expv)
	}
	return pass, expl
}

func (r *tableRunner) explChecker(label string, in interface{}, gotv, expv reflect.Value) string {
	labv := reflect.ValueOf(defaultString(label, "value"))
	expl := expv.MethodByName("Explain").Call([]reflect.Value{labv, gotv})[0]
	return fmt.Sprintf(
		"%s(%v) -> checker returned the following error:\n%s",
		r.label, in, expl.String(),
	)
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

func (r *tableRunner) fail(expl string) {
	r.t.Error(expl)
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

	Check interface{}
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
