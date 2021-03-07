// Package tomath wraps github.com/shopspring/decimal library which implements an arbitrary precision fixed-point decimal.
//
// The zero-value of a Decimal is 0, as you would expect.
// The zero-value name is "?". To set a name for a decimal without a name use the SetName() method.
//
// The best way to create a new Decimal is to use decimal.NewFromStringWithName, ex:
//
//     n, err := decimal.NewFromStringWithName("var1", "1.3")
//     n.String() // output: "1.3"
//
//     vars, formula := n.Add("var2", "1.8")
//                       .SetName("var3")
//                       .Math()
//     // vars:    "var1 + var2 = var3"
//     // formula: "1.3 + 1.8 = 3.1"
//
// To use Decimal as part of a struct:
//
//     type Struct struct {
//         Number Decimal
//     }
//
// Note: This can "only" represent numbers with a maximum of 2^31 digits after the decimal point.
package tomath

import (
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	leftParen  = "("
	rightParen = ")"

	// 	abs        = "abs"
	// 	add        = " + "
	// 	sub        = " - "
	// 	neg        = "neg"
	// 	mul        = " * "
	// 	shift      = "shift"
	// 	div        = " / "
	// 	quoRem     = "quoRem"
	// 	divRound   = "divRound"
	// 	mod        = " % "
	// 	pow        = "^"
	// 	round      = "round"
	// 	roundBank  = "roundBank"
	// 	roundCash  = "roundCash"
	// 	floor      = "floor"
	// 	ceil       = "ceil"
	// 	truncate   = "truncate"
	// 	min        = "min"
	// 	comma      = ", "
	// 	max        = "max"
	// 	sum        = "sum"
	// 	avg        = "avg"
	// 	atan       = "atan"
	// 	sin        = "sin"
	// 	cos        = "cos"
	// 	tan        = "tan"
	equal = " = "
)

var symbols = map[byte]string{
	abs:       "abs",
	neg:       "neg",
	round:     "round",
	roundBank: "roundBank",
	roundCash: "roundCash",
	floor:     "floor",
	ceil:      "ceil",
	truncate:  "truncate",
	atan:      "atan",
	sin:       "sin",
	cos:       "cos",
	tan:       "tan",
	add:       " + ", // second the binary operations
	sub:       " - ",
	mul:       " * ",
	div:       " / ",
	shift:     "shift",
	quoRem:    "quoRem",
	divRound:  "divRound",
	mod:       " % ",
	pow:       "^",
	min:       "min", //finally the operations that can take many operands
	max:       "max",
	sum:       "sum",
	avg:       "avg",
}

var (
	// unary operators with precision
	round     byte = 0
	roundBank byte = 1
	roundCash byte = 2
	shift     byte = 3
	truncate  byte = 4

	// unary operations
	abs   byte = 5
	atan  byte = 6
	ceil  byte = 7
	cos   byte = 8
	floor byte = 9
	neg   byte = 10
	sin   byte = 11
	tan   byte = 12

	// binary operators with precision
	divRound byte = 13
	quoRem   byte = 14

	//  binary operations
	add byte = 15
	div byte = 16
	mod byte = 17
	mul byte = 18
	pow byte = 19
	sub byte = 20

	// variatic operators
	avg byte = 21
	max byte = 22
	min byte = 23
	sum byte = 24
)

func isUnary(b byte) bool {
	return b < divRound
}

func isUnaryWithPrecision(b byte) bool {
	return b < abs
}

func isBinary(b byte) bool {
	return b < avg && b > tan
}

func isBinaryWithPrecision(b byte) bool {
	return b < add && b > tan
}

func isVariatic(b byte) bool {
	return b > sub
}

type (
	// Decimal represents a fixed-point decimal. It is immutable.
	// // number = value * 10 ^ exp
	// Decimal struct {
	// 	parens  bool
	// 	name    string
	// 	decimal decimal.Decimal
	// 	vars    string
	// 	formula string
	// }

	// NullDecimal represents a nullable decimal with compatibility for
	// scanning null values from the database.
	NullDecimal struct {
		name    string
		decimal decimal.NullDecimal
	}

	// Potentially store strings in string table to cut out duplicates
	// parens can be bit shiftend if you know how many uniary operators have gone by
	// Decimal struct {
	// 	ops        []byte            // the operation decides what of the others is to be used
	// 	parens     []bool            // not needed for unary operations
	// 	names      []string          // not needed for unary operations
	// 	decimals   []decimal.Decimal // not needed for unary operations
	// 	valCounts  []uint8           // only needed for operations that can take many operands
	// 	precisions []int32
	// }

	// TODO: rename Decimal
	Decimal struct {
		left      *Decimal
		op        *byte
		right     *Decimal
		precision *int32
		name      *string
		value     *decimal.Decimal
	}
)

// var Zero = Decimal{decimals: []decimal.Decimal{decimal.Zero}, names: []string{"zero"}}

// SetName sets the name of the Decimal
// func (d Decimal) SetName(name string) Decimal {
// 	d.name = name
// 	if d.vars == "" {
// 		d.vars = name
// 	}
// 	return d
// }

// // SetName gets the name of the Decimal
// func (d Decimal) GetName() string {
// 	return d.name
// }

// // Resolve removes the underlying math from the decimal and replaces it with the
// // current name and value.
// func (d Decimal) Resolve() Decimal {
// 	return Decimal{
// 		name:    d.name,
// 		vars:    d.name,
// 		formula: d.String(),
// 		decimal: d.decimal,
// 	}
// }

// ResolveTo is a wrapper around SetName() and Resolve().
// func (d Decimal) ResolveTo(name string) Decimal {
// 	return d.SetName(name).Resolve()
// }

