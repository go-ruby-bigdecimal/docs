# go-ruby-bigdecimal documentation

**Ruby's BigDecimal arbitrary-precision decimals in pure Go — MRI-compatible, no cgo.**

`go-ruby-bigdecimal/bigdecimal` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's BigDecimal,
matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-bigdecimal/bigdecimal`.

It was **extracted from rbgo's prelude/internals into a reusable standalone
library**: the module is standalone and importable by any Go program, and it is
the backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: arithmetic + rounding + formatting complete — MRI byte-exact"
    Faithful port of Ruby's BigDecimal: **arbitrary-precision** add / subtract / multiply, **`div` / `divmod`** honouring an explicit precision argument, the **seven rounding modes** (`ROUND_UP`, `ROUND_DOWN`, `ROUND_HALF_UP`, `ROUND_HALF_DOWN`, `ROUND_HALF_EVEN`, `ROUND_CEILING`, `ROUND_FLOOR`), both the engineering (`"E"`) and plain (`"F"`) `to_s` forms, and `NaN` / `+Infinity` / `-Infinity` handling. Validated by a **differential oracle** against the system `ruby` / `bigdecimal` — values compared byte-for-byte — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
a := bigdecimal.New("0.1")
b := bigdecimal.New("0.2")
s := a.Add(b).String()              // "0.3e0"
q := a.DivPrec(b, 10)               // division to 10 significant digits
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`bigdecimal`](https://github.com/go-ruby-bigdecimal/bigdecimal) | the library — Ruby's BigDecimal in pure Go |
| [`docs`](https://github.com/go-ruby-bigdecimal/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-bigdecimal.github.io`](https://github.com/go-ruby-bigdecimal/go-ruby-bigdecimal.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-bigdecimal/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-bigdecimal/bigdecimal](https://github.com/go-ruby-bigdecimal/bigdecimal).
