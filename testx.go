package testx

import (
	"testing"
)

type gotType interface{}

type getFunc func() gotType

func fail(t *testing.T, msg string) {
	t.Error(msg)
}
