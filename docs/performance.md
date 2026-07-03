# Performance

`go-ruby-bigdecimal/bigdecimal` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `bigdecimal`. This
page records a **comparative benchmark** of that module against the reference
Ruby runtimes, part of the ecosystem-wide per-module parity suite.

## What is measured

The **same** Ruby script — a sequence of arbitrary-precision `div`/`mul`/`round` operations — is run under every runtime. `rbgo`'s
number reflects **this pure-Go library doing the work**; every other column is
that interpreter's own `bigdecimal` stdlib. So the comparison is the **Ruby-visible
operation**, apples-to-apples across interpreters. The script prints a
deterministic checksum and its output is checked **byte-identical to MRI**
before timing.

- **Host:** Apple M4 Max, macOS (darwin/arm64). **Method:** best-of-5 wall time
  (best, not mean, to suppress scheduler noise); single-shot processes, no
  warm-up beyond the script's own loop.
- **Runtimes:** `ruby 4.0.5 +PRISM` (MRI, the oracle) and `ruby --yjit`;
  `jruby 10.1.0.0` (OpenJDK 25); `truffleruby 34.0.1` (GraalVM CE Native).
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules)
  (`bigdecimal.rb` + `run.sh`). Reproduce:
  `RBGO=./rbgo TRUFFLE=truffleruby bash bench/modules/run.sh 5`.

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-bigdecimal) | 570 | 1.90× |
| MRI (ruby 4.0.5) | 300 | 1.00× |
| MRI + YJIT | 300 | 1.00× |
| JRuby 10.1.0.0 | 1890 | 6.30× |
| TruffleRuby 34.0.1 | 7640 | 25.47× |

rbgo runs on **go-ruby-bigdecimal** at ~1.9x MRI's BigDecimal (itself a C extension) — a strong showing for pure-Go arbitrary precision. TruffleRuby pays very heavy cold warm-up on this row (7640 ms).

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are real measured numbers from the
    2026-06-29 run — nothing is cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `bigdecimal`?*
MRI's `bigdecimal` is a **C extension**, so this is pure Go measured against C —
reported honestly. The **same workload, same ~60-digit operands, same Div
precision, same Newton-sqrt algorithm and iteration count** run through the Go
library and through each reference runtime's stdlib; every result's
MRI-canonical `BigDecimal#to_s` string was checked **byte-identical across all
five runtimes** before any timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Operations:** `add` (`a+b`), `mul` (`a*b`), `div` (`a.div(b, 80)`),
  `newton-sqrt` (√2 by `x=(x+2/x)/2`, 20 iterations at 80 digits — a compound
  div+add workload), `parse` (`BigDecimal("1.7320…")`) and `to_s` (MRI-canonical
  `to_s` of the 80-digit `div` result). Operands are the 60-significant-digit
  decimal expansions of √3 and √5, so no machine-word fast path is taken.
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

> **Updated 2026-07-03 (allocation-reduction pass).** A hot-path rework — cached
> powers of ten, an allocation-free digit counter, in-place magnitude handling,
> and a single-pass decimal parser (see
> [bigdecimal#1](https://github.com/go-ruby-bigdecimal/bigdecimal/pull/1)) — moved
> the library from *slower than MRI on three ops* to **faster than MRI on all six**
> and **faster than MRI + YJIT on five of six**. The `vs YJIT` column below is the
> pure-Go ns/op divided by the MRI + YJIT ns/op (`< 1.00×` = faster than YJIT).

#### add

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 59.1 | **0.53×** | **0.64×** |
| MRI | 112.5 | 1.00× | — |
| MRI + YJIT | 93.0 | 0.83× | 1.00× |
| JRuby | 292.9 | 2.60× | — |
| TruffleRuby | 3247.9 | 28.87× | — |

#### mul

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 110.1 | **0.44×** | **0.52×** |
| MRI | 249.5 | 1.00× | — |
| MRI + YJIT | 212.5 | 0.85× | 1.00× |
| JRuby | 361.7 | 1.45× | — |
| TruffleRuby | 2379.1 | 9.54× | — |

#### div

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 347.1 | **0.76×** | **0.80×** |
| MRI | 454.0 | 1.00× | — |
| MRI + YJIT | 436.0 | 0.96× | 1.00× |
| JRuby | 884.8 | 1.95× | — |
| TruffleRuby | 4444.1 | 9.79× | — |

#### newton-sqrt

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 12418.8 | **0.85×** | **0.89×** |
| MRI | 14550.0 | 1.00× | — |
| MRI + YJIT | 13950.0 | 0.96× | 1.00× |
| JRuby | 24740.6 | 1.70× | — |
| TruffleRuby | 209974.0 | 14.43× | — |

#### parse

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 156.4 | **0.88×** | 1.06× |
| MRI | 178.5 | 1.00× | — |
| MRI + YJIT | 147.5 | 0.83× | 1.00× |
| JRuby | 437.0 | 2.45× | — |
| TruffleRuby | 3257.6 | 18.25× | — |

#### to_s

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 242.2 | **0.18×** | **0.19×** |
| MRI | 1325.5 | 1.00× | — |
| MRI + YJIT | 1309.0 | 0.99× | 1.00× |
| JRuby | 1217.1 | 0.92× | — |
| TruffleRuby | 3534.2 | 2.67× | — |

**Reading the tranche.** The pure-Go library is now **faster than MRI's C
extension on all six operations** and **faster than MRI + YJIT on five of the
six** — `add` (0.64× YJIT), `mul` (0.52×), `div` (0.80×), `newton-sqrt` (0.89×,
which it inherits from the `div` win it is built on) and `to_s` (0.19×, over 5×
faster than the C `to_s`). The one row that does not clear YJIT is `parse`, and it
still **beats MRI** (0.88×): it lands ~6 % above YJIT because building a
`math/big` significand requires a **decimal→binary conversion** that MRI's
base-10⁹ `BigDecimal` storage never performs, plus a handful of `math/big` heap
allocations against MRI's C arena — an honest representational floor, not a missing
optimization. Closing it would mean abandoning `math/big` for a decimal-limb
representation, which would regress the arithmetic rows that now *beat* YJIT. Every
operation is also **several times faster than cold JRuby and TruffleRuby**.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-bigdecimal/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via `go.mod`),
    the equivalent `ruby/bigdecimal.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh` (add `GOWORK=off` if a parent `go.work` captures the
    module); env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries. Run either driver
    with `VERIFY=1` to print the `CHECK` lines whose `to_s` strings are
    byte-identical across all five runtimes.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly **TruffleRuby**, which pays heavy cold-JIT cost
    here (e.g. `add` at 29×, entirely a warm-up artefact on a ~100 ns op). Rows in
    the sub-microsecond range carry the most relative noise; treat those ratios as
    order-of-magnitude. Every number here is a **real measured value** from the
    dated run above — nothing is fabricated, estimated, or cherry-picked. The
    go-ruby column is the pure-Go library; every other column is that
    interpreter's own `bigdecimal` stdlib (a C extension under MRI) doing the
    equivalent work.
