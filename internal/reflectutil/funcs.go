package reflectutil

import (
	"fmt"
	"log"
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
func FuncName(fn interface{}) string {
	return funcName(reflect.ValueOf(fn))
}

func funcName(fval reflect.Value) string {
	if !IsFunc(fval) {
		log.Panicf(
			"calling FuncName with a non func argument (%s %v)",
			fval.Type().String(), fval.String(),
		)
	}
	fptr := fval.Pointer()
	path := runtime.FuncForPC(fptr).Name()
	return filepath.Base(path)
}

type FuncSignature struct {
	Name    string
	In, Out []reflect.Kind
}

func (s FuncSignature) Match(ftyp reflect.Type) bool {
	return s.matchIn(ftyp) && s.matchOut(ftyp)
}

func (s FuncSignature) ImplementedBy(v reflect.Value) bool {
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

type Func struct {
	Name string
	Rtyp reflect.Type
	Rval reflect.Value
}

func NewFunc(fn interface{}) (*Func, error) {
	fval := reflect.ValueOf(fn)
	if !IsFunc(fval) {
		return nil, fmt.Errorf("%w: got %v", ErrNotAFunc, fn)
	}
	return &Func{
		Name: funcName(fval),
		Rtyp: fval.Type(),
		Rval: fval,
	}, nil
}
