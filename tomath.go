package tomath

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	abs operator = iota
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
	min
	max
	sum
	avg
)

type (
	operator  int
	operation struct {
		operator operator
		operands []Decimal
	}
	Decimal struct {
		name      string
		precision int
		decimal   decimal.Decimal
		operation *operation
	}
)

// TODO: Implement me to get rid of the ?
func (d Decimal) Result(name string) Decimal {
	d.name = name
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

	if d.name == "" {
		vars = fmt.Sprintf("%s = ?", vars)
	} else {
		vars = fmt.Sprintf("%s = %s", vars, d.name)
	}

	formula = fmt.Sprintf("%s = %s", formula, d.decimal)

	return vars, formula
}

func insertSliceAt(slice []interface{}, at int, insertables ...interface{}) []interface{} {
	for i := 0; i < len(insertables)-1; i++ {
		slice = append(slice, interface{}(nil))
	}

	copy(slice[at+len(insertables)-1:], slice[at:])

	for i, insert := range insertables {
		slice[at+i] = insert
	}

	return slice
}

func shouldUseParensAt(slice []interface{}, at int) bool {
	if len(slice) > at+1 {
		if operator, ok := slice[at+1].(string); ok && operator == " * " || operator == " / " {
			return true
		}
	}
	return false
}

func (d Decimal) math() (string, string) {
	vars := []interface{}{d}
	formula := []interface{}{d}

	for {
		breakLoop := true

		for i := 0; i < len(vars); i++ {
			switch c := vars[i].(type) {
			case Decimal:
				breakLoop = false

				if c.operation != nil {
					var insertables []interface{}

					switch c.operation.operator {
					case abs:
						insertables = []interface{}{"abs(", c.operation.operands[0], ")"}
					case add:
						if shouldUseParensAt(vars, i) {
							insertables = []interface{}{"(", c.operation.operands[0], " + ", c.operation.operands[1], ")"}
						} else {
							insertables = []interface{}{c.operation.operands[0], " + ", c.operation.operands[1]}
						}

					case sub:
						insertables = []interface{}{c.operation.operands[0], " - ", c.operation.operands[1]}
					case neg:
						insertables = []interface{}{"neg(", c.operation.operands[0], ")"}
					case mul:
						if shouldUseParensAt(vars, i) {
							insertables = []interface{}{"(", c.operation.operands[0], " * ", c.operation.operands[1], ")"}
						} else {
							insertables = []interface{}{c.operation.operands[0], " * ", c.operation.operands[1]}
						}

					case shift:
						insertables = []interface{}{"shift(", strconv.Itoa(c.operation.operands[0].precision), ")(", c.operation.operands[0], ")"}
					case div:
						if shouldUseParensAt(vars, i) {
							insertables = []interface{}{"(", c.operation.operands[0], " / ", c.operation.operands[1], ")"}
						} else {
							insertables = []interface{}{c.operation.operands[0], " / ", c.operation.operands[1]}
						}

					case quoRem:
						insertables = []interface{}{"quoRem(", strconv.Itoa(c.operation.operands[0].precision), ")(", c.operation.operands[0], " / ", c.operation.operands[1], ")"}
					case divRound:
						insertables = []interface{}{"divRound(", strconv.Itoa(c.operation.operands[0].precision), ")(", c.operation.operands[0], " / ", c.operation.operands[1], ")"}
					case mod:
						if shouldUseParensAt(vars, i) {
							insertables = []interface{}{"(", c.operation.operands[0], " % ", c.operation.operands[1], ")"}
						} else {
							insertables = []interface{}{c.operation.operands[0], " % ", c.operation.operands[1]}
						}

					case pow:
						if shouldUseParensAt(vars, i) {
							insertables = []interface{}{"(", c.operation.operands[0], "^", c.operation.operands[1], ")"}
						} else {
							insertables = []interface{}{c.operation.operands[0], "^", c.operation.operands[1]}
						}

					case round:
						insertables = []interface{}{"round(", strconv.Itoa(c.operation.operands[0].precision), ")(", c.operation.operands[0], ")"}
					case roundBank:
						insertables = []interface{}{"roundBank(", strconv.Itoa(c.operation.operands[0].precision), ")(", c.operation.operands[0], ")"}
					case roundCash:
						insertables = []interface{}{"roundCash(", strconv.Itoa(c.operation.operands[0].precision), ")(", c.operation.operands[0], ")"}
					case floor:
						insertables = []interface{}{"floor(", c.operation.operands[0], ")"}
					case ceil:
						insertables = []interface{}{"ceil(", c.operation.operands[0], ")"}
					case truncate:
						insertables = []interface{}{"truncate(", strconv.Itoa(c.operation.operands[0].precision), ")(", c.operation.operands[0], ")"}
					case min:
						insertables = make([]interface{}, (len(c.operation.operands)*2)+1)
						insertables[0] = "min("

						for i, operand := range c.operation.operands {
							insertables[(i*2)+1] = operand
							if i != len(c.operation.operands)-1 {
								insertables[(i*2)+2] = ", "
							}
						}
						insertables[(len(c.operation.operands) * 2)] = ")"
					case max:
						insertables = make([]interface{}, (len(c.operation.operands)*2)+1)
						insertables[0] = "max("

						for i, operand := range c.operation.operands {
							insertables[(i*2)+1] = operand
							if i != len(c.operation.operands)-1 {
								insertables[(i*2)+2] = ", "
							}
						}
						insertables[(len(c.operation.operands) * 2)] = ")"
					case sum:
						insertables = make([]interface{}, (len(c.operation.operands)*2)+1)
						insertables[0] = "sum("

						for i, operand := range c.operation.operands {
							insertables[(i*2)+1] = operand
							if i != len(c.operation.operands)-1 {
								insertables[(i*2)+2] = ", "
							}
						}
						insertables[(len(c.operation.operands) * 2)] = ")"
					case avg:
						insertables = make([]interface{}, (len(c.operation.operands)*2)+1)
						insertables[0] = "avg("

						for i, operand := range c.operation.operands {
							insertables[(i*2)+1] = operand
							if i != len(c.operation.operands)-1 {
								insertables[(i*2)+2] = ", "
							}
						}
						insertables[(len(c.operation.operands) * 2)] = ")"
					}

					vars = insertSliceAt(vars, i, insertables...)
					formula = insertSliceAt(formula, i, insertables...)
					i += len(insertables) - 1
				} else {
					vars = insertSliceAt(vars, i, c.name)
					formula = insertSliceAt(formula, i, c.decimal.String())
				}
			}
		}

		if breakLoop {
			break
		}
	}

	outVars := make([]string, len(vars))
	for _, v := range vars {
		outVars = append(outVars, v.(string))
	}

	outFormula := make([]string, len(formula))
	for _, f := range formula {
		outFormula = append(outFormula, f.(string))
	}

	return strings.Join(outVars, ""), strings.Join(outFormula, "")
}

