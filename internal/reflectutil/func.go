package reflectutil

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
)

// IsFunc returns true if v is a func, else false.
func IsFunc(v reflect.Value) bool {
	return v.Kind() == reflect.Func
}

// FuncName returns the name of the given func prefixed with
// the name of the package it is from. It panics if f is not a func.
func FuncName(fn any) string {
	return funcName(reflect.ValueOf(fn))
}

func funcName(fval reflect.Value) string {
	if !IsFunc(fval) {
		panic(fmt.Sprintf(
			"calling FuncName with a non func argument (%s)",
			fval.Type().String(),
		))
	}
	fptr := fval.Pointer()
	path := runtime.FuncForPC(fptr).Name()
	return filepath.Base(path)
}

// Func is a FuncSignature associated with its actual reflect.Value.
type Func struct {
	FuncSignature
	Value reflect.Value
}

// Call calls Func's underlying func with given args and returns the results
// as a slice of empty interfaces.
func (f *Func) Call(args []any) []any {
	return UnwrapValues(f.Value.Call(WrapValues(args)))
}

// NewFunc returns a *Func from the given func input, or a non-nil error
// if fn's kind is not reflect.Func.'
func NewFunc(fn any) (*Func, error) {
	fval := reflect.ValueOf(fn)
	if !IsFunc(fval) {
		return nil, fmt.Errorf("%w: got %v", ErrNotAFunc, fn)
	}
	return &Func{
		FuncSignature: FuncSignature{Name: funcName(fval)},
		Value:         fval,
	}, nil
}

// FuncSignature is a func signature having a name and In/Out types represented
// by slices of reflect.Kind.
type FuncSignature struct {
	Name    string
	In, Out []reflect.Kind
}

// Match returns true if ftyp in/out kinds match the FuncSignature ones.
func (s FuncSignature) Match(ftyp reflect.Type) bool {
	return s.matchIn(ftyp) && s.matchOut(ftyp)
}

// ImplementedBy returns true if v's type has a method that matches
// FuncSignature's Name and In/Out kinds.
func (s FuncSignature) ImplementedBy(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}
	m := v.MethodByName(s.Name)
	if !m.IsValid() {
		return false
	}
	return s.Match(m.Type())
}

func (s FuncSignature) matchIn(t reflect.Type) bool {
	return s.matchValues(t.NumIn(), t.In, s.In)
}

func (s FuncSignature) matchOut(t reflect.Type) bool {
	return s.matchValues(t.NumOut(), t.Out, s.Out)
}

func (s FuncSignature) matchValues(
	numValues int,
	getIthVal func(int) reflect.Type,
	expKinds []reflect.Kind,
) bool {
	if numValues != len(expKinds) {
		return false
	}
	for i := 0; i < numValues; i++ {
		if !s.validKind(getIthVal(i).Kind(), expKinds[i]) {
			return false
		}
	}
	return true
}

func (s FuncSignature) validKind(gotk, expk reflect.Kind) bool {
	return expk == AnyKind || gotk == expk
}
