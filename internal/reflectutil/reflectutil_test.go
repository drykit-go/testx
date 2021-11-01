package reflectutil_test

import (
	"reflect"
	"testing"

	"github.com/drykit-go/testx/internal/reflectutil"
	"github.com/drykit-go/testx/internal/testutil"
)

func TestIsZero(t *testing.T) {
	t.Run("zeros", func(t *testing.T) {
		var zeroMap map[string]any
		var zeroSlice []int
		zeros := []any{
			0, "", 0i + 0, zeroMap, zeroSlice, struct{ n int }{n: 0},
		}

		for _, z := range zeros {
			if !reflectutil.IsZero(z) {
				t.Errorf("exp zero, got false for value: %v", z)
			}
		}
	})

	t.Run("non zeros", func(t *testing.T) {
		nozeros := []any{
			1, "hi", 0i + 1, map[int]bool{}, []float32{}, struct{ n int }{n: -1},
		}

		for _, nz := range nozeros {
			if reflectutil.IsZero(nz) {
				t.Errorf("exp not zero, got true for value: %v", nz)
			}
		}
	})
}

func TestCallUnwrap(t *testing.T) {
	swap := func(x, y float64) (float64, float64) {
		return y, x
	}
	fval := reflect.ValueOf(swap)
	got := reflectutil.CallUnwrap(fval, []any{-1., 1.})
	exp := []any{1., -1.}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("unexpected output: exp %v, got %v", exp, got)
	}
}

func TestMustBeOfKind(t *testing.T) {
	t.Run("bad kind", func(t *testing.T) {
		const kind = reflect.Int8
		badValues := []any{"hi", 1, true, []int8{1}}
		for _, v := range badValues {
			func(v any) {
				defer testutil.AssertPanicMessage(t,
					"expect kind int8, got "+reflect.ValueOf(v).Kind().String(),
				)
				reflectutil.MustBeOfKind(v, kind)
			}(v)
		}
	})

	t.Run("good kind", func(_ *testing.T) {
		for _, tc := range []struct {
			val any
			exp reflect.Kind
		}{
			{val: 42, exp: reflect.Int},
			{val: "hello", exp: reflect.String},
			{val: float32(42), exp: reflect.Float32},
			{val: map[bool][]uint8{}, exp: reflect.Map},
			{val: []byte{68, 65, 108, 108, 111}, exp: reflect.Slice},
		} {
			// no panic means test passes
			reflectutil.MustBeOfKind(tc.val, tc.exp)
		}
	})
}
