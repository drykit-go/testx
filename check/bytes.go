package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type bytesCheck struct {
	passFunc BytesPassFunc
	explFunc ExplainFunc
}

func (c bytesCheck) Pass(got []byte) bool {
	return c.passFunc(got)
}

func (c bytesCheck) Explain(label string, got interface{}) string {
	return c.explFunc(label, got)
}

type bytesValue struct{}

func (bytesValue) Equal(tar []byte) BytesChecker {
	return bytesCheck{
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

// TODO: find a way to handle errors
//
// suggestion:
//
// c := bytesCheck{
//	explainErr: func(label string, err error) string
// }
//
// TODO: ExplainFunc
func (bytesValue) EqualJSON(tar []byte) BytesChecker {
	var decGot, decTar interface{}
	return bytesCheck{
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

func (bytesValue) Len(c IntChecker) BytesChecker {
	return bytesCheck{
		passFunc: func(got []byte) bool {
			return c.Pass(len(got))
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"unexpected %s length: %s",
				label, c.Explain(label, len(got.([]byte))), // TODO: handle assertion error
			)
		},
	}
}
