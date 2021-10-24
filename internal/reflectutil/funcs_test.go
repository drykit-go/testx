package reflectutil_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/drykit-go/testx/internal/reflectutil"
	"github.com/drykit-go/testx/internal/testutil"
)

func validFunc(in int) int {
	return in
}

func TestFuncName(t *testing.T) {
	t.Run("non-func panics", func(t *testing.T) {
		nofunc := struct{}{}
		defer testutil.AssertPanicMessage(t,
			"calling FuncName with a non func argument (struct {})",
		)
		if got, exp := reflectutil.FuncName(nofunc), ""; got != exp {
			t.Errorf("exp %s, got %s", exp, got)
		}
	})

	t.Run("returns func name", func(t *testing.T) {
		got := reflectutil.FuncName(validFunc)
		exp := "reflectutil_test.validFunc"
		if got != exp {
			t.Errorf("exp %s, got %s", exp, got)
		}
	})
}

func TestNewFunc(t *testing.T) {
	t.Run("non-func returns error", func(t *testing.T) {
		nofunc := struct{}{}
		f, err := reflectutil.NewFunc(nofunc)
		if err == nil {
			t.Error("exp non-nil error, got nil")
		}
		if !errors.Is(err, reflectutil.ErrNotAFunc) {
			t.Errorf("got unexpected error: %s", err)
		}
		if f != nil {
			t.Errorf("exp nil *Func, got %v", f)
		}
	})

	t.Run("valid func returns *Func", func(t *testing.T) {
		f, err := reflectutil.NewFunc(validFunc)
		if err != nil {
			t.Errorf("got unexpected error: %s", err)
		}
		if got, exp := f.Name, reflectutil.FuncName(validFunc); got != exp {
			t.Errorf("f.Name: exp %s, got %s", exp, got)
		}
	})
}

func TestFunc_Call(t *testing.T) {
	f, _ := reflectutil.NewFunc(validFunc)
	if got, exp := f.Call([]interface{}{42})[0], 42; got != exp {
		t.Errorf("exp %v, got %v", exp, got)
	}
}

func TestFuncSignature_Match(t *testing.T) {
	ftyp := reflect.TypeOf(validFunc)

	t.Run("bad parameters type", func(t *testing.T) {
		sign := reflectutil.FuncSignature{
			In:  []reflect.Kind{reflect.String},
			Out: []reflect.Kind{reflect.Int},
		}
		if match := sign.Match(ftyp); match {
			t.Error("unexpected match for bad parameters type")
		}
	})

	t.Run("bad parameters len", func(t *testing.T) {
		sign := reflectutil.FuncSignature{
			In:  []reflect.Kind{reflect.Int, reflect.Int},
			Out: []reflect.Kind{reflect.Int},
		}
		if match := sign.Match(ftyp); match {
			t.Error("unexpected match for bad parameters len")
		}
	})

	t.Run("bad output type", func(t *testing.T) {
		sign := reflectutil.FuncSignature{
			In:  []reflect.Kind{reflect.Int},
			Out: []reflect.Kind{reflect.String},
		}
		if match := sign.Match(ftyp); match {
			t.Error("unexpected match for bad output type")
		}
	})

	t.Run("bad output len", func(t *testing.T) {
		sign := reflectutil.FuncSignature{
			In:  []reflect.Kind{reflect.Int},
			Out: []reflect.Kind{reflect.Int, reflect.Int},
		}
		if match := sign.Match(ftyp); match {
			t.Error("unexpected match for bad output len")
		}
	})

	t.Run("match signature", func(t *testing.T) {
		sign := reflectutil.FuncSignature{
			In:  []reflect.Kind{reflect.Int},
			Out: []reflect.Kind{reflect.Int},
		}
		if match := sign.Match(ftyp); !match {
			t.Error("exp to match but did not")
		}
	})

	t.Run("match any kind", func(t *testing.T) {
		sign := reflectutil.FuncSignature{
			In:  []reflect.Kind{reflectutil.AnyKind},
			Out: []reflect.Kind{reflectutil.AnyKind},
		}
		if match := sign.Match(ftyp); !match {
			t.Error("exp to match but did not")
		}
	})
}
