package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/drykit-go/testx/internal/gen"
)

const (
	// binDirPath is the path from testx root to bin directory.
	binDirPath = "/bin"
	// tplDirPath is the relative path from bin/gen to internal/gen/templates.
	tplDirPath = "../internal/gen/templates"
	// tplExt is the extension for all template files
	tplExt = "gotmpl"
	// outExt is the output file extension
	outExt = "go"
)

var (
	name = flag.String("name", "", "template name in internal/gen (without extension)")
	kind = flag.String("kind", "", "data to be generated (types or interfaces)")
)

var kindsFuncs = map[string]func(tpl, out string) error{
	"interfaces": gen.Interfaces,
	"types":      gen.Types,
}

func main() {
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("⏳ Generating %s %s...\n", *name, *kind)
	if err := run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Done")
}

func run() error {
	tpl, out, err := getFilesPaths()
	if err != nil {
		return err
	}

	generate, ok := kindsFuncs[*kind]
	if !ok {
		return fmt.Errorf("unknown gen kind: %s", *kind)
	}

	return generate(tpl, out)
}

func parseFlags() error {
	flag.Parse()
	if *kind == "" {
		return errors.New("missing generation kind (-kind)")
	}
	if *name == "" {
		return errors.New("missing template name (-name)")
	}
	return nil
}

func getFilesPaths() (tplPath, outPath string, err error) {
	workDir, err := os.Getwd()
	if err != nil {
		return
	}
	currExe, err := os.Executable()
	if err != nil {
		return
	}
	currDir := filepath.Dir(currExe)

	if !strings.HasSuffix(currDir, binDirPath) {
		// when run with go run, currDir is an unexploitable temp dir,
		// for that reason we ensure it is run from the executable.
		return "", "", errors.New("must be run from executable testx/bin/gen")
	}

	tplPath = filepath.Join(currDir, tplDirPath, filename(*name, tplExt))
	outPath = filepath.Join(workDir, filename(*name, outExt))
	return
}

func filename(name, ext string) string {
	return fmt.Sprintf("%s.%s", name, ext)
}
