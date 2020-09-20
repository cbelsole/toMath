// Package tomath wraps github.com/shopspring/decimal library which implements an arbitrary precision fixed-point decimal.
//
// The zero-value of a Decimal is 0, as you would expect.
// The zero-value name is "?". To set a name for a decimal without a name use the SetName() method.
//
// The best way to create a new Decimal is to use decimal.NewFromStringWithName, ex:
//
//     n, err := decimal.NewFromStringWithName("var1", "1.3")
//     n.String() // output: "3.1"
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
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

type (
	// Decimal represents a fixed-point decimal. It is immutable.
	// number = value * 10 ^ exp
	Decimal struct {
		parens  bool
		name    string
		decimal decimal.Decimal
		vars    string
		formula string
	}

	// NullDecimal represents a nullable decimal with compatibility for
	// scanning null values from the database.
	NullDecimal struct {
		name    string
		decimal decimal.NullDecimal
	}
)

var (
	Zero = Decimal{
		decimal: decimal.Zero,
		name:    "zero",
		vars:    "zero",
		formula: "0",
	}
)

// SetName sets the name of the Decimal
func (d Decimal) SetName(name string) Decimal {
	d.name = name
	if d.vars == "" {
		d.vars = name
	}
	return d
}

// SetName gets the name of the Decimal
func (d Decimal) GetName() string {
	return d.name
}

// Resolve removes the underlying math from the decimal and replaces it with the
// current name and value.
func (d Decimal) Resolve() Decimal {
	return Decimal{
		name:    d.name,
		vars:    d.name,
		formula: d.String(),
		decimal: d.decimal,
	}
}

// ResolveTo is a wrapper around SetName() and Resolve().
func (d Decimal) ResolveTo(name string) Decimal {
	return d.SetName(name).Resolve()
}

// Decimal ejects the github.com/shopspring/decimal#Decimal
func (d Decimal) Decimal() decimal.Decimal {
	return d.decimal
}

// Math returns two strings representing the formula underlying the decimal. The
// first uses the decimal names. The second uses the decimal values. Both are
// follwed by an equals sign with the current name and value respectively.
func (d Decimal) Math() (string, string) {
	if d.name == "" {
		d.name = "?"
	}

	if d.vars == "" {
		d.vars = "?"
	}

	if d.formula == "" {
		d.formula = d.String()
	}

	d.vars = fmt.Sprintf("%s = %s", d.vars, d.name)
	d.formula = fmt.Sprintf("%s = %s", d.formula, d.decimal)

	return d.vars, d.formula
}

// New returns a new fixed-point decimal, value * 10 ^ exp.
func New(value int64, exp int32) Decimal {
	d := decimal.New(value, exp)
	return Decimal{decimal: d, formula: d.String()}
}

// NewWithName returns a new fixed-point decimal, value * 10 ^ exp with a given name.
func NewWithName(name string, value int64, exp int32) Decimal {
	d := decimal.New(value, exp)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}

// NewFromInt converts a int64 to Decimal.
//
// Example:
//
//     NewFromInt(123).String() // output: "123"
//     NewFromInt(-10).String() // output: "-10"
func NewFromInt(value int64) Decimal {
	d := decimal.NewFromInt(value)
	return Decimal{decimal: d, formula: d.String()}
}

// NewFromIntWithName converts a int64 to Decimal with a given name.
//
// Example:
//
//     NewFromIntWithName("var1", 123).String() // output: "123"
//     NewFromIntWithName("var1", -10).String() // output: "-10"
func NewFromIntWithName(name string, value int64) Decimal {
	d := decimal.NewFromInt(value)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}

// NewFromInt32 converts a int32 to Decimal.
//
// Example:
//
//     NewFromInt(123).String() // output: "123"
//     NewFromInt(-10).String() // output: "-10"
func NewFromInt32(value int32) Decimal {
	d := decimal.NewFromInt32(value)
	return Decimal{decimal: d, formula: d.String()}
}

