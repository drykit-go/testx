package main

import (
	"errors"
	"flag"
	"log"

	"github.com/drykit-go/testx/internal/gen"
)

var (
	tpl  = flag.String("t", "", "go template source file")
	out  = flag.String("o", "", "output path")
	kind = flag.String("k", "", "data to be generated (types or interfaces)")
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := parseFlags(); err != nil {
		return err
	}

	var generate func(string, string) error
	switch *kind {
	case "interfaces":
		generate = gen.Interfaces
	case "types":
		generate = gen.Types
	default:
		return errors.New("invalid value for -k: expect 'types' or 'interfaces'")
	}
	return generate(*tpl, *out)
}

func parseFlags() error {
	flag.Parse()
	if *kind == "" {
		return errors.New("missing data kind (-k)")
	}
	if *tpl == "" {
		return errors.New("missing template source file (-t)")
	}
	if *out == "" {
		return errors.New("missing output path (-o)")
	}
	return nil
}
