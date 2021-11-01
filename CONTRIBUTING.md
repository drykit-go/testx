# Contributing to testx

This page gathers some information about this repository
that can be valuable to someone about to submit a contribution.

Feel free to [open a discussion 💬](https://github.com/drykit-go/testx/discussions/new)
for any question or suggestion regarding this repository.

## Table of contents

- [External documentation](#external-documentation)
- [Internal documentation](#internal-documentation)
  - [Repository structure](#repository-structure)
  - [Code generation](#code-generation)
  - [Implementing a new checker provider](#implementing-a-new-checker-provider)
- [Contribution suggestions](#contribution-suggestions)
- [Dev environment](#dev-environment)
- [Conventions](#conventions)
  - [Code style](#code-style)
  - [Unit tests](#unit-tests)

## External documentation

The main documentation can help understand how the repository globally works:

- [Go package documentation](https://pkg.go.dev/github.com/drykit-go/testx#section-documentation)
- [Main Readme](./README.md)
- [Package `check` Readme](./check/README.md)
- [Package `checkconv` Readme](./checkconv/README.md)

## Internal documentation

### Repository structure

```sh
.                     # Package testx (test runners)
├── bin               # Binary files (for code generation)
├── check             # Package check (checkers interfaces & providers)
├── checkconv         # Package checkconv (checkers conversion)
├── cmd               # Runnable source files
│   └── gen           # Code generation main command
└── internal          # Cross-packages utilities for internal use only
    ├── fmtexpl       # Formatting helpers for checkers ExplainFuncs
    ├── gen           # Code generation related packages
    │   ├── ...
    │   └── templates # Template files for code generation
    ├── httpconv      # Conversion helpers for net/http types
    ├── ioutil        # Helpers for package io
    └── reflectutil   # Convenience API around package reflect
```


### Code generation

We use code generation to reduce code repetition, and consequently
reduce the complexity of implementations and the risks of errors.

We have 1 use case for that: generate public interfaces of checker providers
from their implementation.

#### Generated files

The following files are generated:

- `check/providers.go`

These files should **never** be manually edited, as specified
in their header (your IDE will likely inform you it shouldn't be edited
if you attempt to).
If an update is needed, re-generated them using the command described below.

#### Command

All generated files are made using a single command:

```sh
make gen
```

This builds `cmd/gen/main.go` and outputs the binary to `bin/gen`
(with your default `GOOS` and `GOARCH` environment variables),
then runs all `//go:generate` directives found in the whole repository.
Located in `gen.go` files, these directives execute `bin/gen` binary
with various arguments.

Run this command each time you implement a new `check.<Type>CheckerProvider`
or update one in a way that changes its public interface, which includes:
  - Adding/removing a method
  - Changing a method signature (a parameter name change counts)
  - Editing a doc comment for a method

### Implementing a new checker provider

To illustrate the process described above, let's implement
`check.Complex128`, a `check.Complex128CheckerProvider`
that performs checks on type `complex128`:

1. Create file `check/providers_complex128.go` and implement
`complex128CheckerProvider` following the existing models.

1. Run `make gen`  
This generates `Complex128CheckerProvider` interface that declares all public
methods implemented

1. Now the methods are accessible from outside, it's time to write unit tests :)

## Contribution suggestions

Here are some contributing suggestions:

- Implement new checkers providers, such as `check.URL`
- Add new methods (checkers) to existing checker providers, such as `check.HTTPRequest.URL`
- Improve test coverage, especially regarding checker providers' `Explain` output

## Dev environment

This branches uses Go 1.18 type parameters, which is yet to be released.
As a consequence it has specific requirements:

- Go 1.18 (we recommand using the wrapper [`gotip`](https://pkg.go.dev/golang.org/dl/gotip))
- `make` commands available (native on Unix-based systems)
- `golangci-lint` to run linters locally ([installation](https://golangci-lint.run/usage/install/#local-installation))
  Note: it currently doesn't work properly with type parameters syntax.

## Conventions

### Code style

Code style conventions are enforced by `golangci-lint`.
Run `make lint` to ensure your code is compliant.

Note: the linter currently doesn't work properly with type parameters syntax.

### Unit tests

We try to maintain a high level of test coverage, so we encourage you
to write relevant test as much as possible when implementing a feature.

Some simple rules apply:
- Location: same directory as the tested file
- Package name: the current package suffixed with `_test`
- File name: the name of the tested file suffixed with `_test`
- Function name: the name of the tested function prefixed with `Test`

To summarize:
```go
// file: mypackage/myfile.go

package mypackage

func MyFunc() {}

// file: mypackage/myfile_test.go

package mypackage_test

func TestMyFunc(t *testing.T) {}
```
