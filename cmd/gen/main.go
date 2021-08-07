package main

import (
	"errors"
	"flag"
	"log"

	"github.com/drykit-go/testix/check/gen"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	tpl, out, err := readFlags()
	if err != nil {
		return err
	}

	return gen.File(tpl, out)
}

func readFlags() (tpl, out string, err error) {
	flag.StringVar(&tpl, "t", "", "go template source file")
	flag.StringVar(&out, "o", "", "output filename")
	flag.Parse()

	if tpl == "" || out == "" {
		err = errors.New("need template file (-t) and output path (-o)")
	}

	return
}
