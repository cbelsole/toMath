package tomath

import (
	"database/sql/driver"
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
)

const (
	new operation = iota
	abs
	add
	sub
	neg
	mul
	shift
	div
	quoRem
	divRound
	mod
	pow
	round
	roundBank
	roundCash
	floor
	ceil
	truncate
)

type (
	operation int
	Decimal   struct {
		vars, formula         string
		name, leftValue       string
		rightName, rightValue string
		precision             int
		decimal               decimal.Decimal
		operations            []operation
		resultName            string
	}
)

// TODO: Implement me to get rid of the ?
func (d Decimal) Result(name string) Decimal {
	d.resultName = name
	return d
}

func (d Decimal) Name() string {
	return d.name
}

func (d Decimal) Decimal() decimal.Decimal {
	return d.decimal
}

func (d Decimal) Math() (string, string) {
	vars, formula := d.math()

	if len(d.operations) != 1 {
		if d.resultName == "" {
			vars = fmt.Sprintf("%s = ?", vars)
		} else {
			vars = fmt.Sprintf("%s = %s", vars, d.resultName)
		}

		formula = fmt.Sprintf("%s = %s", formula, d.decimal)
	}

	return vars, formula
}

func (d Decimal) math() (string, string) {
	paren := true
	vars := d.vars
	formula := d.formula

	if vars == "" {
		paren = false
		vars = d.name
		formula = d.leftValue
	}

	for _, operation := range d.operations {
		switch operation {
		case new:
		case abs:
			vars = fmt.Sprintf("abs(%s)", vars)
			formula = fmt.Sprintf("abs(%s)", formula)
		case add:
			if paren {
				vars = fmt.Sprintf("(%s)", vars)
				formula = fmt.Sprintf("(%s)", formula)
			}

			vars = fmt.Sprintf("%s + %s", vars, d.rightName)
			formula = fmt.Sprintf("%s + %s", formula, d.rightValue)
		case sub:
			if paren {
				vars = fmt.Sprintf("(%s)", vars)
				formula = fmt.Sprintf("(%s)", formula)
			}

			vars = fmt.Sprintf("%s - %s", vars, d.rightName)
			formula = fmt.Sprintf("%s - %s", formula, d.rightValue)
		case neg:
			vars = fmt.Sprintf("neg(%s)", vars)
			formula = fmt.Sprintf("neg(%s)", formula)
		case mul:
			if paren {
				vars = fmt.Sprintf("(%s)", vars)
				formula = fmt.Sprintf("(%s)", formula)
			}

			vars = fmt.Sprintf("%s * %s", vars, d.rightName)
			formula = fmt.Sprintf("%s * %s", formula, d.rightValue)
		case shift:
			vars = fmt.Sprintf("shift(%d)(%s)", d.precision, vars)
			formula = fmt.Sprintf("shift(%d)(%s)", d.precision, formula)
		case div:
			if paren {
				vars = fmt.Sprintf("(%s)", vars)
				formula = fmt.Sprintf("(%s)", formula)
			}

			vars = fmt.Sprintf("%s / %s", vars, d.rightName)
			formula = fmt.Sprintf("%s / %s", formula, d.rightValue)
		case quoRem:
			if paren {
				vars = fmt.Sprintf("(%s)", vars)
				formula = fmt.Sprintf("(%s)", formula)
			}

			vars = fmt.Sprintf("quoRem(%d)(%s / %s)", d.precision, vars, d.rightName)
			formula = fmt.Sprintf("quoRem(%d)(%s / %s)", d.precision, formula, d.rightValue)
		case divRound:
			vars = fmt.Sprintf("round(%d)(%s / %s)", d.precision, vars, d.rightName)
			formula = fmt.Sprintf("round(%d)(%s / %s)", d.precision, formula, d.rightValue)
		case mod:
			if paren {
				vars = fmt.Sprintf("(%s)", vars)
				formula = fmt.Sprintf("(%s)", formula)
			}

			vars = fmt.Sprintf("%s %% %s", vars, d.rightName)
			formula = fmt.Sprintf("%s %% %s", formula, d.rightValue)
		case pow:
			if paren {
				vars = fmt.Sprintf("(%s)", vars)
				formula = fmt.Sprintf("(%s)", formula)
			}

			vars = fmt.Sprintf("%s^%s", vars, d.rightName)
			formula = fmt.Sprintf("%s^%s", formula, d.rightValue)
		case round:
			vars = fmt.Sprintf("round(%d)(%s)", d.precision, vars)
			formula = fmt.Sprintf("round(%d)(%s)", d.precision, formula)
		case roundBank:
			vars = fmt.Sprintf("roundBank(%d)(%s)", d.precision, vars)
			formula = fmt.Sprintf("roundBank(%d)(%s)", d.precision, formula)
		case roundCash:
			vars = fmt.Sprintf("roundCash(%d)(%s)", d.precision, vars)
			formula = fmt.Sprintf("roundCash(%d)(%s)", d.precision, formula)
		case floor:
			vars = fmt.Sprintf("floor(%s)", vars)
			formula = fmt.Sprintf("floor(%s)", formula)
		case ceil:
			vars = fmt.Sprintf("ceil(%s)", vars)
			formula = fmt.Sprintf("ceil(%s)", formula)
		case truncate:
			vars = fmt.Sprintf("truncate(%d)(%s)", d.precision, vars)
			formula = fmt.Sprintf("truncate(%d)(%s)", d.precision, formula)
		}
	}

	return vars, formula
}