func New(name string, value int64, exp int32) Decimal {
	d := decimal.New(value, exp)
	return Decimal{name: name, decimal: d}
}

func NewFromInt(name string, value int64) Decimal {
	d := decimal.NewFromInt(value)
	return Decimal{name: name, decimal: d}
}

func NewFromInt32(name string, value int32) Decimal {
	d := decimal.NewFromInt32(value)
	return Decimal{name: name, decimal: d}
}
func NewFromBigInt(name string, value *big.Int, exp int32) Decimal {
	d := decimal.NewFromBigInt(value, exp)
	return Decimal{name: name, decimal: d}
}
func NewFromString(name string, value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{name: name, decimal: d}, nil
}
func RequireFromString(name string, value string) Decimal {
	d := decimal.RequireFromString(value)
	return Decimal{name: name, decimal: d}
}
func NewFromFloat(name string, value float64) Decimal {
	d := decimal.NewFromFloat(value)
	return Decimal{name: name, decimal: d}
}
func NewFromFloat32(name string, value float32) Decimal {
	d := decimal.NewFromFloat32(value)
	return Decimal{name: name, decimal: d}
}
func NewFromFloatWithExponent(name string, value float64, exp int32) Decimal {
	d := decimal.NewFromFloatWithExponent(value, exp)
	return Decimal{name: name, decimal: d}
}

func (d Decimal) Abs() Decimal {
	return Decimal{
		decimal:   d.decimal.Abs(),
		operation: &operation{abs, []Decimal{d}},
	}
}

func (d Decimal) Add(d2 Decimal) Decimal {
	return Decimal{
		decimal:   d.decimal.Add(d2.decimal),
		operation: &operation{add, []Decimal{d, d2}},
	}
}

func (d Decimal) Sub(d2 Decimal) Decimal {
	return Decimal{
		decimal:   d.decimal.Sub(d2.decimal),
		operation: &operation{sub, []Decimal{d, d2}},
	}
}

func (d Decimal) Neg() Decimal {
	return Decimal{
		decimal:   d.decimal.Neg(),
		operation: &operation{neg, []Decimal{d}},
	}
}

func (d Decimal) Mul(d2 Decimal) Decimal {
	return Decimal{
		decimal:   d.decimal.Mul(d2.decimal),
		operation: &operation{mul, []Decimal{d, d2}},
	}
}

func (d Decimal) Shift(s int32) Decimal {
	d.precision = int(s)
	return Decimal{
		decimal:   d.decimal.Shift(s),
		operation: &operation{shift, []Decimal{d}},
	}
}

func (d Decimal) Div(d2 Decimal) Decimal {
	return Decimal{
		decimal:   d.decimal.Div(d2.decimal),
		operation: &operation{div, []Decimal{d, d2}},
	}
}

func (d Decimal) QuoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
	d3, d4 := d.decimal.QuoRem(d2.decimal, precision)
	d.precision = int(precision)

	return Decimal{
			name:      d.name + d2.name + "Quotient",
			operation: &operation{quoRem, []Decimal{d, d2}},
			decimal:   d3,
		}, Decimal{
			name:      d.name + d2.name + "Remainder",
			operation: &operation{quoRem, []Decimal{d, d2}},
			decimal:   d4,
		}
}

