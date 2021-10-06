<table align="center">
	<tr>
		<td><a href="../README.md">testx</a></td>
		<td>testx/check</td>
		<td><a href="../checkconv/README.md">testx/checkconv</a></td>
	</tr>
</table>

# testx/check


## Table of contents

- [Checker interfaces](#checker-interfaces)
- [Provided checkers](#provided-checkers)
- [Custom checkers](#custom-checkers)
  - [Extend provided checkers](#extend-provided-checkers)
    - [Using a `New` function](#using-a-new-function)
    - [Using `check.Value.Custom`](#using-checkvaluecustom)
  - [Custom checkers of any type](#custom-checkers-of-any-type)
  - [Contribute!](#contribute)

## Checker interfaces

A checker is a type-specific interface composed of 2 funcs `Pass` and `Explain`:

```go
// IntChecker provides methods to run checks on type int.
type IntChecker interface {
    IntPasser
    Explainer
}

// IntPasser provides a method Pass(got int) bool to determinate
// whether the got int passes a check.
type IntPasser interface { Pass(got int) bool }


// Explainer provides a method Explain to describe the reason of a failed check.
type Explainer interface { Explain(label string, got interface{}) string }
```

## Provided checkers

Package `check` provides a collection of basic checkers for common types and kinds:

<details>
  <summary>Covered types</summary>
  <table>
    <thead>
      <tr>
        <th>Go type or kind</th>
        <th>Checker provider</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>bool</code></td>
        <td><code>check.Bool</code></td>
      </tr>
      <tr>
        <td><code>int</code></td>
        <td><code>check.Int</code></td>
      </tr>
      <tr>
        <td><code>float64</code></td>
        <td><code>check.Float64</code></td>
      </tr>
      <tr>
        <td><code>string</code></td>
        <td><code>check.String</code></td>
      </tr>
      <tr>
        <td><code>[]byte</code></td>
        <td><code>check.Bytes</code></td>
      </tr>
      <tr>
        <td><code>time.Duration</code></td>
        <td><code>check.Duration</code></td>
      </tr>
      <tr>
        <td><code>http.Header</code></td>
        <td><code>check.HTTPHeader</code></td>
      </tr>
      <tr>
        <td><code>slice</code></td>
        <td><code>check.Slice</code></td>
      </tr>
      <tr>
        <td><code>map</code></td>
        <td><code>check.Map</code></td>
      </tr>
      <tr>
        <td><code>struct</code></td>
        <td><code>check.Bool</code></td>
      </tr>
      <tr>
        <td><code>interface{}</code></td>
        <td><code>check.Value</code></td>
      </tr>
    </tbody>
  </table>
</details>

To acknowledge the available checkers for each provider, see their respective interfaces in file [providers.go](./providers.go)

## Custom checkers

This section is about creating custom checkers, either extending the provided ones or full-fledged customs.

### Extend provided checkers

There are 2 ways to extend a provided checker.

#### Using a `New` function

With every `Type` covered by this package comes a `New<Type>Checker` function.

It takes a `PassFunc` for the corresponding type (`func(Type) bool`)
and an `ExplainFunc`, then returns a checker for that type.

```go
func Test42IsEven(t *testing.T) {
    checkIsEven := check.NewIntChecker(
        func(got int) bool { return got&1 == 0 },
        func(label string, got interface{}) string {
            return fmt.Sprintf("%s: expect even int, got %v", label, got)
        },
      )

    testx.
        Value(42).
        Pass(checkconv.FromInt(checkIsEven)).
        Run(t)
}
```

#### Using `check.Value.Custom`

```go
func Test42IsEven(t *testing.T) {
    checkIsEven := check.Value.Custom(
        "even int",
        func(got interface{}) bool {
          gotInt, ok := got.(int)
          ok && got&1 == 0
        },
      )

    testx.
        Value(42).
        Pass(checkIsEven).
        Run(t)
}
```

### Custom checkers of any type

At some point you may need a checker for a type that is not (yet?) covered
by `check`, such as `float32` or `complex128`.  
You may also want custom checkers for your own local types, like `User`.  

Luckily, because _checkers_ are interfaces, nothing prevents you
from creating them.  
In practice, any type that implements `Pass` and `Explain` can be recognized
as a valid checker and used in `testx` runners via `checkconv.Cast`.

For some practical examples, see:

- [example_custom_unknown_test.go](./example_custom_unknown_test.go):
  implementation of a custom checker for a local type `MyType` struct

- [example_custom_closures_test.go](./example_custom_closures_test.go):
  advanced implementation of a parameterized custom checker
  for uncovered type `complex128`, using closures.

### Contribute!

If you believe a particuliar checker should be natively implemented,
feel free to contribute!

- [📋 File an issue](https://github.com/drykit-go/testx/issues/new)
- [🔀 Open a Pull request](https://github.com/drykit-go/testx/pulls)