// // Decimal ejects the github.com/shopspring/decimal#Decimal
// func (d Decimal) Decimal() decimal.Decimal {
// 	return d.decimal
// }

// Math returns two strings representing the formula underlying the decimal. The
// first uses the decimal names. The second uses the decimal values. Both are
// follwed by an equals sign with the current name and value respectively.
func (d Decimal) Math() (string, string) {
	var vars, formula strings.Builder

	// handle single value without ops first
	if d.op == nil {
		if d.name == nil {
			vars.WriteRune('?')
		} else {
			vars.WriteString(*d.name)
		}

		value := d.value.String()
		formula.WriteString(value)

		vars.WriteString(equal)

		if d.name == nil {
			vars.WriteRune('?')
		} else {
			vars.WriteString(*d.name)
		}

		formula.WriteString(equal)
		formula.WriteString(value)

		return vars.String(), formula.String()
	}

	curDecimal := &d
	var parents []*Decimal
	visited := make(map[*Decimal]bool)
	values := make([]*decimal.Decimal, 0, 1) // we should have at least 1 value

	for curDecimal != nil {
		if curDecimal.op == nil {
			writeValue(&vars, &formula, curDecimal)
			values = append(values, curDecimal.value)
			visited[curDecimal] = true

			curDecimal = parents[len(parents)-1]
			parents = parents[:len(parents)-1]
			continue
		} else if !visited[curDecimal] && isUnary(*curDecimal.op) {
			write(&vars, &formula, symbols[*curDecimal.op])
			write(&vars, &formula, leftParen)

			visited[curDecimal] = true
		}

		if curDecimal.left != nil && !visited[curDecimal.left] {
			parents = append(parents, curDecimal)
			curDecimal = curDecimal.left
			continue
		}

		if isUnary(*curDecimal.op) {
			write(&vars, &formula, rightParen)
		} else if !visited[curDecimal] && isBinary(*curDecimal.op) {
			write(&vars, &formula, symbols[*curDecimal.op])
			visited[curDecimal] = true
		}

		if curDecimal.right != nil && !visited[curDecimal.right] {
			parents = append(parents, curDecimal)
			curDecimal = curDecimal.right
			continue
		}

		switch *curDecimal.op {
		case abs:
			val := values[len(values)-1]
			values = values[:len(values)-1]

			result := val.Abs()

			values = append(values, &result)
		case add:
			val2 := values[len(values)-1]
			val1 := values[len(values)-2]

			values = values[:len(values)-2]

			result := val1.Add(*val2)

			values = append(values, &result)
		}

		if len(parents) > 0 {
			curDecimal = parents[len(parents)-1]
			parents = parents[:len(parents)-1]
		} else {
			curDecimal = nil
		}
	}

	vars.WriteString(equal)
	// TODO: implement final name
	vars.WriteRune('?')

	formula.WriteString(equal)
	formula.WriteString(values[0].String())

	return vars.String(), formula.String()
}

func writeValue(vars, formula *strings.Builder, d *Decimal) {
	if d.name != nil && *d.name != "" {
		vars.WriteString(*d.name)
	} else {
		vars.WriteRune('?')
	}

	formula.WriteString(d.value.String())
}

func write(vars, formula *strings.Builder, s string) {
	vars.WriteString(s)
	formula.WriteString(s)
}

// New returns a new fixed-point decimal, value * 10 ^ exp.
func New(value int64, exp int32) Decimal {
	d := decimal.New(value, exp)
	return Decimal{value: &d}
}

// NewWithName returns a new fixed-point decimal, value * 10 ^ exp with a given name.
func NewWithName(name string, value int64, exp int32) Decimal {
	d := decimal.New(value, exp)
	return Decimal{name: &name, value: &d}
}

// NewFromInt converts a int64 to Decimal.
//
// Example:
//
//     NewFromInt(123).String() // output: "123"
//     NewFromInt(-10).String() // output: "-10"
func NewFromInt(value int64) Decimal {
	d := decimal.NewFromInt(value)
	return Decimal{value: &d}
}

// NewFromIntWithName converts a int64 to Decimal with a given name.
//
// Example:
//
//     NewFromIntWithName("var1", 123).String() // output: "123"
//     NewFromIntWithName("var1", -10).String() // output: "-10"
func NewFromIntWithName(name string, value int64) Decimal {
	d := decimal.NewFromInt(value)
	return Decimal{name: &name, value: &d}
}

// NewFromInt32 converts a int32 to Decimal.
//
// Example:
//
//     NewFromInt(123).String() // output: "123"
//     NewFromInt(-10).String() // output: "-10"
func NewFromInt32(value int32) Decimal {
	d := decimal.NewFromInt32(value)
	return Decimal{value: &d}
}

// NewFromInt32WithName converts a int32 to Decimal with a given name.
//
// Example:
//
//     NewFromInt32WithName("var1", 123).String() // output: "123"
//     NewFromInt32WithName("var1", -10).String() // output: "-10"
func NewFromInt32WithName(name string, value int32) Decimal {
	d := decimal.NewFromInt32(value)
	return Decimal{name: &name, value: &d}
}

// NewFromBigInt returns a new Decimal from a big.Int, value * 10 ^ exp
func NewFromBigInt(value *big.Int, exp int32) Decimal {
	d := decimal.NewFromBigInt(value, exp)
	return Decimal{value: &d}
}

// NewFromBigIntWithName returns a new Decimal from a big.Int, value * 10 ^ exp
// with a given name
func NewFromBigIntWithName(name string, value *big.Int, exp int32) Decimal {
	d := decimal.NewFromBigInt(value, exp)
	return Decimal{name: &name, value: &d}
}