func New(name string, value int64, exp int32) Decimal {
	d := decimal.New(value, exp)
	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}
}

func NewFromInt(name string, value int64) Decimal {
	d := decimal.NewFromInt(value)
	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}
}

func NewFromInt32(name string, value int32) Decimal {
	d := decimal.NewFromInt32(value)
	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}
}
func NewFromBigInt(name string, value *big.Int, exp int32) Decimal {
	d := decimal.NewFromBigInt(value, exp)
	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}
}
func NewFromString(name string, value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}, nil
}
func RequireFromString(name string, value string) Decimal {
	d := decimal.RequireFromString(value)
	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}
}
func NewFromFloat(name string, value float64) Decimal {
	d := decimal.NewFromFloat(value)
	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}
}
func NewFromFloat32(name string, value float32) Decimal {
	d := decimal.NewFromFloat32(value)
	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}
}
func NewFromFloatWithExponent(name string, value float64, exp int32) Decimal {
	d := decimal.NewFromFloatWithExponent(value, exp)
	return Decimal{name: name, decimal: d, leftValue: d.String(), operations: []operation{new}}
}

func (d Decimal) Abs() Decimal {
	return Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		decimal:    d.decimal.Abs(),
		operations: append(d.operations, abs),
	}
}

func (d Decimal) Add(d2 Decimal) Decimal {
	out := Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		decimal:    d.decimal.Add(d2.decimal),
		operations: append(d.operations, add),
		vars:       d.vars,
		formula:    d.formula,
	}

	out.vars, out.formula = out.math()
	out.operations = []operation{}

	return out
}

func (d Decimal) Sub(d2 Decimal) Decimal {
	out := Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		decimal:    d.decimal.Sub(d2.decimal),
		operations: append(d.operations, sub),
		vars:       d.vars,
		formula:    d.formula,
	}

	out.vars, out.formula = out.math()
	out.operations = []operation{}

	return out
}

func (d Decimal) Neg() Decimal {
	return Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		decimal:    d.decimal.Neg(),
		operations: append(d.operations, neg),
	}
}

func (d Decimal) Mul(d2 Decimal) Decimal {
	out := Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		decimal:    d.decimal.Mul(d2.decimal),
		operations: append(d.operations, mul),
		vars:       d.vars,
		formula:    d.formula,
	}

	out.vars, out.formula = out.math()
	out.operations = []operation{}

	return out
}

