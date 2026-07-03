<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-bigdecimal` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-bigdecimal`
library** against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby,
TruffleRuby). It measures the **library primitive** through its Go API, isolated
from the `rbgo` interpreter, so the numbers answer: *is the pure-Go
implementation as fast as the reference runtime's own `bigdecimal`?* MRI's
`bigdecimal` is itself a C extension, so this is pure Go vs C — reported honestly.

## Layout

- `go/`               — self-contained Go driver; `go.mod` pins the published library.
- `ruby/bigdecimal.rb`  — the equivalent workload; `ruby/_harness.rb` is the shared timer.
- `run.sh`            — runs every available runtime and prints one Markdown table per
  sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries, and
`GOWORK=off` if a parent `go.work` would otherwise capture the `go/` module.

## Operations

Six representative arbitrary-precision operations on two ~60-significant-digit
operands (decimal expansions of √3 and √5), so a machine-word fast path is never
taken:

- **add** — `a + b`
- **mul** — `a * b`
- **div** — `a.div(b, 80)` (80 significant digits)
- **newton-sqrt** — √2 by `x = (x + 2/x)/2`, 20 iterations at 80 digits
  (a compound div+add precision workload)
- **parse** — `BigDecimal("1.7320…")` / `bigdecimal.New`
- **to_s** — MRI-canonical `BigDecimal#to_s` of the 80-digit `div` result

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script build **identical inputs**, run
the **identical Newton-sqrt algorithm**, and every result's MRI-canonical
`BigDecimal#to_s` string is checked **byte-identical across all runtimes** before
timing (run either driver with `VERIFY=1` to print the `CHECK` lines and `diff`
them). Results are published, dated, in `../docs/performance.md`.
