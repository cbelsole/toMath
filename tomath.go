package tomath

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

type (
	Decimal struct {
		parens  bool
		name    string
		decimal decimal.Decimal
		vars    string
		formula string
	}
	NullDecimal struct {
		name    string
		decimal decimal.NullDecimal
	}
)

func (d Decimal) SetName(name string) Decimal {
	d.name = name
	if d.vars == "" || d.vars == "?" {
		d.vars = name
	}
	return d
}

func (d Decimal) GetName() string {
	return d.name
}

func (d Decimal) Resolve() Decimal {
	return Decimal{
		name:    d.name,
		vars:    d.name,
		formula: d.String(),
		decimal: d.decimal,
	}
}

func (d Decimal) Decimal() decimal.Decimal {
	return d.decimal
}

func (d Decimal) Math() (string, string) {
	if d.name == "" {
		d.vars = fmt.Sprintf("%s = ?", d.vars)
	} else {
		d.vars = fmt.Sprintf("%s = %s", d.vars, d.name)
	}

	d.formula = fmt.Sprintf("%s = %s", d.formula, d.decimal)

	return d.vars, d.formula
}

func New(name string, value int64, exp int32) Decimal {
	d := decimal.New(value, exp)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}

func NewFromInt(name string, value int64) Decimal {
	d := decimal.NewFromInt(value)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}

func NewFromInt32(name string, value int32) Decimal {
	d := decimal.NewFromInt32(value)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}
func NewFromBigInt(name string, value *big.Int, exp int32) Decimal {
	d := decimal.NewFromBigInt(value, exp)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}
func NewFromString(name string, value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}, nil
}
func RequireFromString(name string, value string) Decimal {
	d := decimal.RequireFromString(value)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}
func NewFromFloat(name string, value float64) Decimal {
	d := decimal.NewFromFloat(value)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}
func NewFromFloat32(name string, value float32) Decimal {
	d := decimal.NewFromFloat32(value)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}
func NewFromFloatWithExponent(name string, value float64, exp int32) Decimal {
	d := decimal.NewFromFloatWithExponent(value, exp)
	return Decimal{name: name, decimal: d, vars: name, formula: d.String()}
}

func (d Decimal) Abs() Decimal {
	return Decimal{
		decimal: d.decimal.Abs(),
		vars:    fmt.Sprintf("abs(%s)", d.vars),
		formula: fmt.Sprintf("abs(%s)", d.formula),
	}
}

func (d Decimal) Add(d2 Decimal) Decimal {
	return Decimal{
		decimal: d.decimal.Add(d2.decimal),
		vars:    fmt.Sprintf("%s + %s", d.vars, d2.name),
		formula: fmt.Sprintf("%s + %s", d.formula, d2.String()),
		parens:  true,
	}
}

func (d Decimal) Sub(d2 Decimal) Decimal {
	return Decimal{
		decimal: d.decimal.Sub(d2.decimal),
		vars:    fmt.Sprintf("%s - %s", d.vars, d2.name),
		formula: fmt.Sprintf("%s - %s", d.formula, d2.String()),
		parens:  true,
	}
}

func (d Decimal) Neg() Decimal {
	return Decimal{
		decimal: d.decimal.Neg(),
		vars:    fmt.Sprintf("neg(%s)", d.vars),
		formula: fmt.Sprintf("neg(%s)", d.formula),
	}
}

func (d Decimal) Mul(d2 Decimal) Decimal {
	dec := Decimal{decimal: d.decimal.Mul(d2.decimal)}
	if d.parens {
		dec.vars = fmt.Sprintf("(%s) * %s", d.vars, d2.name)
		dec.formula = fmt.Sprintf("(%s) * %s", d.formula, d2.String())
	} else {
		dec.vars = fmt.Sprintf("%s * %s", d.vars, d2.name)
		dec.formula = fmt.Sprintf("%s * %s", d.formula, d2.String())
	}

	return dec
}

func (d Decimal) Shift(s int32) Decimal {
	return Decimal{
		decimal: d.decimal.Shift(s),
		vars:    fmt.Sprintf("shift(%d)(%s)", s, d.vars),
		formula: fmt.Sprintf("shift(%d)(%s)", s, d.formula),
	}
}

func (d Decimal) Div(d2 Decimal) Decimal {
	dec := Decimal{decimal: d.decimal.Div(d2.decimal)}
	if d.parens {
		dec.vars = fmt.Sprintf("(%s) / %s", d.vars, d2.name)
		dec.formula = fmt.Sprintf("(%s) / %s", d.formula, d2.String())
	} else {
		dec.vars = fmt.Sprintf("%s / %s", d.vars, d2.name)
		dec.formula = fmt.Sprintf("%s / %s", d.formula, d2.String())
	}

	return dec
}

func (d Decimal) QuoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
	d3, d4 := d.decimal.QuoRem(d2.decimal, precision)

	var vars, formula string
	if d.parens {
		vars = fmt.Sprintf("quoRem(%d)((%s) / %s)", precision, d.vars, d2.name)
		formula = fmt.Sprintf("quoRem(%d)((%s) / %s)", precision, d.formula, d2.String())
	} else {
		vars = fmt.Sprintf("quoRem(%d)(%s / %s)", precision, d.vars, d2.name)
		formula = fmt.Sprintf("quoRem(%d)(%s / %s)", precision, d.formula, d2.String())
	}

	return Decimal{name: d.name + d2.name + "Quotient", decimal: d3, vars: vars, formula: formula},
		Decimal{name: d.name + d2.name + "Remainder", decimal: d4, vars: vars, formula: formula}
}

