# Usage & API

The public API lives at the module root (`github.com/go-ruby-bigdecimal/bigdecimal`). It is **Ruby-shaped but Go-idiomatic**: `New` / `Add` / `DivPrec` / `Round` / `String` mirror Ruby's `BigDecimal(...)`, `+`, `div`, `round` and `to_s`, while the surface follows Go conventions — value-returning methods, explicit precision and rounding-mode arguments, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-bigdecimal/bigdecimal`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-bigdecimal/bigdecimal
```

## Worked example

```go
a := bigdecimal.New("0.1")
b := bigdecimal.New("0.2")
s := a.Add(b).String()              // "0.3e0"
q := a.DivPrec(b, 10)               // division to 10 significant digits
```

## Shape

```go
// New parses a decimal literal into a BigDecimal (Ruby's BigDecimal("...")).
func New(s string) *BigDecimal

// Add returns d + o at arbitrary precision (Ruby's BigDecimal#+).
func (d *BigDecimal) Add(o *BigDecimal) *BigDecimal

// DivPrec divides d by o to prec significant digits (div with precision).
func (d *BigDecimal) DivPrec(o *BigDecimal, prec int) *BigDecimal

// Round rounds to prec digits using one of the seven rounding modes.
func (d *BigDecimal) Round(prec int, mode RoundMode) *BigDecimal

// String renders the engineering ("E") to_s form.
func (d *BigDecimal) String() string
```

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a wide
value corpus through both the system `ruby` (with `bigdecimal`) and this library
and compares the results **byte-for-byte** — not approximated from memory. The
oracle tests skip themselves where `ruby` is not on `PATH` (e.g. the qemu arch
lanes), so the cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-bigdecimal/bigdecimal` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs
the other way: this library has no dependency on the Ruby runtime.
