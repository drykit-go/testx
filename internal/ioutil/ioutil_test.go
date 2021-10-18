package ioutil_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/drykit-go/testx/internal/ioutil"
	"github.com/drykit-go/testx/internal/testutil"
)

func TestNopRead(t *testing.T) {
	t.Run("can be read again", func(t *testing.T) {
		r := io.NopCloser(bytes.NewReader([]byte("some bytes here")))

		firstRead := ioutil.NopRead(&r)
		secondRead := ioutil.NopRead(&r)

		if !reflect.DeepEqual(firstRead, secondRead) {
			t.Errorf("bad second read: exp %v, got %v", firstRead, secondRead)
		}
	})

	t.Run("panic on read err", func(t *testing.T) {
		r := io.NopCloser(readErrorer{})
		defer testutil.AssertPanic(t, "read error: oh hi Marc")
		ioutil.NopRead(&r)
	})
}

type readErrorer struct{}

func (r readErrorer) Read(b []byte) (int, error) {
	return len(b), errors.New("oh hi Marc")
}