// NewFromString returns a new Decimal from a string representation.
// Trailing zeroes are not trimmed.
//
// Example:
//
//     d, err := NewFromString("-123.45")
//     d2, err := NewFromString(".0001")
//     d3, err := NewFromString("1.47000")
//
func NewFromString(value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{value: &d}, nil
}

// NewFromStringWithName returns a new Decimal from a string representation with
// a given name. Trailing zeroes are not trimmed.
//
// Example:
//
//     d, err := NewFromStringWithName("var1", "-123.45")
//     d2, err := NewFromStringWithName("var1", ".0001")
//     d3, err := NewFromStringWithName("var1", "1.47000")
//
func NewFromStringWithName(name string, value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{name: &name, value: &d}, nil
}

// RequireFromString returns a new Decimal from a string representation
// or panics if NewFromString would have returned an error.
//
// Example:
//
//     d := RequireFromString("-123.45")
//     d2 := RequireFromString(".0001")
//
func RequireFromString(value string) Decimal {
	d := decimal.RequireFromString(value)
	return Decimal{value: &d}
}

// RequireFromStringWithName returns a new Decimal from a string representation
// with a given name or panics if NewFromString would have returned an error.
//
// Example:
//
//     d := RequireFromStringWithName("var1", "-123.45")
//     d2 := RequireFromStringWithName("var1", ".0001")
//
func RequireFromStringWithName(name string, value string) Decimal {
	d := decimal.RequireFromString(value)
	return Decimal{name: &name, value: &d}
}

// NewFromFloat converts a float64 to Decimal.
//
// The converted number will contain the number of significant digits that can be
// represented in a float with reliable roundtrip.
// This is typically 15 digits, but may be more in some cases.
// See https://www.exploringbinary.com/decimal-precision-of-binary-floating-point-numbers/ for more information.
//
// For slightly faster conversion, use NewFromFloatWithExponent where you can specify the precision in absolute terms.
//
// NOTE: this will panic on NaN, +/-inf
func NewFromFloat(value float64) Decimal {
	d := decimal.NewFromFloat(value)
	return Decimal{value: &d}
}

// NewFromFloatWithName converts a float64 to Decimal with a given name.
//
// The converted number will contain the number of significant digits that can be
// represented in a float with reliable roundtrip.
// This is typically 15 digits, but may be more in some cases.
// See https://www.exploringbinary.com/decimal-precision-of-binary-floating-point-numbers/ for more information.
//
// For slightly faster conversion, use NewFromFloatWithNameWithExponent where you can specify the precision in absolute terms.
//
// NOTE: this will panic on NaN, +/-inf
func NewFromFloatWithName(name string, value float64) Decimal {
	d := decimal.NewFromFloat(value)
	return Decimal{name: &name, value: &d}
}

// NewFromFloat32 converts a float32 to Decimal.
//
// The converted number will contain the number of significant digits that can be
// represented in a float with reliable roundtrip.
// This is typically 6-8 digits depending on the input.
// See https://www.exploringbinary.com/decimal-precision-of-binary-floating-point-numbers/ for more information.
//
// For slightly faster conversion, use NewFromFloatWithExponent where you can specify the precision in absolute terms.
//
// NOTE: this will panic on NaN, +/-inf
func NewFromFloat32(value float32) Decimal {
	d := decimal.NewFromFloat32(value)
	return Decimal{value: &d}
}

// NewFromFloat32WithName converts a float32 to Decimal with a given name.
//
// The converted number will contain the number of significant digits that can be
// represented in a float with reliable roundtrip.
// This is typically 6-8 digits depending on the input.
// See https://www.exploringbinary.com/decimal-precision-of-binary-floating-point-numbers/ for more information.
//
// For slightly faster conversion, use NewFromFloatWithExponent where you can specify the precision in absolute terms.
//
// NOTE: this will panic on NaN, +/-inf
func NewFromFloat32WithName(name string, value float32) Decimal {
	d := decimal.NewFromFloat32(value)
	return Decimal{name: &name, value: &d}
}

// NewFromFloatWithExponent converts a float64 to Decimal, with an arbitrary
// number of fractional digits.
//
// Example:
//
//     NewFromFloatWithExponent(123.456, -2).String() // output: "123.46"
//
func NewFromFloatWithExponent(value float64, exp int32) Decimal {
	d := decimal.NewFromFloatWithExponent(value, exp)
	return Decimal{value: &d}
}

// NewFromFloatWithExponentWithName converts a float64 to Decimal with a given name, with an arbitrary
// number of fractional digits.
//
// Example:
//
//     NewFromFloatWithExponentWithName("var1", 123.456, -2).String() // output: "123.46"
//
func NewFromFloatWithExponentWithName(name string, value float64, exp int32) Decimal {
	d := decimal.NewFromFloatWithExponent(value, exp)
	return Decimal{name: &name, value: &d}
}

// NewFromDecimal returns a new Decimal from github.com/shopspring/decimal#Decimal.
func NewFromDecimal(d decimal.Decimal) Decimal {
	return Decimal{value: &d}
}

// NewFromDecimalWithName returns a new Decimal from github.com/shopspring/decimal#Decimal
// with a given name.
func NewFromDecimalWithName(name string, d decimal.Decimal) Decimal {
	return Decimal{name: &name, value: &d}
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	return Decimal{
		op:   &abs,
		left: &d,
	}
}

// Add returns d + d2.
func (d Decimal) Add(d2 Decimal) Decimal {
	return Decimal{
		op:    &add,
		left:  &d,
		right: &d2,
	}
}