func (d Decimal) Shift(s int32) Decimal {
	return Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		precision:  int(s),
		decimal:    d.decimal.Shift(s),
		operations: append(d.operations, shift),
	}
}

func (d Decimal) Div(d2 Decimal) Decimal {
	out := Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		decimal:    d.decimal.Div(d2.decimal),
		operations: append(d.operations, div),
		vars:       d.vars,
		formula:    d.formula,
	}

	out.vars, out.formula = out.math()
	out.operations = []operation{}

	return out
}

func (d Decimal) QuoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
	d3, d4 := d.decimal.QuoRem(d2.decimal, precision)

	out1 := Decimal{
		name:       d.name,
		resultName: d.name + d2.name + "Quotient",
		precision:  int(precision),
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		operations: append(d.operations, quoRem),
		decimal:    d3,
		vars:       d.vars,
		formula:    d.formula,
	}

	out1.vars, out1.formula = out1.math()
	out1.operations = []operation{}

	out2 := Decimal{
		name:       d.name,
		resultName: d.name + d2.name + "Remainder",
		precision:  int(precision),
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		operations: append(d.operations, quoRem),
		decimal:    d4,
		vars:       d.vars,
		formula:    d.formula,
	}

	out2.vars, out2.formula = out2.math()
	out2.operations = []operation{}

	return out1, out2
}

func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
	out := Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		precision:  int(precision),
		decimal:    d.decimal.DivRound(d2.decimal, precision),
		operations: append(d.operations, divRound),
		vars:       d.vars,
		formula:    d.formula,
	}

	out.vars, out.formula = out.math()
	out.operations = []operation{}

	return out
}

func (d Decimal) Mod(d2 Decimal) Decimal {
	out := Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		decimal:    d.decimal.Mod(d2.decimal),
		operations: append(d.operations, mod),
		vars:       d.vars,
		formula:    d.formula,
	}

	out.vars, out.formula = out.math()
	out.operations = []operation{}

	return out
}

func (d Decimal) Pow(d2 Decimal) Decimal {
	out := Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		rightValue: d2.decimal.String(),
		rightName:  d2.name,
		decimal:    d.decimal.Pow(d2.decimal),
		operations: append(d.operations, pow),
		vars:       d.vars,
		formula:    d.formula,
	}

	out.vars, out.formula = out.math()
	out.operations = []operation{}

	return out
}

func (d Decimal) Cmp(d2 Decimal) int {
	return d.decimal.Cmp(d2.decimal)
}

func (d Decimal) Equal(d2 Decimal) bool {
	return d.decimal.Equal(d2.decimal)
}

func (d Decimal) Equals(d2 Decimal) bool {
	return d.decimal.Equals(d2.decimal)
}

func (d Decimal) GreaterThan(d2 Decimal) bool {
	return d.decimal.GreaterThan(d2.decimal)
}

func (d Decimal) GreaterThanOrEqual(d2 Decimal) bool {
	return d.decimal.GreaterThanOrEqual(d2.decimal)
}

func (d Decimal) LessThan(d2 Decimal) bool {
	return d.decimal.LessThan(d2.decimal)
}

func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	return d.decimal.LessThanOrEqual(d2.decimal)
}

func (d Decimal) Sign() int {
	return d.decimal.Sign()
}

func (d Decimal) IsPositive() bool {
	return d.decimal.IsPositive()
}

func (d Decimal) IsNegative() bool {
	return d.decimal.IsNegative()
}

func (d Decimal) IsZero() bool {
	return d.decimal.IsZero()
}

func (d Decimal) Exponent() int32 {
	return d.decimal.Exponent()
}

func (d Decimal) Coefficient() *big.Int {
	return d.decimal.Coefficient()
}

func (d Decimal) IntPart() int64 {
	return d.decimal.IntPart()
}

func (d Decimal) BigInt() *big.Int {
	return d.decimal.BigInt()
}

func (d Decimal) BigFloat() *big.Float {
	return d.decimal.BigFloat()
}

func (d Decimal) Rat() *big.Rat {
	return d.decimal.Rat()
}

