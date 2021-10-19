// Package fmtexpl provides functions to format failed test explanations.
package fmtexpl

import (
	"fmt"

	"github.com/drykit-go/cond"
)

// Default computes and return an explain string in the default format.
func Default(label string, exp, got interface{}) string {
	return fmt.Sprintf("%s:\nexp %v\ngot %v", label, exp, got)
}

// Pretty computes and return an explain string in a pretty output.
// Compatibility is not guaranteed, usage should be restricted to
// local tests -- not in checkers explain func.
func Pretty(label string, exp, got interface{}) string {
	return "❌ " + Default(label, exp, got)
}

// Checker computes and return an explain string based on a gotten
// checker explanation.
func Checker(label, expStr, gotExpl string) string {
	return Default(label, expStr, "explanation: "+gotExpl)
}

// TableCaseLabel returns the label for a testx.Table test case
// in format: Case <caseID> "<caseLab>" <fname>(<caseIn>)
//
// Examples:
// 	`Table.Cases[2] isEven(42)`
// 	`Table.Cases[2] "even number" isEven(42)`
func TableCaseLabel(
	fname string,
	caseID int,
	caseLab string,
	args fmt.Stringer,
) string {
	fcall := fmt.Sprintf("%s(%v)", fname, args)
	label := cond.String(fmt.Sprintf(` "%s"`, caseLab), "", caseLab != "")
	return fmt.Sprintf("Table.Cases[%d]%s %s", caseID, label, fcall)
}