// Sub returns d - d2.
// func (d Decimal) Sub(d2 Decimal) Decimal {
// 	return Decimal{
// 		ops:        append(append(d.ops, d2.ops...), sub),
// 		decimals:   append(d.decimals, d2.decimals...),
// 		names:      append(append(d.names, d2.names...)),
// 		parens:     append(d.parens, true),
// 		precisions: d.precisions,
// 	}
// }

// Neg returns -d.
// func (d Decimal) Neg() Decimal {
// 	return Decimal{
// 		ops:        append(d.ops, neg),
// 		decimals:   d.decimals,
// 		names:      d.names,
// 		parens:     d.parens,
// 		precisions: d.precisions,
// 	}
// }

// Mul returns d * d2.
// func (d Decimal) Mul(d2 Decimal) Decimal {
// 	return Decimal{
// 		ops:        append(append(d.ops, d2.ops...), mul),
// 		decimals:   append(d.decimals, d2.decimals...),
// 		names:      append(append(d.names, d2.names...)),
// 		parens:     append(d.parens, true),
// 		precisions: d.precisions,
// 	}

// 	// dec := Decimal{decimal: d.decimal.Mul(d2.decimal)}
// 	// var vars, formula string

// 	// if d.parens {
// 	// 	vars += leftParen + d.vars + rightParen + mul
// 	// 	formula += leftParen + d.formula + rightParen + mul
// 	// } else {
// 	// 	vars += d.vars + mul
// 	// 	formula += d.formula + mul
// 	// }

// 	// if d2.parens {
// 	// 	vars += leftParen + d2.vars + rightParen
// 	// 	formula += leftParen + d2.formula + rightParen
// 	// } else {
// 	// 	vars += d2.vars
// 	// 	formula += d2.formula
// 	// }

// 	// dec.vars = vars
// 	// dec.formula = formula

// 	// return dec
// }

// Shift shifts the decimal in base 10.
// It shifts left when shift is positive and right if shift is negative.
// In simpler terms, the given value for shift is added to the exponent
// of the decimal.
// func (d Decimal) Shift(s int32) Decimal {
// 	return Decimal{
// 		ops:        append(d.ops, shift),
// 		decimals:   d.decimals,
// 		names:      d.names,
// 		precisions: append(d.precisions, s),
// 	}
// }

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
// func (d Decimal) Div(d2 Decimal) Decimal {
// 	return Decimal{
// 		ops:        append(append(d.ops, d2.ops...), div),
// 		decimals:   append(d.decimals, d2.decimals...),
// 		names:      append(append(d.names, d2.names...)),
// 		parens:     append(d.parens, true),
// 		precisions: d.precisions,
// 	}
// 	// dec := Decimal{decimal: d.decimal.Div(d2.decimal)}

// 	// var vars, formula string
// 	// if d.parens {
// 	// 	vars += leftParen + d.vars + rightParen + div
// 	// 	formula += leftParen + d.formula + rightParen + div
// 	// } else {
// 	// 	vars += d.vars + div
// 	// 	formula += d.formula + div
// 	// }

// 	// if d2.parens {
// 	// 	vars += leftParen + d2.vars + rightParen
// 	// 	formula += leftParen + d2.formula + rightParen
// 	// } else {
// 	// 	vars += d2.vars
// 	// 	formula += d2.formula
// 	// }

// 	// dec.vars = vars
// 	// dec.formula = formula

// 	// return dec
// }

// QuoRem does divsion with remainder
// d.QuoRem(d2,precision) returns quotient q and remainder r such that
//   d = d2 * q + r, q an integer multiple of 10^(-precision)
//   0 <= r < abs(d2) * 10 ^(-precision) if d>=0
//   0 >= r > -abs(d2) * 10 ^(-precision) if d<0
// Note that precision<0 is allowed as input.
// func (d Decimal) QuoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
// 	return Decimal{
// 			ops:      append(append(d.ops, d2.ops...), quoRem),
// 			decimals: append(d.decimals, d2.decimals...),
// 			names: append(
// 				append(d.names, d2.names...),
// 				d.names[len(d.names)-1]+d2.names[len(d2.names)-1]+"Quotient"),
// 			parens:     append(d.parens, true),
// 			precisions: append(d.precisions, precision),
// 		}, Decimal{
// 			ops:      append(append(d.ops, d2.ops...), quoRem),
// 			decimals: append(d.decimals, d2.decimals...),
// 			names: append(
// 				append(d.names, d2.names...),
// 				d.names[len(d.names)-1]+d2.names[len(d2.names)-1]+"Remainder"),
// 			parens:     append(d.parens, true),
// 			precisions: append(d.precisions, precision),
// 		}

// 	// d3, d4 := d.decimal.QuoRem(d2.decimal, precision)
// 	// p := strconv.Itoa(int(precision))

// 	// var vars, formula string
// 	// if d.parens {
// 	// 	vars += quoRem + leftParen + p + rightParen + leftParen + leftParen + d.vars + rightParen + div
// 	// 	formula += quoRem + leftParen + p + rightParen + leftParen + leftParen + d.formula + rightParen + div
// 	// } else {
// 	// 	vars += quoRem + leftParen + p + rightParen + leftParen + d.vars + div
// 	// 	formula += quoRem + leftParen + p + rightParen + leftParen + d.formula + div
// 	// }

// 	// if d2.parens {
// 	// 	vars += leftParen + d2.vars + rightParen + rightParen
// 	// 	formula += leftParen + d2.formula + rightParen + rightParen
// 	// } else {
// 	// 	vars += d2.vars + rightParen
// 	// 	formula += d2.formula + rightParen
// 	// }

