// Package fmtexpl provides functions to format failed test explanations.
package fmtexpl

import "fmt"

// Default computes and return an explain string in the default format.
func Default(label string, exp, got interface{}) string {
	return fmt.Sprintf("%s:\nexp %v\ngot %v", label, exp, got)
}

// Debug computes and return a detailed explain string based on the
// default format.
func Debug(label string, exp, got interface{}) string {
	return Default(label, fmt.Sprintf("%#v", exp), fmt.Sprintf("%#v", got))
}

// Checker computes and return an explain string based on a gotten
// checker explanation.
func Checker(label, expStr, gotExpl string) string {
	return Default(label, expStr, "explanation: "+gotExpl)
}

// Pretty computes and return an explain string in a pretty output.
// Compatibility is not guaranteed, usage should be restricted to
// local tests -- not in checkers explain func.
func Pretty(label string, exp, got interface{}) string {
	return "❌ " + Default(label, exp, got)
}

// FuncResult computes and return an explain string for an output
// from a func call.
// TODO: handle multiple func args/results
func FuncResult(fname, desc string, arg, exp, got interface{}) string {
	label := fmt.Sprintf("%s(%v)", fname, arg)
	if desc != "" {
		label = "[" + desc + "] " + label
	}
	return Default(label, exp, got)
}