// NewFromInt32WithName converts a int32 to Decimal with a given name.
//
// Example:
//
//     NewFromInt32WithName("var1", 123).String() // output: "123"
//     NewFromInt32WithName("var1", -10).String() // output: "-10"
func NewFromInt32WithName(name string, value int32) Decimal {
	d := decimal.NewFromInt32(value)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}

// NewFromBigInt returns a new Decimal from a big.Int, value * 10 ^ exp
func NewFromBigInt(value *big.Int, exp int32) Decimal {
	d := decimal.NewFromBigInt(value, exp)
	return Decimal{decimal: d, formula: d.String()}
}

// NewFromBigIntWithName returns a new Decimal from a big.Int, value * 10 ^ exp
// with a given name
func NewFromBigIntWithName(name string, value *big.Int, exp int32) Decimal {
	d := decimal.NewFromBigInt(value, exp)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
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

	return Decimal{decimal: d, formula: d.String()}, nil
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

	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}, nil
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
	return Decimal{decimal: d, formula: d.String()}
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
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
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
	return Decimal{decimal: d, formula: d.String()}
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
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
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
	return Decimal{decimal: d, formula: d.String()}
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
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
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
	return Decimal{decimal: d, formula: d.String()}
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
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}

// NewFromDecimal returns a new Decimal from github.com/shopspring/decimal#Decimal.
func NewFromDecimal(d decimal.Decimal) Decimal {
	return Decimal{decimal: d, formula: d.String()}
}

// NewFromDecimalWithName returns a new Decimal from github.com/shopspring/decimal#Decimal
// with a given name.
func NewFromDecimalWithName(name string, d decimal.Decimal) Decimal {
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	return Decimal{
		decimal: d.decimal.Abs(),
		vars:    fmt.Sprintf("abs(%s)", d.vars),
		formula: fmt.Sprintf("abs(%s)", d.formula),
	}
}

// Add returns d + d2.
func (d Decimal) Add(d2 Decimal) Decimal {
	return Decimal{
		decimal: d.decimal.Add(d2.decimal),
		vars:    fmt.Sprintf("%s + %s", d.vars, d2.vars),
		formula: fmt.Sprintf("%s + %s", d.formula, d2.formula),
		parens:  true,
	}
}

// Sub returns d - d2.
func (d Decimal) Sub(d2 Decimal) Decimal {
	return Decimal{
		decimal: d.decimal.Sub(d2.decimal),
		vars:    fmt.Sprintf("%s - %s", d.vars, d2.vars),
		formula: fmt.Sprintf("%s - %s", d.formula, d2.formula),
		parens:  true,
	}
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	return Decimal{
		decimal: d.decimal.Neg(),
		vars:    fmt.Sprintf("neg(%s)", d.vars),
		formula: fmt.Sprintf("neg(%s)", d.formula),
	}
}

// Mul returns d * d2.
func (d Decimal) Mul(d2 Decimal) Decimal {
	dec := Decimal{decimal: d.decimal.Mul(d2.decimal)}

	var format string
	if d.parens && d2.parens {
		format = "(%s) * (%s)"
	} else if d.parens {
		format = "(%s) * %s"
	} else if d2.parens {
		format = "%s * (%s)"
	} else {
		format = "%s * %s"
	}

	dec.vars = fmt.Sprintf(format, d.vars, d2.vars)
	dec.formula = fmt.Sprintf(format, d.formula, d2.formula)

	return dec
}