// 	// return Decimal{name: d.name + d2.name + "Quotient", decimal: d3, vars: vars, formula: formula},
// 	// 	Decimal{name: d.name + d2.name + "Remainder", decimal: d4, vars: vars, formula: formula}
// }

// // DivRound divides and rounds to a given precision
// // i.e. to an integer multiple of 10^(-precision)
// //   for a positive quotient digit 5 is rounded up, away from 0
// //   if the quotient is negative then digit 5 is rounded down, away from 0
// // Note that precision<0 is allowed as input.
// func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
// 	dec := Decimal{decimal: d.decimal.DivRound(d2.decimal, precision)}
// 	p := strconv.Itoa(int(precision))

// 	var vars, formula string
// 	if d.parens {
// 		vars += divRound + leftParen + p + rightParen + leftParen + leftParen + d.vars + rightParen + div
// 		formula += divRound + leftParen + p + rightParen + leftParen + leftParen + d.formula + rightParen + div
// 	} else {
// 		vars += divRound + leftParen + p + rightParen + leftParen + d.vars + div
// 		formula += divRound + leftParen + p + rightParen + leftParen + d.formula + div
// 	}

// 	if d2.parens {
// 		vars += leftParen + d2.vars + rightParen + rightParen
// 		formula += leftParen + d2.formula + rightParen + rightParen
// 	} else {
// 		vars += d2.vars + rightParen
// 		formula += d2.formula + rightParen
// 	}

// 	dec.vars = vars
// 	dec.formula = formula

// 	return dec
// }

// // Mod returns d % d2.
// func (d Decimal) Mod(d2 Decimal) Decimal {
// 	dec := Decimal{decimal: d.decimal.Mod(d2.decimal)}

// 	var vars, formula string
// 	if d.parens {
// 		vars += leftParen + d.vars + rightParen + mod
// 		formula += leftParen + d.formula + rightParen + mod
// 	} else {
// 		vars += d.vars + mod
// 		formula += d.formula + mod
// 	}

// 	if d2.parens {
// 		vars += leftParen + d2.vars + rightParen
// 		formula += leftParen + d2.formula + rightParen
// 	} else {
// 		vars += d2.vars
// 		formula += d2.formula
// 	}

// 	dec.vars = vars
// 	dec.formula = formula

// 	return dec
// }

// // Pow returns d to the power d2
// func (d Decimal) Pow(d2 Decimal) Decimal {
// 	dec := Decimal{decimal: d.decimal.Pow(d2.decimal)}

// 	var vars, formula string
// 	if d.parens {
// 		vars += leftParen + d.vars + rightParen + pow
// 		formula += leftParen + d.formula + rightParen + pow
// 	} else {
// 		vars += d.vars + pow
// 		formula += d.formula + pow
// 	}

// 	if d2.parens {
// 		vars += leftParen + d2.vars + rightParen
// 		formula += leftParen + d2.formula + rightParen
// 	} else {
// 		vars += d2.vars
// 		formula += d2.formula
// 	}

// 	dec.vars = vars
// 	dec.formula = formula

// 	return dec
// }

// // Cmp compares the numbers represented by d and d2 and returns:
// //
// //     -1 if d <  d2
// //      0 if d == d2
// //     +1 if d >  d2
// //
// func (d Decimal) Cmp(d2 Decimal) int {
// 	return d.decimal.Cmp(d2.decimal)
// }

// // Equal returns whether the numbers represented by d and d2 are equal.
// func (d Decimal) Equal(d2 Decimal) bool {
// 	return d.decimal.Equal(d2.decimal)
// }

// // Equals is deprecated, please use Equal method instead
// func (d Decimal) Equals(d2 Decimal) bool {
// 	return d.decimal.Equals(d2.decimal)
// }

// // GreaterThan (GT) returns true when d is greater than d2.
// func (d Decimal) GreaterThan(d2 Decimal) bool {
// 	return d.decimal.GreaterThan(d2.decimal)
// }

// // GreaterThanOrEqual (GTE) returns true when d is greater than or equal to d2.
// func (d Decimal) GreaterThanOrEqual(d2 Decimal) bool {
// 	return d.decimal.GreaterThanOrEqual(d2.decimal)
// }

// // LessThan (LT) returns true when d is less than d2.
// func (d Decimal) LessThan(d2 Decimal) bool {
// 	return d.decimal.LessThan(d2.decimal)
// }

// // LessThanOrEqual (LTE) returns true when d is less than or equal to d2.
// func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
// 	return d.decimal.LessThanOrEqual(d2.decimal)
// }

// // Sign returns:
// //
// //	-1 if d <  0
// //	 0 if d == 0
// //	+1 if d >  0
// //
// func (d Decimal) Sign() int {
// 	return d.decimal.Sign()
// }

// // IsPositive return
// //
// //	true if d > 0
// //	false if d == 0
// //	false if d < 0
// func (d Decimal) IsPositive() bool {
// 	return d.decimal.IsPositive()
// }

// // IsNegative return
// //
// //	true if d < 0
// //	false if d == 0
// //	false if d > 0
// func (d Decimal) IsNegative() bool {
// 	return d.decimal.IsNegative()
// }

// // IsZero return
// //
// //	true if d == 0
// //	false if d > 0
// //	false if d < 0
// func (d Decimal) IsZero() bool {
// 	return d.decimal.IsZero()
// }

// // Exponent returns the exponent, or scale component of the decimal.
// func (d Decimal) Exponent() int32 {
// 	return d.decimal.Exponent()
// }

