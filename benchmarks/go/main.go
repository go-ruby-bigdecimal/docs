// SPDX-License-Identifier: BSD-3-Clause
//
// Library-level benchmark workload (Go side) for go-ruby-bigdecimal/bigdecimal.
//
// The same arbitrary-precision operations are run here through the pure-Go
// library's Go API and, in ruby/bigdecimal.rb, through each reference runtime's
// own `bigdecimal` stdlib. Inputs, iteration counts and the Newton-sqrt
// algorithm are identical on both sides, and every result's MRI-canonical
// BigDecimal#to_s string is byte-checked equal across runtimes before timing
// (run with VERIFY=1 to print the CHECK lines both drivers emit).
package main

import (
	"fmt"
	"os"

	bd "github.com/go-ruby-bigdecimal/bigdecimal"
)

// Two ~60-significant-digit operands (decimal expansions of √3 and √5) exercise
// real arbitrary precision rather than a machine-word fast path.
const (
	aStr = "1.732050807568877293527446341505872366942805253810380628055806"
	bStr = "2.236067977499789696409173668731276235440618359611525724270897"
	// prec is the significant-digit budget handed to Div, matching Ruby's
	// a.div(b, prec) second argument.
	prec = 80
	// newtonIters fixes the Newton–Raphson sqrt loop length so both sides do
	// exactly the same compound div+add work.
	newtonIters = 20
)

func mustNew(s string) *bd.Decimal {
	d, err := bd.New(s)
	if err != nil {
		panic(err)
	}
	return d
}

// newtonSqrt2 computes √2 by x_{n+1} = (x_n + 2/x_n)/2 at `prec` significant
// digits for `newtonIters` iterations — a div-and-add heavy precision workload.
func newtonSqrt2() *bd.Decimal {
	two := mustNew("2")
	half := mustNew("2")
	x := mustNew("1.5")
	for i := 0; i < newtonIters; i++ {
		x = x.Add(two.Div(x, prec)).Div(half, prec)
	}
	return x
}

func main() {
	a := mustNew(aStr)
	b := mustNew(bStr)
	divRes := a.Div(b, prec)

	if os.Getenv("VERIFY") != "" {
		fmt.Printf("CHECK\tadd\t%s\n", a.Add(b).String())
		fmt.Printf("CHECK\tmul\t%s\n", a.Mul(b).String())
		fmt.Printf("CHECK\tdiv\t%s\n", divRes.String())
		fmt.Printf("CHECK\tnewton-sqrt\t%s\n", newtonSqrt2().String())
		fmt.Printf("CHECK\tparse\t%s\n", mustNew(aStr).String())
		fmt.Printf("CHECK\tto_s\t%s\n", divRes.String())
		return
	}

	bench("add", 2000, func() { sink = a.Add(b) })
	bench("mul", 2000, func() { sink = a.Mul(b) })
	bench("div", 500, func() { sink = a.Div(b, prec) })
	bench("newton-sqrt", 40, func() { sink = newtonSqrt2() })
	bench("parse", 2000, func() { sink = mustNew(aStr) })
	bench("to_s", 2000, func() { sink = divRes.String() })
}
