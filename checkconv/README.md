<h1 align="center">testx/checkconv</h1>
<table align="center">
  <tr>
    <td><a href="../README.md">testx</a></td>
    <td><a href="../check/README.md">testx/check</a></td>
    <td>testx/checkconv</td>
  </tr>
</table>

Package `checkconv` provides conversion utilities to convert any typed checker to a `check.ValueChecker`

## Table of contents

- [Why convert checkers](#why-convert-checkers)
- [Available functions](#from-functions)
  - [`From` functions](#from-functions)
  - [`Assert` functions](#assert-functions)
  - [`Cast` functions](#cast-functions)
- [`Assert` vs `Cast`](#assert-vs-cast)
  - [Should I use `Assert` or `Cast`](#should-i-use-assert-or-cast)

## Why convert checkers

Some runners from `testx` have methods that require generic checkers as parameters,
such as `testx.ValueRunner.Pass`:

```go
type ValueRunner interface {
    // Pass adds checkers on the tested value.
    Pass(checkers ...check.ValueChecker) ValueRunner
    // ...
}
```

These methods theorically accept any checker, either from package `check`
or locally implemented with custom types.

However, without a language support for generic types,
package `check` cannot provide a generic checker interface because
they have distinct signatures for method `Pass(got T) bool`.

```go
// What we need:
type GenericChecker[T any] interface {
    Pass(got T) bool
    Explainer
}

// What we have:
type IntChecker interface {
    Pass(got int) bool
    Explainer
}
type StringChecker interface {
    Pass(got string) bool
    Explainer
}
```

As a consequence, there is no way to combine `IntChecker` and `StringChecker`
into a generic `Checker` interface.
For that reason, we use `check.ValueChecker` checker (checker on type `interface{}`),
because it can wrap any checker type.

That's what this package provides: functions to wrap any typed checker
into a generic `check.ValueChecker`.

## Available functions

### `From` functions

`From<Type>` functions return a generic `check.ValueChecker` that wraps
the input `check.<Type>Checker`.

It is typically used for runners that expect a `check.ValueRunner`
when one wants to use a typed checker:

```go
testx.
    Value(42).
    // Pass expects a check.ValueChecker
    Pass(
        // We use check.IntChecker then we convert it using checkconv.FromInt.
        checkconv.FromInt(check.Int.InRange(41, 43)), 
    ).
    Run(t)
)
```

### Assert functions

- `Assert(checker interface{}) check.ValueChecker`
- `AssertMany(checkers ...interface{}) []check.ValueChecker`

Assert functions basically return `From<Type>(inputChecker)`
if that `From<Type>` function exists for the input checker.
Else, it panics.

```go
testx.
    Value(42).
    // Pass expects a check.ValueChecker
    Pass(
        // We use check.IntChecker then we convert it using checkconv.Assert.
        checkconv.Assert(check.Int.InRange(41, 43)),
    ).
    Run(t)
)
```

Alternatively, `AssertMany` can be used to convert several checkers at once:

```go
testx.
    Value(42).
    // Pass expects a check.ValueChecker
    Pass(
        // We use several check.IntChecker then we convert them
        // using checkconv.AssertMany.
        checkconv.AssertMany(
            check.Int.InRange(41, 43),
            check.Int.Not(-1),
            check.Int.GTE(99),
        )...,
    ).
    Run(t)
)
```

### Cast functions

- `Cast(checker interface{}) (check.ValueChecker, bool)`
- `CastMany(checkers ...interface{}) []check.ValueChecker`

Cast functions serve the same purpose as Assert functions:
they wrap the given checker in a `check.ValueChecker`.
The difference is that it works with _any_ type that implement
a checker interface (`Pass(T) bool` and `Explain(string, interface{} string`)
while Assert functions are restricted to the types defined in package `check`.

### `Assert` vs `Cast`

There is a fundamental difference between `Assert` and `Cast` implementations:

- `Assert` uses `From<Type>` functions that rely on type assertion
  to wrap the input checker:
  ```go
  func Assert(knownChecker interface{}) check.ValueChecker {
    switch c := knownChecker.(type) {
    case check.StringChecker:
        return FromString(c)
    // ...
    default:
        // Assert panics if the input checker is not defined in package check.
        log.Panic("assert from unknown checker type")
        return nil
    }
  }

  func FromString(c check.StringChecker) check.ValueChecker {
    return check.NewValueChecker(
      func(got interface{}) bool { return c.Pass(got.(string)) },
      c.Explain,                       // ^^^^^^^^^^^^^^^^^^^ got is safely asserted
    )
  }
  ```
  This is faster, but requires the input checker to implement one of the native
  `check.<Type>Checker` interfaces.

- `Cast` uses **reflection** to determinate whether a checker is valid and call
its methods in the resulting checker. As a consequence it is slower than
`Assert`, however it remains the only way to wrap a `check.ValueChecker`
around a full-fledged custom checker that performs checks on types that are not 
defined in package `check`.

#### Should I use `Assert` or `Cast`

- Use `Assert`/`From<Type>` to convert any checker from package `check`,
  or custom checkers that implement any `check.<Type>Checker` interface
  (see [provided checkers](../check/README.md#provided-checkers)).

  ```go
  checkconv.AssertMany(
      check.Int.Range(41, 43), // satisfies check.IntChecker
      check.NewIntChecker(isEven, explainIsEven), // satisfies check.IntChecker
      myCustomIntChecker, // satisfies check.IntChecker
  )
  ```

- Use `Cast` to convert any checker on types that are not defined
  in package `check`.

  ```go
  checkconv.CastMany(
      myCustomComplex128Checker,
      myCustomUserChecker,
  )
  ```