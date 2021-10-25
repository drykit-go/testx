package reflectutil_test

import (
	"testing"

	"github.com/drykit-go/testx/internal/reflectutil"
)

func TestIsZero(t *testing.T) {
	t.Run("zeros", func(t *testing.T) {
		var zeroMap map[string]interface{}
		var zeroSlice []int
		zeros := []interface{}{
			0, "", 0i + 0, zeroMap, zeroSlice, struct{ n int }{n: 0},
		}

		for _, z := range zeros {
			if !reflectutil.IsZero(z) {
				t.Errorf("exp zero, got false for value: %v", z)
			}
		}
	})

	t.Run("non zeros", func(t *testing.T) {
		nozeros := []interface{}{
			1, "hi", 0i + 1, map[int]bool{}, []float32{}, struct{ n int }{n: -1},
		}

		for _, nz := range nozeros {
			if reflectutil.IsZero(nz) {
				t.Errorf("exp not zero, got true for value: %v", nz)
			}
		}
	})
}