// // Coefficient returns the coefficient of the decimal.  It is scaled by 10^Exponent()
// func (d Decimal) Coefficient() *big.Int {
// 	return d.decimal.Coefficient()
// }

// // IntPart returns the integer component of the decimal.
// func (d Decimal) IntPart() int64 {
// 	return d.decimal.IntPart()
// }

// // BigInt returns integer component of the decimal as a BigInt.
// func (d Decimal) BigInt() *big.Int {
// 	return d.decimal.BigInt()
// }

// // BigFloat returns decimal as BigFloat.
// // Be aware that casting decimal to BigFloat might cause a loss of precision.
// func (d Decimal) BigFloat() *big.Float {
// 	return d.decimal.BigFloat()
// }

// // Rat returns a rational number representation of the decimal.
// func (d Decimal) Rat() *big.Rat {
// 	return d.decimal.Rat()
// }

// // Float64 returns the nearest float64 value for d and a bool indicating
// // whether f represents d exactly.
// // For more details, see the documentation for big.Rat.Float64
// func (d Decimal) Float64() (f float64, exact bool) {
// 	return d.decimal.Float64()
// }

// String returns the string representation of the decimal
// with the fixed point.
//
// Example:
//
//     d := New(-12345, -3)
//     println(d.String())
//
// Output:
//
//     -12.345
//
// func (d Decimal) String() string {
// 	if len(d.decimals) == 0 {
// 		return "0"
// 	}

// 	return d.decimals[len(d.decimals)-1].String()
// }

// // StringFixed returns a rounded fixed-point string with places digits after
// // the decimal point.
// //
// // Example:
// //
// // 	   NewFromFloat(0).StringFixed(2) // output: "0.00"
// // 	   NewFromFloat(0).StringFixed(0) // output: "0"
// // 	   NewFromFloat(5.45).StringFixed(0) // output: "5"
// // 	   NewFromFloat(5.45).StringFixed(1) // output: "5.5"
// // 	   NewFromFloat(5.45).StringFixed(2) // output: "5.45"
// // 	   NewFromFloat(5.45).StringFixed(3) // output: "5.450"
// // 	   NewFromFloat(545).StringFixed(-1) // output: "550"
// //
// func (d Decimal) StringFixed(places int32) string {
// 	return d.decimal.StringFixed(places)
// }

// // StringFixedBank returns a banker rounded fixed-point string with places digits
// // after the decimal point.
// //
// // Example:
// //
// // 	   NewFromFloat(0).StringFixedBank(2) // output: "0.00"
// // 	   NewFromFloat(0).StringFixedBank(0) // output: "0"
// // 	   NewFromFloat(5.45).StringFixedBank(0) // output: "5"
// // 	   NewFromFloat(5.45).StringFixedBank(1) // output: "5.4"
// // 	   NewFromFloat(5.45).StringFixedBank(2) // output: "5.45"
// // 	   NewFromFloat(5.45).StringFixedBank(3) // output: "5.450"
// // 	   NewFromFloat(545).StringFixedBank(-1) // output: "540"
// //
// func (d Decimal) StringFixedBank(places int32) string {
// 	return d.decimal.StringFixedBank(places)
// }

// // StringFixedCash returns a Swedish/Cash rounded fixed-point string. For
// // more details see the documentation at function RoundCash.
// func (d Decimal) StringFixedCash(interval uint8) string {
// 	return d.decimal.StringFixedCash(interval)
// }

// // Round rounds the decimal to places decimal places.
// // If places < 0, it will round the integer part to the nearest 10^(-places).
// //
// // Example:
// //
// // 	   NewFromFloat(5.45).Round(1).String() // output: "5.5"
// // 	   NewFromFloat(545).Round(-1).String() // output: "550"
// //
// func (d Decimal) Round(places int32) Decimal {
// 	p := strconv.Itoa(int(places))

// 	return Decimal{
// 		decimal: d.decimal.Round(places),
// 		vars:    round + leftParen + p + rightParen + leftParen + d.vars + rightParen,
// 		formula: round + leftParen + p + rightParen + leftParen + d.formula + rightParen,
// 	}
// }

// // RoundBank rounds the decimal to places decimal places.
// // If the final digit to round is equidistant from the nearest two integers the
// // rounded value is taken as the even number
// //
// // If places < 0, it will round the integer part to the nearest 10^(-places).
// //
// // Examples:
// //
// // 	   NewFromFloat(5.45).Round(1).String() // output: "5.4"
// // 	   NewFromFloat(545).Round(-1).String() // output: "540"
// // 	   NewFromFloat(5.46).Round(1).String() // output: "5.5"
// // 	   NewFromFloat(546).Round(-1).String() // output: "550"
// // 	   NewFromFloat(5.55).Round(1).String() // output: "5.6"
// // 	   NewFromFloat(555).Round(-1).String() // output: "560"
// //
// func (d Decimal) RoundBank(places int32) Decimal {
// 	p := strconv.Itoa(int(places))

// 	return Decimal{
// 		decimal: d.decimal.RoundBank(places),
// 		vars:    roundBank + leftParen + p + rightParen + leftParen + d.vars + rightParen,
// 		formula: roundBank + leftParen + p + rightParen + leftParen + d.formula + rightParen,
// 	}
// }

