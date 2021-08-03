package check

import "fmt"

type untypedCheck struct {
	passFunc UntypedPassFunc
	explFunc ExplainFunc
}

func (c untypedCheck) Pass(got interface{}) bool {
	return c.passFunc(got)
}

func (c untypedCheck) Explain(label string, got interface{}) string {
	return c.explFunc(label, got)
}

type untypedValue struct{}

func (untypedValue) Custom(desc string, f UntypedPassFunc) UntypedChecker {
	return untypedCheck{
		passFunc: f,
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to %s, got %v",
				label, desc, got,
			)
		},
	}
}
