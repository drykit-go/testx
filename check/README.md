<h1 align="center">testx/check</h1>

<table align="center">
	<tr>
		<td><a href="../README.md">testx</a></td>
		<td>testx/check</td>
		<td><a href="../checkconv/README.md">testx/checkconv</a></td>
	</tr>
</table>

Package `check` provides _checkers_ interfaces and implementations
for `testx` test runners.

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
        <th>Go type</th>
        <th>Checker provider</th>
        <th>Interface</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>bool</code></td>
        <td><code>check.Bool</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#BoolCheckerProvider">
            <code>TypeCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>int</code></td>
        <td><code>check.Int</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#IntCheckerProvider">
            <code>IntCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>float64</code></td>
        <td><code>check.Float64</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#Float64CheckerProvider">
            <code>Float64CheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>string</code></td>
        <td><code>check.String</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#StringCheckerProvider">
            <code>StringCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>[]byte</code></td>
        <td><code>check.Bytes</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#BytesCheckerProvider">
            <code>BytesCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>time.Duration</code></td>
        <td><code>check.Duration</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#DurationCheckerProvider">
            <code>DurationCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>context.Context</code></td>
        <td><code>check.Context</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#ContextCheckerProvider">
            <code>ContextCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>http.Header</code></td>
        <td><code>check.HTTPHeader</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#HTTPHeaderCheckerProvider">
            <code>HTTPHeaderCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>*http.Request</code></td>
        <td><code>check.HTTPRequest</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#HTTPRequestCheckerProvider">
            <code>HTTPRequestCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>*http.Response</code></td>
        <td><code>check.HTTPResponse</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#HTTPResponseCheckerProvider">
            <code>HTTPResponseCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>interface{}</code></td>
        <td><code>check.Value</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#ValueCheckerProvider">
            <code>ValueCheckerProvider</code>
          </a>
        </td>
      </tr>
    </tbody>
    <thead>
      <tr>
        <th>Go kind</th>
        <th>Checker provider</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>slice</code></td>
        <td><code>check.Slice</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#SliceCheckerProvider">
            <code>SliceCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>map</code></td>
        <td><code>check.Map</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#MapCheckerProvider">
            <code>MapCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>struct</code></td>
        <td><code>check.Bool</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#StructCheckerProvider">
            <code>StructCheckerProvider</code>
          </a>
        </td>
      </tr>
    </tbody>
  </table>
</details>

## Custom checkers

This section is about creating custom checkers, either extending the provided ones or full-fledged customs.

### Extend provided checkers

There are 2 ways to extend a provided checker.

#### Using a `New` function

With every `<Type>` covered by this package comes a `New<Type>Checker` function.

It takes a `PassFunc` for the corresponding type (`func(got Type) bool`)
and an `ExplainFunc`, then returns a checker for that type.

```go
func Test42IsEven(t *testing.T) {
    checkIsEven := check.NewIntChecker(
        func(got int) bool { return got&1 == 0 },
        func(label string, got interface{}) string {
            return fmt.Sprintf("%s: expect even int, got %v", label, got)
        },
      )

    testx.Value(42).
        Pass(checkconv.FromInt(checkIsEven)).
        Run(t)
}
```

Related examples:

- [NewIntChecker](https://pkg.go.dev/github.com/drykit-go/testx/check#example-package-NewIntChecker)

#### Using `check.Value.Custom`

```go
func Test42IsEven(t *testing.T) {
    checkIsEven := check.Value.Custom(
        "even int",
        func(got interface{}) bool {
          gotInt, ok := got.(int)
          return ok && got&1 == 0
        },
      )

    testx.Value(42).
        Pass(checkIsEven).
        Run(t)
}
```

### Custom checkers of any type

At some point you may need a checker for a type that is not (yet?) covered
by `check`, such as `float32` or `complex128`.  
You may also want custom checkers for your own local types, like `User`.  

Fortunately, because _checkers_ are interfaces, it is easy to implement
a custom checker for any type.

In practice, any type that implements `func Pass(got Type) bool`
and `Explainer` interface is recognized as a valid checker and
can be used for `testx` runners via `checkconv.Cast`.

Related examples:

- [CustomChecker](https://pkg.go.dev/github.com/drykit-go/testx/check#example-package-CustomChecker):
  implementation of a custom checker that implements `IntChecker`
- [CustomCheckerUnknownType](https://pkg.go.dev/github.com/drykit-go/testx/check#example-package-CustomCheckerUnknownType):
  implementation of a custom checker for a local type `MyType` struct
- [CustomCheckerClosures](https://pkg.go.dev/github.com/drykit-go/testx/check#example-package-CustomCheckerClosures):
  advanced implementation of a parameterized custom checker
  for uncovered type `complex128`, using closures.

### Contribute!

If you believe a particuliar checker should be natively implemented,
feel free to contribute!

- [📋 File an issue](https://github.com/drykit-go/testx/issues/new)
- [🔀 Open a Pull request](https://github.com/drykit-go/testx/pulls)
- [💬 Open a discussion](https://github.com/drykit-go/testx/discussions/new)