func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
	dec := Decimal{decimal: d.decimal.DivRound(d2.decimal, precision)}
	if d.parens {
		dec.vars = fmt.Sprintf("divRound(%d)((%s) / %s)", precision, d.vars, d2.name)
		dec.formula = fmt.Sprintf("divRound(%d)((%s) / %s)", precision, d.formula, d2.String())
	} else {
		dec.vars = fmt.Sprintf("divRound(%d)(%s / %s)", precision, d.vars, d2.name)
		dec.formula = fmt.Sprintf("divRound(%d)(%s / %s)", precision, d.formula, d2.String())
	}

	return dec
}

func (d Decimal) Mod(d2 Decimal) Decimal {
	dec := Decimal{decimal: d.decimal.Mod(d2.decimal)}
	if d.parens {
		dec.vars = fmt.Sprintf("(%s) %% %s", d.vars, d2.name)
		dec.formula = fmt.Sprintf("(%s) %% %s", d.formula, d2.String())
	} else {
		dec.vars = fmt.Sprintf("%s %% %s", d.vars, d2.name)
		dec.formula = fmt.Sprintf("%s %% %s", d.formula, d2.String())
	}

	return dec
}

func (d Decimal) Pow(d2 Decimal) Decimal {
	dec := Decimal{decimal: d.decimal.Pow(d2.decimal)}
	if d.parens {
		dec.vars = fmt.Sprintf("(%s)^%s", d.vars, d2.name)
		dec.formula = fmt.Sprintf("(%s)^%s", d.formula, d2.String())
	} else {
		dec.vars = fmt.Sprintf("%s^%s", d.vars, d2.name)
		dec.formula = fmt.Sprintf("%s^%s", d.formula, d2.String())
	}

	return dec
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
	return Decimal{
		decimal: d.decimal.Round(places),
		vars:    fmt.Sprintf("round(%d)(%s)", places, d.vars),
		formula: fmt.Sprintf("round(%d)(%s)", places, d.formula),
	}
}

func (d Decimal) RoundBank(places int32) Decimal {
	return Decimal{
		decimal: d.decimal.RoundBank(places),
		vars:    fmt.Sprintf("roundBank(%d)(%s)", places, d.vars),
		formula: fmt.Sprintf("roundBank(%d)(%s)", places, d.formula),
	}
}

func (d Decimal) RoundCash(interval uint8) Decimal {
	return Decimal{
		decimal: d.decimal.RoundCash(interval),
		vars:    fmt.Sprintf("roundCash(%d)(%s)", interval, d.vars),
		formula: fmt.Sprintf("roundCash(%d)(%s)", interval, d.formula),
	}
}

func (d Decimal) Floor() Decimal {
	return Decimal{
		decimal: d.decimal.Floor(),
		vars:    fmt.Sprintf("floor(%s)", d.vars),
		formula: fmt.Sprintf("floor(%s)", d.formula),
	}
}

func (d Decimal) Ceil() Decimal {
	return Decimal{
		decimal: d.decimal.Ceil(),
		vars:    fmt.Sprintf("ceil(%s)", d.vars),
		formula: fmt.Sprintf("ceil(%s)", d.formula),
	}
}

func (d Decimal) Truncate(precision int32) Decimal {
	return Decimal{
		decimal: d.decimal.Truncate(precision),
		vars:    fmt.Sprintf("truncate(%d)(%s)", precision, d.vars),
		formula: fmt.Sprintf("truncate(%d)(%s)", precision, d.formula),
	}
}

func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	if err := d.decimal.UnmarshalJSON(decimalBytes); err != nil {
		return err
	}
	d.vars = "?"
	d.formula = d.String()

	return nil
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.decimal.MarshalJSON()
}

func (d *Decimal) UnmarshalBinary(data []byte) error {
	if err := d.decimal.UnmarshalBinary(data); err != nil {
		return err
	}
	d.vars = "?"
	d.formula = d.String()
	return nil
}

func (d Decimal) MarshalBinary() (data []byte, err error) {
	return d.decimal.MarshalBinary()
}

func (d *Decimal) Scan(value interface{}) error {
	if err := d.decimal.Scan(value); err != nil {
		return err
	}
	d.vars = "?"
	d.formula = d.String()
	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	return d.decimal.Value()
}

func (d *Decimal) UnmarshalText(text []byte) error {
	if err := d.decimal.UnmarshalText(text); err != nil {
		return err
	}
	d.vars = "?"
	d.formula = d.String()
	return nil
}

func (d Decimal) MarshalText() (text []byte, err error) {
	return d.decimal.MarshalText()
}

func (d Decimal) GobEncode() ([]byte, error) {
	return d.decimal.GobEncode()
}

func (d *Decimal) GobDecode(data []byte) error {
	if err := d.decimal.UnmarshalBinary(data); err != nil {
		return err
	}
	d.vars = "?"
	d.formula = d.String()
	return d.decimal.GobDecode(data)
}

func (d Decimal) StringScaled(exp int32) string {
	return d.decimal.StringScaled(exp)
}

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

func (d *NullDecimal) Scan(value interface{}) error {
	if err := d.decimal.Scan(value); err != nil {
		return err
	}
	d.name = "?"
	return nil
}
func (d NullDecimal) Value() (driver.Value, error) {
	return d.decimal.Value()
}
func (d *NullDecimal) UnmarshalJSON(decimalBytes []byte) error {
	if err := d.decimal.UnmarshalJSON(decimalBytes); err != nil {
		return err
	}
	d.name = "?"
	return nil
}
func (d NullDecimal) MarshalJSON() ([]byte, error) {
	return d.decimal.MarshalJSON()
}

// func (d Decimal) Atan() Decimal {}
// func (d Decimal) Sin() Decimal {}
// func (d Decimal) Cos() Decimal {}
// func (d Decimal) Tan() Decimal {}