// Shift shifts the decimal in base 10.
// It shifts left when shift is positive and right if shift is negative.
// In simpler terms, the given value for shift is added to the exponent
// of the decimal.
func (d Decimal) Shift(s int32) Decimal {
	return Decimal{
		decimal: d.decimal.Shift(s),
		vars:    fmt.Sprintf("shift(%d)(%s)", s, d.vars),
		formula: fmt.Sprintf("shift(%d)(%s)", s, d.formula),
	}
}

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
func (d Decimal) Div(d2 Decimal) Decimal {
	dec := Decimal{decimal: d.decimal.Div(d2.decimal)}

	var format string
	if d.parens && d2.parens {
		format = "(%s) / (%s)"
	} else if d.parens {
		format = "(%s) / %s"
	} else if d2.parens {
		format = "%s / (%s)"
	} else {
		format = "%s / %s"
	}

	dec.vars = fmt.Sprintf(format, d.vars, d2.vars)
	dec.formula = fmt.Sprintf(format, d.formula, d2.formula)

	return dec
}

// QuoRem does divsion with remainder
// d.QuoRem(d2,precision) returns quotient q and remainder r such that
//   d = d2 * q + r, q an integer multiple of 10^(-precision)
//   0 <= r < abs(d2) * 10 ^(-precision) if d>=0
//   0 >= r > -abs(d2) * 10 ^(-precision) if d<0
// Note that precision<0 is allowed as input.
func (d Decimal) QuoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
	d3, d4 := d.decimal.QuoRem(d2.decimal, precision)

	var format string
	if d.parens && d2.parens {
		format = "quoRem(%d)((%s) / (%s))"
	} else if d.parens {
		format = "quoRem(%d)((%s) / %s)"
	} else if d2.parens {
		format = "quoRem(%d)(%s / (%s))"
	} else {
		format = "quoRem(%d)(%s / %s)"
	}

	vars := fmt.Sprintf(format, precision, d.vars, d2.vars)
	formula := fmt.Sprintf(format, precision, d.formula, d2.formula)

	return Decimal{name: d.name + d2.name + "Quotient", decimal: d3, vars: vars, formula: formula},
		Decimal{name: d.name + d2.name + "Remainder", decimal: d4, vars: vars, formula: formula}
}

// DivRound divides and rounds to a given precision
// i.e. to an integer multiple of 10^(-precision)
//   for a positive quotient digit 5 is rounded up, away from 0
//   if the quotient is negative then digit 5 is rounded down, away from 0
// Note that precision<0 is allowed as input.
func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
	dec := Decimal{decimal: d.decimal.DivRound(d2.decimal, precision)}

	var format string
	if d.parens && d2.parens {
		format = "divRound(%d)((%s) / (%s))"
	} else if d.parens {
		format = "divRound(%d)((%s) / %s)"
	} else if d2.parens {
		format = "divRound(%d)(%s / (%s))"
	} else {
		format = "divRound(%d)(%s / %s)"
	}

	dec.vars = fmt.Sprintf(format, precision, d.vars, d2.vars)
	dec.formula = fmt.Sprintf(format, precision, d.formula, d2.formula)

	return dec
}

// Mod returns d % d2.
func (d Decimal) Mod(d2 Decimal) Decimal {
	dec := Decimal{decimal: d.decimal.Mod(d2.decimal)}

	var format string
	if d.parens && d2.parens {
		format = "(%s) %% (%s)"
	} else if d.parens {
		format = "(%s) %% %s"
	} else if d2.parens {
		format = "%s %% (%s)"
	} else {
		format = "%s %% %s"
	}

	dec.vars = fmt.Sprintf(format, d.vars, d2.vars)
	dec.formula = fmt.Sprintf(format, d.formula, d2.formula)

	return dec
}