func (d Decimal) Float64() (f float64, exact bool) {
	return d.decimal.Float64()
}

func (d Decimal) String() string {
	return d.decimal.String()
}

func (d Decimal) StringFixed(places int32) string {
	return d.decimal.StringFixed(places)
}

func (d Decimal) StringFixedBank(places int32) string {
	return d.decimal.StringFixedBank(places)
}

func (d Decimal) StringFixedCash(interval uint8) string {
	return d.decimal.StringFixedCash(interval)
}

func (d Decimal) Round(places int32) Decimal {
	out := Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		precision:  int(places),
		decimal:    d.decimal.Round(places),
		operations: append(d.operations, round),
		vars:       d.vars,
		formula:    d.formula,
	}

	out.vars, out.formula = out.math()
	out.operations = []operation{}

	return out
}

func (d Decimal) RoundBank(places int32) Decimal {
	return Decimal{
		name:       d.name,
		leftValue:  d.decimal.String(),
		precision:  int(places),
		decimal:    d.decimal.RoundBank(places),
		operations: append(d.operations, roundBank),
	}
}

func (d Decimal) RoundCash(interval uint8) Decimal {
	d.leftValue = d.decimal.String()
	d.precision = int(interval)

	d.decimal = d.decimal.RoundCash(interval)
	d.operations = append(d.operations, roundCash)

	return d
}

func (d Decimal) Floor() Decimal {
	d.leftValue = d.decimal.String()

	d.decimal = d.decimal.Floor()
	d.operations = append(d.operations, floor)

	return d
}

func (d Decimal) Ceil() Decimal {
	d.leftValue = d.decimal.String()

	d.decimal = d.decimal.Ceil()
	d.operations = append(d.operations, ceil)

	return d
}

func (d Decimal) Truncate(precision int32) Decimal {
	d.leftValue = d.decimal.String()
	d.precision = int(precision)

	d.decimal = d.decimal.Truncate(precision)
	d.operations = append(d.operations, truncate)

	return d
}

func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	return d.decimal.UnmarshalJSON(decimalBytes)
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.decimal.MarshalJSON()
}

func (d *Decimal) UnmarshalBinary(data []byte) error {
	return d.decimal.UnmarshalBinary(data)
}

func (d Decimal) MarshalBinary() (data []byte, err error) {
	return d.decimal.MarshalBinary()
}

func (d *Decimal) Scan(value interface{}) error {
	return d.decimal.Scan(value)
}

func (d Decimal) Value() (driver.Value, error) {
	return d.decimal.Value()
}

func (d *Decimal) UnmarshalText(text []byte) error {
	return d.decimal.UnmarshalText(text)
}

func (d Decimal) MarshalText() (text []byte, err error) {
	return d.decimal.MarshalText()
}

func (d Decimal) GobEncode() ([]byte, error) {
	return d.decimal.GobEncode()
}

func (d *Decimal) GobDecode(data []byte) error {
	return d.decimal.GobDecode(data)
}

func (d Decimal) StringScaled(exp int32) string {
	return d.decimal.StringScaled(exp)
}

// func Min(first Decimal, rest ...Decimal) Decimal {}
// func Max(first Decimal, rest ...Decimal) Decimal {}
// func Sum(first Decimal, rest ...Decimal) Decimal {}
// func Avg(first Decimal, rest ...Decimal) Decimal {}
// func RescalePair(d1 Decimal, d2 Decimal) (Decimal, Decimal) {}
// func (d *NullDecimal) Scan(value interface{}}) error {}
// func (d NullDecimal) Value() (driver.Value, error) {}
// func (d *NullDecimal) UnmarshalJSON(decimalBytes []byte) error {}
// func (d NullDecimal) MarshalJSON() ([]byte, error) {}
// func (d Decimal) Atan() Decimal {}
// func (d Decimal) Sin() Decimal {}
// func (d Decimal) Cos() Decimal {}
// func (d Decimal) Tan() Decimal {}