// // RoundCash aka Cash/Penny/öre rounding rounds decimal to a specific
// // interval. The amount payable for a cash transaction is rounded to the nearest
// // multiple of the minimum currency unit available. The following intervals are
// // available: 5, 10, 25, 50 and 100; any other number throws a panic.
// //	    5:   5 cent rounding 3.43 => 3.45
// // 	   10:  10 cent rounding 3.45 => 3.50 (5 gets rounded up)
// // 	   25:  25 cent rounding 3.41 => 3.50
// // 	   50:  50 cent rounding 3.75 => 4.00
// // 	  100: 100 cent rounding 3.50 => 4.00
// // For more details: https://en.wikipedia.org/wiki/Cash_rounding
// func (d Decimal) RoundCash(interval uint8) Decimal {
// 	i := strconv.Itoa(int(interval))

// 	return Decimal{
// 		decimal: d.decimal.RoundCash(interval),
// 		vars:    roundCash + leftParen + i + rightParen + leftParen + d.vars + rightParen,
// 		formula: roundCash + leftParen + i + rightParen + leftParen + d.formula + rightParen,
// 	}
// }

// // Floor returns the nearest integer value less than or equal to d.
// func (d Decimal) Floor() Decimal {
// 	return Decimal{
// 		decimal: d.decimal.Floor(),
// 		vars:    floor + leftParen + d.vars + rightParen,
// 		formula: floor + leftParen + d.formula + rightParen,
// 	}
// }

// // Ceil returns the nearest integer value greater than or equal to d.
// func (d Decimal) Ceil() Decimal {
// 	return Decimal{
// 		decimal: d.decimal.Ceil(),
// 		vars:    ceil + leftParen + d.vars + rightParen,
// 		formula: ceil + leftParen + d.formula + rightParen,
// 	}
// }

// // Truncate truncates off digits from the number, without rounding.
// //
// // NOTE: precision is the last digit that will not be truncated (must be >= 0).
// //
// // Example:
// //
// //     decimal.NewFromString("123.456").Truncate(2).String() // "123.45"
// //
// func (d Decimal) Truncate(precision int32) Decimal {
// 	p := strconv.Itoa(int(precision))

// 	return Decimal{
// 		decimal: d.decimal.Truncate(precision),
// 		vars:    truncate + leftParen + p + rightParen + leftParen + d.vars + rightParen,
// 		formula: truncate + leftParen + p + rightParen + leftParen + d.formula + rightParen,
// 	}
// }

// // UnmarshalJSON implements the json.Unmarshaler interface.
// func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
// 	if err := d.decimal.UnmarshalJSON(decimalBytes); err != nil {
// 		return err
// 	}
// 	d.formula = d.String()

// 	return nil
// }

// // MarshalJSON implements the json.Marshaler interface.
// func (d Decimal) MarshalJSON() ([]byte, error) {
// 	return d.decimal.MarshalJSON()
// }

// // UnmarshalBinary implements the encoding.BinaryUnmarshaler interface. As a string representation
// // is already used when encoding to text, this method stores that string as []byte
// func (d *Decimal) UnmarshalBinary(data []byte) error {
// 	if err := d.decimal.UnmarshalBinary(data); err != nil {
// 		return err
// 	}
// 	d.formula = d.String()
// 	return nil
// }

// // MarshalBinary implements the encoding.BinaryMarshaler interface.
// func (d Decimal) MarshalBinary() (data []byte, err error) {
// 	return d.decimal.MarshalBinary()
// }

// // Scan implements the sql.Scanner interface for database deserialization.
// func (d *Decimal) Scan(value interface{}) error {
// 	if err := d.decimal.Scan(value); err != nil {
// 		return err
// 	}
// 	d.formula = d.String()
// 	return nil
// }

// // Value implements the driver.Valuer interface for database serialization.
// func (d Decimal) Value() (driver.Value, error) {
// 	return d.decimal.Value()
// }

// // UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// // deserialization.
// func (d *Decimal) UnmarshalText(text []byte) error {
// 	if err := d.decimal.UnmarshalText(text); err != nil {
// 		return err
// 	}
// 	d.formula = d.String()
// 	return nil
// }

// // MarshalText implements the encoding.TextMarshaler interface for XML
// // serialization.
// func (d Decimal) MarshalText() (text []byte, err error) {
// 	return d.decimal.MarshalText()
// }

// // GobEncode implements the gob.GobEncoder interface for gob serialization.
// func (d Decimal) GobEncode() ([]byte, error) {
// 	return d.decimal.GobEncode()
// }

// // GobDecode implements the gob.GobDecoder interface for gob serialization.
// func (d *Decimal) GobDecode(data []byte) error {
// 	if err := d.decimal.UnmarshalBinary(data); err != nil {
// 		return err
// 	}
// 	d.formula = d.String()
// 	return d.decimal.GobDecode(data)
// }

// // StringScaled first scales the decimal then calls .String() on it.
// // NOTE: buggy, unintuitive, and DEPRECATED! Use StringFixed instead.
// func (d Decimal) StringScaled(exp int32) string {
// 	return d.decimal.StringScaled(exp)
// }

// // Min returns the smallest Decimal that was passed in the arguments.
// //
// // To call this function with an array, you must do:
// //
// //     Min(arr[0], arr[1:]...)
// //
// // This makes it harder to accidentally call Min with 0 arguments.
// func Min(first Decimal, rest ...Decimal) Decimal {
// 	varsList := make([]string, 1+len(rest))
// 	varsList[0] = first.vars
// 	formulaList := make([]string, 1+len(rest))
// 	formulaList[0] = first.formula

// 	newRest := make([]decimal.Decimal, len(rest))
// 	for i, r := range rest {
// 		newRest[i] = r.decimal
// 		varsList[i+1] = r.vars
// 		formulaList[i+1] = r.formula
// 	}

// 	return Decimal{
// 		decimal: decimal.Min(first.decimal, newRest...),
// 		vars:    min + leftParen + strings.Join(varsList, comma) + rightParen,
// 		formula: min + leftParen + strings.Join(formulaList, comma) + rightParen,
// 	}
// }