// Pow returns d to the power d2
func (d Decimal) Pow(d2 Decimal) Decimal {
	dec := Decimal{decimal: d.decimal.Pow(d2.decimal)}

	var format string
	if d.parens && d2.parens {
		format = "(%s)^(%s)"
	} else if d.parens {
		format = "(%s)^%s"
	} else if d2.parens {
		format = "%s^(%s)"
	} else {
		format = "%s^%s"
	}

	dec.vars = fmt.Sprintf(format, d.vars, d2.vars)
	dec.formula = fmt.Sprintf(format, d.formula, d2.formula)

	return dec
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//     -1 if d <  d2
//      0 if d == d2
//     +1 if d >  d2
//
func (d Decimal) Cmp(d2 Decimal) int {
	return d.decimal.Cmp(d2.decimal)
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (d Decimal) Equal(d2 Decimal) bool {
	return d.decimal.Equal(d2.decimal)
}

// Equals is deprecated, please use Equal method instead
func (d Decimal) Equals(d2 Decimal) bool {
	return d.decimal.Equals(d2.decimal)
}

// GreaterThan (GT) returns true when d is greater than d2.
func (d Decimal) GreaterThan(d2 Decimal) bool {
	return d.decimal.GreaterThan(d2.decimal)
}

// GreaterThanOrEqual (GTE) returns true when d is greater than or equal to d2.
func (d Decimal) GreaterThanOrEqual(d2 Decimal) bool {
	return d.decimal.GreaterThanOrEqual(d2.decimal)
}

// LessThan (LT) returns true when d is less than d2.
func (d Decimal) LessThan(d2 Decimal) bool {
	return d.decimal.LessThan(d2.decimal)
}

// LessThanOrEqual (LTE) returns true when d is less than or equal to d2.
func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	return d.decimal.LessThanOrEqual(d2.decimal)
}

// Sign returns:
//
//	-1 if d <  0
//	 0 if d == 0
//	+1 if d >  0
//
func (d Decimal) Sign() int {
	return d.decimal.Sign()
}

// IsPositive return
//
//	true if d > 0
//	false if d == 0
//	false if d < 0
func (d Decimal) IsPositive() bool {
	return d.decimal.IsPositive()
}

// IsNegative return
//
//	true if d < 0
//	false if d == 0
//	false if d > 0
func (d Decimal) IsNegative() bool {
	return d.decimal.IsNegative()
}

// IsZero return
//
//	true if d == 0
//	false if d > 0
//	false if d < 0
func (d Decimal) IsZero() bool {
	return d.decimal.IsZero()
}

// Exponent returns the exponent, or scale component of the decimal.
func (d Decimal) Exponent() int32 {
	return d.decimal.Exponent()
}

// Coefficient returns the coefficient of the decimal.  It is scaled by 10^Exponent()
func (d Decimal) Coefficient() *big.Int {
	return d.decimal.Coefficient()
}

// IntPart returns the integer component of the decimal.
func (d Decimal) IntPart() int64 {
	return d.decimal.IntPart()
}

// BigInt returns integer component of the decimal as a BigInt.
func (d Decimal) BigInt() *big.Int {
	return d.decimal.BigInt()
}

// BigFloat returns decimal as BigFloat.
// Be aware that casting decimal to BigFloat might cause a loss of precision.
func (d Decimal) BigFloat() *big.Float {
	return d.decimal.BigFloat()
}

// Rat returns a rational number representation of the decimal.
func (d Decimal) Rat() *big.Rat {
	return d.decimal.Rat()
}

// Float64 returns the nearest float64 value for d and a bool indicating
// whether f represents d exactly.
// For more details, see the documentation for big.Rat.Float64
func (d Decimal) Float64() (f float64, exact bool) {
	return d.decimal.Float64()
}

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
func (d Decimal) String() string {
	return d.decimal.String()
}

// StringFixed returns a rounded fixed-point string with places digits after
// the decimal point.
//
// Example:
//
// 	   NewFromFloat(0).StringFixed(2) // output: "0.00"
// 	   NewFromFloat(0).StringFixed(0) // output: "0"
// 	   NewFromFloat(5.45).StringFixed(0) // output: "5"
// 	   NewFromFloat(5.45).StringFixed(1) // output: "5.5"
// 	   NewFromFloat(5.45).StringFixed(2) // output: "5.45"
// 	   NewFromFloat(5.45).StringFixed(3) // output: "5.450"
// 	   NewFromFloat(545).StringFixed(-1) // output: "550"
//
func (d Decimal) StringFixed(places int32) string {
	return d.decimal.StringFixed(places)
}

// StringFixedBank returns a banker rounded fixed-point string with places digits
// after the decimal point.
//
// Example:
//
// 	   NewFromFloat(0).StringFixedBank(2) // output: "0.00"
// 	   NewFromFloat(0).StringFixedBank(0) // output: "0"
// 	   NewFromFloat(5.45).StringFixedBank(0) // output: "5"
// 	   NewFromFloat(5.45).StringFixedBank(1) // output: "5.4"
// 	   NewFromFloat(5.45).StringFixedBank(2) // output: "5.45"
// 	   NewFromFloat(5.45).StringFixedBank(3) // output: "5.450"
// 	   NewFromFloat(545).StringFixedBank(-1) // output: "540"
//
func (d Decimal) StringFixedBank(places int32) string {
	return d.decimal.StringFixedBank(places)
}

// StringFixedCash returns a Swedish/Cash rounded fixed-point string. For
// more details see the documentation at function RoundCash.
func (d Decimal) StringFixedCash(interval uint8) string {
	return d.decimal.StringFixedCash(interval)
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Example:
//
// 	   NewFromFloat(5.45).Round(1).String() // output: "5.5"
// 	   NewFromFloat(545).Round(-1).String() // output: "550"
//
func (d Decimal) Round(places int32) Decimal {
	return Decimal{
		decimal: d.decimal.Round(places),
		vars:    fmt.Sprintf("round(%d)(%s)", places, d.vars),
		formula: fmt.Sprintf("round(%d)(%s)", places, d.formula),
	}
}

// RoundBank rounds the decimal to places decimal places.
// If the final digit to round is equidistant from the nearest two integers the
// rounded value is taken as the even number
//
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Examples:
//
// 	   NewFromFloat(5.45).Round(1).String() // output: "5.4"
// 	   NewFromFloat(545).Round(-1).String() // output: "540"
// 	   NewFromFloat(5.46).Round(1).String() // output: "5.5"
// 	   NewFromFloat(546).Round(-1).String() // output: "550"
// 	   NewFromFloat(5.55).Round(1).String() // output: "5.6"
// 	   NewFromFloat(555).Round(-1).String() // output: "560"
//
func (d Decimal) RoundBank(places int32) Decimal {
	return Decimal{
		decimal: d.decimal.RoundBank(places),
		vars:    fmt.Sprintf("roundBank(%d)(%s)", places, d.vars),
		formula: fmt.Sprintf("roundBank(%d)(%s)", places, d.formula),
	}
}

// RoundCash aka Cash/Penny/öre rounding rounds decimal to a specific
// interval. The amount payable for a cash transaction is rounded to the nearest
// multiple of the minimum currency unit available. The following intervals are
// available: 5, 10, 25, 50 and 100; any other number throws a panic.
//	    5:   5 cent rounding 3.43 => 3.45
// 	   10:  10 cent rounding 3.45 => 3.50 (5 gets rounded up)
// 	   25:  25 cent rounding 3.41 => 3.50
// 	   50:  50 cent rounding 3.75 => 4.00
// 	  100: 100 cent rounding 3.50 => 4.00
// For more details: https://en.wikipedia.org/wiki/Cash_rounding
func (d Decimal) RoundCash(interval uint8) Decimal {
	return Decimal{
		decimal: d.decimal.RoundCash(interval),
		vars:    fmt.Sprintf("roundCash(%d)(%s)", interval, d.vars),
		formula: fmt.Sprintf("roundCash(%d)(%s)", interval, d.formula),
	}
}

// Floor returns the nearest integer value less than or equal to d.
func (d Decimal) Floor() Decimal {
	return Decimal{
		decimal: d.decimal.Floor(),
		vars:    fmt.Sprintf("floor(%s)", d.vars),
		formula: fmt.Sprintf("floor(%s)", d.formula),
	}
}

// Ceil returns the nearest integer value greater than or equal to d.
func (d Decimal) Ceil() Decimal {
	return Decimal{
		decimal: d.decimal.Ceil(),
		vars:    fmt.Sprintf("ceil(%s)", d.vars),
		formula: fmt.Sprintf("ceil(%s)", d.formula),
	}
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//     decimal.NewFromString("123.456").Truncate(2).String() // "123.45"
//
func (d Decimal) Truncate(precision int32) Decimal {
	return Decimal{
		decimal: d.decimal.Truncate(precision),
		vars:    fmt.Sprintf("truncate(%d)(%s)", precision, d.vars),
		formula: fmt.Sprintf("truncate(%d)(%s)", precision, d.formula),
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	if err := d.decimal.UnmarshalJSON(decimalBytes); err != nil {
		return err
	}
	d.formula = d.String()

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.decimal.MarshalJSON()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface. As a string representation
// is already used when encoding to text, this method stores that string as []byte
func (d *Decimal) UnmarshalBinary(data []byte) error {
	if err := d.decimal.UnmarshalBinary(data); err != nil {
		return err
	}
	d.formula = d.String()
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Decimal) MarshalBinary() (data []byte, err error) {
	return d.decimal.MarshalBinary()
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Decimal) Scan(value interface{}) error {
	if err := d.decimal.Scan(value); err != nil {
		return err
	}
	d.formula = d.String()
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (d Decimal) Value() (driver.Value, error) {
	return d.decimal.Value()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (d *Decimal) UnmarshalText(text []byte) error {
	if err := d.decimal.UnmarshalText(text); err != nil {
		return err
	}
	d.formula = d.String()
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (d Decimal) MarshalText() (text []byte, err error) {
	return d.decimal.MarshalText()
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (d Decimal) GobEncode() ([]byte, error) {
	return d.decimal.GobEncode()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (d *Decimal) GobDecode(data []byte) error {
	if err := d.decimal.UnmarshalBinary(data); err != nil {
		return err
	}
	d.formula = d.String()
	return d.decimal.GobDecode(data)
}

// StringScaled first scales the decimal then calls .String() on it.
// NOTE: buggy, unintuitive, and DEPRECATED! Use StringFixed instead.
func (d Decimal) StringScaled(exp int32) string {
	return d.decimal.StringScaled(exp)
}

// Min returns the smallest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//     Min(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Min with 0 arguments.
func Min(first Decimal, rest ...Decimal) Decimal {
	varsList := make([]string, 1+len(rest))
	varsList[0] = first.vars
	formulaList := make([]string, 1+len(rest))
	formulaList[0] = first.formula

	newRest := make([]decimal.Decimal, len(rest))
	for i, r := range rest {
		newRest[i] = r.decimal
		varsList[i+1] = r.vars
		formulaList[i+1] = r.formula
	}

	return Decimal{
		decimal: decimal.Min(first.decimal, newRest...),
		vars:    fmt.Sprintf("min(%s)", strings.Join(varsList, ", ")),
		formula: fmt.Sprintf("min(%s)", strings.Join(formulaList, ", ")),
	}
}

// Max returns the largest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//     Max(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Max with 0 arguments.
func Max(first Decimal, rest ...Decimal) Decimal {
	varsList := make([]string, 1+len(rest))
	varsList[0] = first.vars
	formulaList := make([]string, 1+len(rest))
	formulaList[0] = first.formula

	newRest := make([]decimal.Decimal, len(rest))
	for i, r := range rest {
		newRest[i] = r.decimal
		varsList[i+1] = r.vars
		formulaList[i+1] = r.formula
	}

	return Decimal{
		decimal: decimal.Max(first.decimal, newRest...),
		vars:    fmt.Sprintf("max(%s)", strings.Join(varsList, ", ")),
		formula: fmt.Sprintf("max(%s)", strings.Join(formulaList, ", ")),
	}
}

// Sum returns the combined total of the provided first and rest Decimals
func Sum(first Decimal, rest ...Decimal) Decimal {
	varsList := make([]string, 1+len(rest))
	varsList[0] = first.vars
	formulaList := make([]string, 1+len(rest))
	formulaList[0] = first.formula

	newRest := make([]decimal.Decimal, len(rest))
	for i, r := range rest {
		newRest[i] = r.decimal
		varsList[i+1] = r.vars
		formulaList[i+1] = r.formula
	}

	return Decimal{
		decimal: decimal.Sum(first.decimal, newRest...),
		vars:    fmt.Sprintf("sum(%s)", strings.Join(varsList, ", ")),
		formula: fmt.Sprintf("sum(%s)", strings.Join(formulaList, ", ")),
	}
}

// Avg returns the average value of the provided first and rest Decimals
func Avg(first Decimal, rest ...Decimal) Decimal {
	varsList := make([]string, 1+len(rest))
	varsList[0] = first.vars
	formulaList := make([]string, 1+len(rest))
	formulaList[0] = first.formula

	newRest := make([]decimal.Decimal, len(rest))
	for i, r := range rest {
		newRest[i] = r.decimal
		varsList[i+1] = r.vars
		formulaList[i+1] = r.formula
	}

	return Decimal{
		decimal: decimal.Avg(first.decimal, newRest...),
		vars:    fmt.Sprintf("avg(%s)", strings.Join(varsList, ", ")),
		formula: fmt.Sprintf("avg(%s)", strings.Join(formulaList, ", ")),
	}
}

// RescalePair rescales two decimals to common exponential value (minimal exp of both decimals)
func RescalePair(d1 Decimal, d2 Decimal) (Decimal, Decimal) {
	d3, d4 := decimal.RescalePair(d1.decimal, d2.decimal)
	return Decimal{name: d1.name, decimal: d3, vars: d1.name, formula: d3.String()},
		Decimal{name: d2.name, decimal: d4, vars: d2.name, formula: d4.String()}
}

func (d NullDecimal) Valid() bool {
	return d.decimal.Valid
}

func (d NullDecimal) Decimal() Decimal {
	return Decimal{
		name:    d.name,
		decimal: d.decimal.Decimal,
		vars:    d.name,
		formula: d.decimal.Decimal.String(),
	}
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *NullDecimal) Scan(value interface{}) error {
	return d.decimal.Scan(value)
}

// Value implements the driver.Valuer interface for database serialization.
func (d NullDecimal) Value() (driver.Value, error) {
	return d.decimal.Value()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *NullDecimal) UnmarshalJSON(decimalBytes []byte) error {
	return d.decimal.UnmarshalJSON(decimalBytes)
}

// MarshalJSON implements the json.Marshaler interface.
func (d NullDecimal) MarshalJSON() ([]byte, error) {
	return d.decimal.MarshalJSON()
}

// Atan returns the arctangent, in radians, of x.
func (d Decimal) Atan() Decimal {
	return Decimal{
		decimal: d.decimal.Atan(),
		vars:    fmt.Sprintf("atan(%s)", d.vars),
		formula: fmt.Sprintf("atan(%s)", d.formula),
	}
}

// Sin returns the sine of the radian argument x.
func (d Decimal) Sin() Decimal {
	return Decimal{
		decimal: d.decimal.Sin(),
		vars:    fmt.Sprintf("sin(%s)", d.vars),
		formula: fmt.Sprintf("sin(%s)", d.formula),
	}
}

// Cos returns the cosine of the radian argument x.
func (d Decimal) Cos() Decimal {
	return Decimal{
		decimal: d.decimal.Cos(),
		vars:    fmt.Sprintf("cos(%s)", d.vars),
		formula: fmt.Sprintf("cos(%s)", d.formula),
	}
}

// Tan returns the tangent of the radian argument x.
func (d Decimal) Tan() Decimal {
	return Decimal{
		decimal: d.decimal.Tan(),
		vars:    fmt.Sprintf("tan(%s)", d.vars),
		formula: fmt.Sprintf("tan(%s)", d.formula),
	}
}
