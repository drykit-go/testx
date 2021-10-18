package ioutil

import (
	"bytes"
	"fmt"
	"io"
)

// NopRead calls io.ReadAll(*r) and restores r so it can be read again.
func NopRead(r *io.ReadCloser) []byte { // nolint:gocritic // ptrToRefParam
	b, err := io.ReadAll(*r)
	if err != nil {
		panic(fmt.Sprintf("read error: %s", err))
	}
	*r = io.NopCloser(bytes.NewReader(b))
	return b
}