func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
	d.precision = int(precision)

	return Decimal{
		decimal:   d.decimal.DivRound(d2.decimal, precision),
		operation: &operation{divRound, []Decimal{d, d2}},
	}
}

func (d Decimal) Mod(d2 Decimal) Decimal {
	return Decimal{
		decimal:   d.decimal.Mod(d2.decimal),
		operation: &operation{mod, []Decimal{d, d2}},
	}
}

func (d Decimal) Pow(d2 Decimal) Decimal {
	return Decimal{
		decimal:   d.decimal.Pow(d2.decimal),
		operation: &operation{pow, []Decimal{d, d2}},
	}
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
	d.precision = int(places)

	return Decimal{
		decimal:   d.decimal.Round(places),
		operation: &operation{round, []Decimal{d}},
	}
}

func (d Decimal) RoundBank(places int32) Decimal {
	d.precision = int(places)

	return Decimal{
		decimal:   d.decimal.RoundBank(places),
		operation: &operation{roundBank, []Decimal{d}},
	}
}

func (d Decimal) RoundCash(interval uint8) Decimal {
	d.precision = int(interval)

	return Decimal{
		decimal:   d.decimal.RoundCash(interval),
		operation: &operation{roundCash, []Decimal{d}},
	}
}

func (d Decimal) Floor() Decimal {
	return Decimal{
		decimal:   d.decimal.Floor(),
		operation: &operation{floor, []Decimal{d}},
	}
}

func (d Decimal) Ceil() Decimal {
	return Decimal{
		decimal:   d.decimal.Ceil(),
		operation: &operation{ceil, []Decimal{d}},
	}
}

func (d Decimal) Truncate(precision int32) Decimal {
	d.precision = int(precision)

	return Decimal{
		decimal:   d.decimal.Truncate(precision),
		operation: &operation{truncate, []Decimal{d}},
	}
}

// func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
// 	return d.decimal.UnmarshalJSON(decimalBytes)
// }

// func (d Decimal) MarshalJSON() ([]byte, error) {
// 	return d.decimal.MarshalJSON()
// }

// func (d *Decimal) UnmarshalBinary(data []byte) error {
// 	return d.decimal.UnmarshalBinary(data)
// }

// func (d Decimal) MarshalBinary() (data []byte, err error) {
// 	return d.decimal.MarshalBinary()
// }

// func (d *Decimal) Scan(value interface{}) error {
// 	return d.decimal.Scan(value)
// }

// func (d Decimal) Value() (driver.Value, error) {
// 	return d.decimal.Value()
// }

// func (d *Decimal) UnmarshalText(text []byte) error {
// 	return d.decimal.UnmarshalText(text)
// }

// func (d Decimal) MarshalText() (text []byte, err error) {
// 	return d.decimal.MarshalText()
// }

// func (d Decimal) GobEncode() ([]byte, error) {
// 	return d.decimal.GobEncode()
// }

// func (d *Decimal) GobDecode(data []byte) error {
// 	return d.decimal.GobDecode(data)
// }

// func (d Decimal) StringScaled(exp int32) string {
// 	return d.decimal.StringScaled(exp)
// }

func Min(first Decimal, rest ...Decimal) Decimal {
	newRest := make([]decimal.Decimal, len(rest))
	for i, r := range rest {
		newRest[i] = r.decimal
	}

	return Decimal{
		decimal:   decimal.Min(first.decimal, newRest...),
		operation: &operation{min, append([]Decimal{first}, rest...)},
	}
}

func Max(first Decimal, rest ...Decimal) Decimal {
	newRest := make([]decimal.Decimal, len(rest))
	for i, r := range rest {
		newRest[i] = r.decimal
	}

	return Decimal{
		decimal:   decimal.Max(first.decimal, newRest...),
		operation: &operation{max, append([]Decimal{first}, rest...)},
	}
}

func Sum(first Decimal, rest ...Decimal) Decimal {
	newRest := make([]decimal.Decimal, len(rest))
	for i, r := range rest {
		newRest[i] = r.decimal
	}

	return Decimal{
		decimal:   decimal.Sum(first.decimal, newRest...),
		operation: &operation{sum, append([]Decimal{first}, rest...)},
	}
}

func Avg(first Decimal, rest ...Decimal) Decimal {
	newRest := make([]decimal.Decimal, len(rest))
	for i, r := range rest {
		newRest[i] = r.decimal
	}

	return Decimal{
		decimal:   decimal.Avg(first.decimal, newRest...),
		operation: &operation{avg, append([]Decimal{first}, rest...)},
	}
}

// func RescalePair(d1 Decimal, d2 Decimal) (Decimal, Decimal) {}
// func (d *NullDecimal) Scan(value interface{}}) error {}
// func (d NullDecimal) Value() (driver.Value, error) {}
// func (d *NullDecimal) UnmarshalJSON(decimalBytes []byte) error {}
// func (d NullDecimal) MarshalJSON() ([]byte, error) {}
// func (d Decimal) Atan() Decimal {}
// func (d Decimal) Sin() Decimal {}
// func (d Decimal) Cos() Decimal {}
// func (d Decimal) Tan() Decimal {}
