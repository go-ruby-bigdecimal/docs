# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Library-level benchmark workload (Ruby side) for go-ruby-bigdecimal/bigdecimal.
#
# Mirrors benchmarks/go/main.go exactly: same ~60-digit operands, same Div
# precision, same Newton-sqrt iteration count. Each runtime uses its OWN
# `bigdecimal` stdlib (a C extension under MRI), so every column is the
# Ruby-visible operation done natively. Run with VERIFY=1 to print the CHECK
# lines whose BigDecimal#to_s strings are byte-identical to the Go driver's.
require "bigdecimal"
require_relative "_harness"

A_STR = "1.732050807568877293527446341505872366942805253810380628055806"
B_STR = "2.236067977499789696409173668731276235440618359611525724270897"
PREC  = 80
NEWTON_ITERS = 20

A = BigDecimal(A_STR)
B = BigDecimal(B_STR)
TWO  = BigDecimal("2")
HALF = BigDecimal("2")
DIV_RES = A.div(B, PREC)

# newton_sqrt2 computes √2 by x = (x + 2/x)/2 at PREC significant digits for
# NEWTON_ITERS iterations — the same compound div+add loop as the Go driver.
def newton_sqrt2
  x = BigDecimal("1.5")
  NEWTON_ITERS.times { x = (x + TWO.div(x, PREC)).div(HALF, PREC) }
  x
end

if ENV["VERIFY"] && !ENV["VERIFY"].empty?
  puts "CHECK\tadd\t#{(A + B)}"
  puts "CHECK\tmul\t#{(A * B)}"
  puts "CHECK\tdiv\t#{DIV_RES}"
  puts "CHECK\tnewton-sqrt\t#{newton_sqrt2}"
  puts "CHECK\tparse\t#{BigDecimal(A_STR)}"
  puts "CHECK\tto_s\t#{DIV_RES.to_s}"
  exit
end

bench("add",         2000) { A + B }
bench("mul",         2000) { A * B }
bench("div",          500) { A.div(B, PREC) }
bench("newton-sqrt",   40) { newton_sqrt2 }
bench("parse",       2000) { BigDecimal(A_STR) }
bench("to_s",        2000) { DIV_RES.to_s }
