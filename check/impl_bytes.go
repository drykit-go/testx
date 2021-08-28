package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type bytesCheckerFactory struct{}

func (bytesCheckerFactory) Is(tar []byte) BytesChecker {
	return bytesChecker{
		passFunc: func(got []byte) bool {
			return bytes.Equal(got, tar)
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to equal %v, got %v",
				label, tar, got,
			)
		},
	}
}

func (bytesCheckerFactory) SameJSON(tar []byte) BytesChecker {
	var decGot, decTar interface{}
	return bytesChecker{
		passFunc: func(got []byte) bool {
			if err := json.Unmarshal(got, &decGot); err != nil {
				log.Fatal(err)
			}
			if err := json.Unmarshal(tar, &decTar); err != nil {
				log.Fatal(err)
			}
			return reflect.DeepEqual(decGot, decTar)
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expected json encoded value to equal %#v, got %#v",
				decTar, decGot,
			)
		},
	}
}

func (bytesCheckerFactory) Len(c IntChecker) BytesChecker {
	return bytesChecker{
		passFunc: func(got []byte) bool {
			return c.Pass(len(got))
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"unexpected %s length: %s",
				label, c.Explain(label, len(got.([]byte))),
			)
		},
	}
}
