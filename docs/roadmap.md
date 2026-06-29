# Roadmap

`go-ruby-bigdecimal/bigdecimal` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's BigDecimal — the
deterministic, interpreter-independent slice extracted from rbgo's internals — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| Arbitrary-precision arithmetic | Add, subtract and multiply at unbounded precision, with the sign and scale rules of Ruby's BigDecimal. | **Done** |
| Division with precision | `div`/`divmod` honouring an explicit precision argument, matching MRI's significant-digit semantics rather than truncating. | **Done** |
| Seven rounding modes | `ROUND_UP`, `ROUND_DOWN`, `ROUND_HALF_UP`, `ROUND_HALF_DOWN`, `ROUND_HALF_EVEN`, `ROUND_CEILING`, `ROUND_FLOOR`, each matching reference Ruby exactly. | **Done** |
| `to_s` "E" and "F" forms | Both the engineering (`"E"`) and plain (`"F"`) string forms, with MRI's grouping and exponent formatting byte-for-byte. | **Done** |
| NaN and infinities | `NaN`, `+Infinity` and `-Infinity` produced and propagated with Ruby's comparison and arithmetic rules. | **Done** |
| Differential oracle & coverage | A wide value corpus exercised here and by the system `ruby`/`bigdecimal`, compared byte-for-byte; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Anything that needs a live binding or evaluation is
  the consumer's job — that is why `rbgo` binds this module rather than the
  reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  behaviour; differences across MRI releases are matched to the reference used by
  the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