// // Max returns the largest Decimal that was passed in the arguments.
// //
// // To call this function with an array, you must do:
// //
// //     Max(arr[0], arr[1:]...)
// //
// // This makes it harder to accidentally call Max with 0 arguments.
// func Max(first Decimal, rest ...Decimal) Decimal {
// 	varsList := make([]string, 1+len(rest))
// 	varsList[0] = first.vars
// 	formulaList := make([]string, 1+len(rest))
// 	formulaList[0] = first.formula

// 	newRest := make([]decimal.Decimal, len(rest))
// 	for i, r := range rest {
// 		newRest[i] = r.decimal
// 		varsList[i+1] = r.vars
// 		formulaList[i+1] = r.formula
// 	}

// 	return Decimal{
// 		decimal: decimal.Max(first.decimal, newRest...),
// 		vars:    max + leftParen + strings.Join(varsList, comma) + rightParen,
// 		formula: max + leftParen + strings.Join(formulaList, comma) + rightParen,
// 	}
// }

// // Sum returns the combined total of the provided first and rest Decimals
// func Sum(first Decimal, rest ...Decimal) Decimal {
// 	varsList := make([]string, 1+len(rest))
// 	varsList[0] = first.vars
// 	formulaList := make([]string, 1+len(rest))
// 	formulaList[0] = first.formula

// 	newRest := make([]decimal.Decimal, len(rest))
// 	for i, r := range rest {
// 		newRest[i] = r.decimal
// 		varsList[i+1] = r.vars
// 		formulaList[i+1] = r.formula
// 	}

// 	return Decimal{
// 		decimal: decimal.Sum(first.decimal, newRest...),
// 		vars:    sum + leftParen + strings.Join(varsList, comma) + rightParen,
// 		formula: sum + leftParen + strings.Join(formulaList, comma) + rightParen,
// 	}
// }

// // Avg returns the average value of the provided first and rest Decimals
// func Avg(first Decimal, rest ...Decimal) Decimal {
// 	varsList := make([]string, 1+len(rest))
// 	varsList[0] = first.vars
// 	formulaList := make([]string, 1+len(rest))
// 	formulaList[0] = first.formula

// 	newRest := make([]decimal.Decimal, len(rest))
// 	for i, r := range rest {
// 		newRest[i] = r.decimal
// 		varsList[i+1] = r.vars
// 		formulaList[i+1] = r.formula
// 	}

// 	return Decimal{
// 		decimal: decimal.Avg(first.decimal, newRest...),
// 		vars:    avg + leftParen + strings.Join(varsList, comma) + rightParen,
// 		formula: avg + leftParen + strings.Join(formulaList, comma) + rightParen,
// 	}
// }

// // RescalePair rescales two decimals to common exponential value (minimal exp of both decimals)
// func RescalePair(d1 Decimal, d2 Decimal) (Decimal, Decimal) {
// 	d3, d4 := decimal.RescalePair(d1.decimal, d2.decimal)
// 	return Decimal{name: d1.name, decimal: d3, vars: d1.name, formula: d3.String()},
// 		Decimal{name: d2.name, decimal: d4, vars: d2.name, formula: d4.String()}
// }

// func (d NullDecimal) Valid() bool {
// 	return d.decimal.Valid
// }

// func (d NullDecimal) Decimal() Decimal {
// 	return Decimal{
// 		name:    d.name,
// 		decimal: d.decimal.Decimal,
// 		vars:    d.name,
// 		formula: d.decimal.Decimal.String(),
// 	}
// }

// // Scan implements the sql.Scanner interface for database deserialization.
// func (d *NullDecimal) Scan(value interface{}) error {
// 	return d.decimal.Scan(value)
// }

// // Value implements the driver.Valuer interface for database serialization.
// func (d NullDecimal) Value() (driver.Value, error) {
// 	return d.decimal.Value()
// }

// // UnmarshalJSON implements the json.Unmarshaler interface.
// func (d *NullDecimal) UnmarshalJSON(decimalBytes []byte) error {
// 	return d.decimal.UnmarshalJSON(decimalBytes)
// }

// // MarshalJSON implements the json.Marshaler interface.
// func (d NullDecimal) MarshalJSON() ([]byte, error) {
// 	return d.decimal.MarshalJSON()
// }

// // Atan returns the arctangent, in radians, of x.
// func (d Decimal) Atan() Decimal {
// 	return Decimal{
// 		decimal: d.decimal.Atan(),
// 		vars:    atan + leftParen + d.vars + rightParen,
// 		formula: atan + leftParen + d.formula + rightParen,
// 	}
// }

// // Sin returns the sine of the radian argument x.
// func (d Decimal) Sin() Decimal {
// 	return Decimal{
// 		decimal: d.decimal.Sin(),
// 		vars:    sin + leftParen + d.vars + rightParen,
// 		formula: sin + leftParen + d.formula + rightParen,
// 	}
// }

// // Cos returns the cosine of the radian argument x.
// func (d Decimal) Cos() Decimal {
// 	return Decimal{
// 		decimal: d.decimal.Cos(),
// 		vars:    cos + leftParen + d.vars + rightParen,
// 		formula: cos + leftParen + d.formula + rightParen,
// 	}
// }

// // Tan returns the tangent of the radian argument x.
// func (d Decimal) Tan() Decimal {
// 	return Decimal{
// 		decimal: d.decimal.Tan(),
// 		vars:    tan + leftParen + d.vars + rightParen,
// 		formula: tan + leftParen + d.formula + rightParen,
// 	}
// }
