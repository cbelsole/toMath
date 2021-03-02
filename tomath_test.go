package tomath

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	d := Decimal{}
	assert.Equal(t, "0", d.String())
	formula, vars := d.Math()
	assert.Equal(t, "? = ?", formula)
	assert.Equal(t, "0 = 0", vars)
}

func TestNew(t *testing.T) {
	d := New(0, 0)
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromInt(t *testing.T) {
	d := NewFromInt(0)
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromInt32(t *testing.T) {
	d := NewFromInt32(0)
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromBigInt(t *testing.T) {
	d := NewFromBigInt(big.NewInt(0), 0)
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromString(t *testing.T) {
	d, err := NewFromString("0")
	assert.NoError(t, err)
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestRequireFromString(t *testing.T) {
	d := RequireFromString("0")
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromFloat(t *testing.T) {
	d := NewFromFloat(0)
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromFloat32(t *testing.T) {
	d := NewFromFloat32(0)
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromFloatWithExponent(t *testing.T) {
	d := NewFromFloatWithExponent(0, 0)
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromDecimal(t *testing.T) {
	d := NewFromDecimal(decimal.NewFromFloat(1))
	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "1 = 1", formula)
}

func TestNewWithName(t *testing.T) {
	d := NewWithName("var1", 0, 0)
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromIntWithName(t *testing.T) {
	d := NewFromIntWithName("var1", 0)
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromInt32WithName(t *testing.T) {
	d := NewFromInt32WithName("var1", 0)
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromBigIntWithName(t *testing.T) {
	d := NewFromBigIntWithName("var1", big.NewInt(0), 0)
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromStringWithName(t *testing.T) {
	d, err := NewFromStringWithName("var1", "0")
	assert.NoError(t, err)
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestRequireFromStringWithName(t *testing.T) {
	d := RequireFromStringWithName("var1", "0")
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromFloatWithName(t *testing.T) {
	d := NewFromFloatWithName("var1", 0)
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromFloat32WithName(t *testing.T) {
	d := NewFromFloat32WithName("var1", 0)
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromFloatWithExponentWithName(t *testing.T) {
	d := NewFromFloatWithExponentWithName("var1", 0, 0)
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "0 = 0", formula)
}

func TestNewFromDecimalWithName(t *testing.T) {
	d := NewFromDecimalWithName("var1", decimal.NewFromFloat(1))
	vars, formula := d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "1 = 1", formula)
}

func TestAbs(t *testing.T) {
	d := NewWithName("var1", -1, 0).Abs()
	vars, formula := d.Math()
	assert.Equal(t, "abs(var1) = ?", vars)
	assert.Equal(t, "abs(-1) = 1", formula)

	d2 := d.Abs()
	vars, formula = d2.Math()
	assert.Equal(t, "abs(abs(var1)) = ?", vars)
	assert.Equal(t, "abs(abs(-1)) = 1", formula)
}

func TestAdd(t *testing.T) {
	d := NewWithName("var1", -1, 0).Add(NewWithName("var2", 0, 0))
	vars, formula := d.Math()
	assert.Equal(t, "var1 + var2 = ?", vars)
	assert.Equal(t, "-1 + 0 = -1", formula)

	d2 := NewWithName("var3", -1, 0).Add(NewWithName("var4", 0, 0))
	d3 := d.Add(d2)

	vars, formula = d3.Math()
	assert.Equal(t, "var1 + var2 + var3 + var4 = ?", vars)
	assert.Equal(t, "-1 + 0 + -1 + 0 = -2", formula)
}

func TestSub(t *testing.T) {
	d := NewWithName("var1", -1, 0).Sub(NewWithName("var2", 0, 0))
	vars, formula := d.Math()
	assert.Equal(t, "var1 - var2 = ?", vars)
	assert.Equal(t, "-1 - 0 = -1", formula)

	d2 := NewWithName("var3", -1, 0).Sub(NewWithName("var4", 0, 0))
	d3 := d.Sub(d2)

	vars, formula = d3.Math()
	assert.Equal(t, "var1 - var2 - var3 - var4 = ?", vars)
	assert.Equal(t, "-1 - 0 - -1 - 0 = 0", formula)
}

func TestNeg(t *testing.T) {
	d := NewWithName("var1", 1, 0).Neg()
	vars, formula := d.Math()
	assert.Equal(t, "neg(var1) = ?", vars)
	assert.Equal(t, "neg(1) = -1", formula)

	d2 := d.Neg()
	vars, formula = d2.Math()
	assert.Equal(t, "neg(neg(var1)) = ?", vars)
	assert.Equal(t, "neg(neg(1)) = 1", formula)
}

func TestMul(t *testing.T) {
	d := NewWithName("var1", 1, 0).Mul(NewWithName("var2", 2, 0))
	vars, formula := d.Math()
	assert.Equal(t, "var1 * var2 = ?", vars)
	assert.Equal(t, "1 * 2 = 2", formula)
}

// func TestShift(t *testing.T) {
// 	d := NewWithName("var1", 1, 0).Shift(1)
// 	vars, formula := d.Math()
// 	assert.Equal(t, "shift(1)(var1) = ?", vars)
// 	assert.Equal(t, "shift(1)(1) = 10", formula)
// }

// func TestDiv(t *testing.T) {
// 	d := NewWithName("var1", 4, 0).Div(NewWithName("var2", 2, 0))
// 	vars, formula := d.Math()
// 	assert.Equal(t, "var1 / var2 = ?", vars)
// 	assert.Equal(t, "4 / 2 = 2", formula)
// }

// func TestDivRound(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 4.333).DivRound(NewFromFloatWithName("var2", 2.7), 3)
// 	vars, formula := d.Math()
// 	assert.Equal(t, "divRound(3)(var1 / var2) = ?", vars)
// 	assert.Equal(t, "divRound(3)(4.333 / 2.7) = 1.605", formula)
// }

// func TestQuoRem(t *testing.T) {
// 	d1, d2 := NewFromFloatWithName("var1", 4.333).QuoRem(NewFromFloatWithName("var2", 2.7), 3)

// 	vars, formula := d1.Math()
// 	assert.Equal(t, "quoRem(3)(var1 / var2) = var1var2Quotient", vars)
// 	assert.Equal(t, "quoRem(3)(4.333 / 2.7) = 1.604", formula)

// 	vars, formula = d2.Math()
// 	assert.Equal(t, "quoRem(3)(var1 / var2) = var1var2Remainder", vars)
// 	assert.Equal(t, "quoRem(3)(4.333 / 2.7) = 0.0022", formula)
// }

// func TestMod(t *testing.T) {
// 	d := NewWithName("var1", 4, 0).Mod(NewWithName("var2", 2, 0))
// 	vars, formula := d.Math()
// 	assert.Equal(t, "var1 % var2 = ?", vars)
// 	assert.Equal(t, "4 % 2 = 0", formula)
// }

// func TestPow(t *testing.T) {
// 	d := NewWithName("var1", 4, 0).Pow(NewWithName("var2", 2, 0))
// 	vars, formula := d.Math()
// 	assert.Equal(t, "var1^var2 = ?", vars)
// 	assert.Equal(t, "4^2 = 16", formula)
// }

// func TestRound(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 4.333).Round(2)
// 	vars, formula := d.Math()
// 	assert.Equal(t, "round(2)(var1) = ?", vars)
// 	assert.Equal(t, "round(2)(4.333) = 4.33", formula)
// }

// func TestRoundBank(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 4.333).RoundBank(2)
// 	vars, formula := d.Math()
// 	assert.Equal(t, "roundBank(2)(var1) = ?", vars)
// 	assert.Equal(t, "roundBank(2)(4.333) = 4.33", formula)
// }

// func TestRoundCash(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 4.333).RoundCash(5)
// 	vars, formula := d.Math()
// 	assert.Equal(t, "roundCash(5)(var1) = ?", vars)
// 	assert.Equal(t, "roundCash(5)(4.333) = 4.35", formula)
// }

// func TestFloor(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 4.333).Floor()
// 	vars, formula := d.Math()
// 	assert.Equal(t, "floor(var1) = ?", vars)
// 	assert.Equal(t, "floor(4.333) = 4", formula)
// }

// func TestCeil(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 4.333).Ceil()
// 	vars, formula := d.Math()
// 	assert.Equal(t, "ceil(var1) = ?", vars)
// 	assert.Equal(t, "ceil(4.333) = 5", formula)
// }

// func TestTruncate(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 4.333).Truncate(0)
// 	vars, formula := d.Math()
// 	assert.Equal(t, "truncate(0)(var1) = ?", vars)
// 	assert.Equal(t, "truncate(0)(4.333) = 4", formula)
// }

// func TestSetName(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 1).Add(NewFromFloatWithName("var2", 1)).SetName("var3")
// 	vars, formula := d.Math()
// 	assert.Equal(t, "var1 + var2 = var3", vars)
// 	assert.Equal(t, "1 + 1 = 2", formula)
// }

// func TestResolve(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 1).Add(NewFromFloatWithName("var2", 1)).SetName("var3")
// 	vars, formula := d.Math()
// 	assert.Equal(t, "var1 + var2 = var3", vars)
// 	assert.Equal(t, "1 + 1 = 2", formula)

// 	d = d.Resolve()
// 	vars, formula = d.Math()
// 	assert.Equal(t, "var3 = var3", vars)
// 	assert.Equal(t, "2 = 2", formula)
// }

// func TestResolveTo(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 1).Add(NewFromFloatWithName("var2", 1)).SetName("var3")
// 	vars, formula := d.Math()
// 	assert.Equal(t, "var1 + var2 = var3", vars)
// 	assert.Equal(t, "1 + 1 = 2", formula)

// 	d = d.ResolveTo("var4")
// 	vars, formula = d.Math()
// 	assert.Equal(t, "var4 = var4", vars)
// 	assert.Equal(t, "2 = 2", formula)
// }

// func TestComplexExample(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 1.1).
// 		Round(1).
// 		Add(NewFromFloatWithName("var2", 1)).
// 		Add(NewFromFloatWithName("var2", 1)).
// 		Div(NewFromFloatWithName("var3", 2)).
// 		Mul(NewFromFloatWithName("OneHundred", 100).Add(NewFromFloatWithName("var4", 3))).
// 		SetName("var5")

// 	vars, formula := d.Math()
// 	assert.Equal(t, "(round(1)(var1) + var2 + var2) / var3 * (OneHundred + var4) = var5", vars)
// 	assert.Equal(t, "(round(1)(1.1) + 1 + 1) / 2 * (100 + 3) = 159.65", formula)

// 	// the Math() function is idempotent
// 	vars, formula = d.Math()
// 	assert.Equal(t, "(round(1)(var1) + var2 + var2) / var3 * (OneHundred + var4) = var5", vars)
// 	assert.Equal(t, "(round(1)(1.1) + 1 + 1) / 2 * (100 + 3) = 159.65", formula)
// }

// func TestMin(t *testing.T) {
// 	d := Min(NewFromFloatWithName("var1", 1), NewFromFloatWithName("var2", 2), NewFromFloatWithName("var3", 100))

// 	vars, formula := d.Math()
// 	assert.Equal(t, "min(var1, var2, var3) = ?", vars)
// 	assert.Equal(t, "min(1, 2, 100) = 1", formula)
// }

// func TestMax(t *testing.T) {
// 	d := Max(NewFromFloatWithName("var1", 1), NewFromFloatWithName("var2", 2), NewFromFloatWithName("var3", 100))

// 	vars, formula := d.Math()
// 	assert.Equal(t, "max(var1, var2, var3) = ?", vars)
// 	assert.Equal(t, "max(1, 2, 100) = 100", formula)
// }

// func TestSum(t *testing.T) {
// 	d := Sum(NewFromFloatWithName("var1", 1), NewFromFloatWithName("var2", 2), NewFromFloatWithName("var3", 100))

// 	vars, formula := d.Math()
// 	assert.Equal(t, "sum(var1, var2, var3) = ?", vars)
// 	assert.Equal(t, "sum(1, 2, 100) = 103", formula)
// }

// func TestAvg(t *testing.T) {
// 	d := Avg(NewFromFloatWithName("var1", 1), NewFromFloatWithName("var2", 2), NewFromFloatWithName("var3", 100))

// 	vars, formula := d.Math()
// 	assert.Equal(t, "avg(var1, var2, var3) = ?", vars)
// 	assert.Equal(t, "avg(1, 2, 100) = 34.3333333333333333", formula)
// }

// func TestRescalePair(t *testing.T) {
// 	d1, d2 := RescalePair(NewWithName("var1", 111111, -5), NewWithName("var2", 2111, -3))

// 	vars, formula := d1.Math()
// 	assert.Equal(t, "var1 = var1", vars)
// 	assert.Equal(t, "1.11111 = 1.11111", formula)

// 	vars, formula = d2.Math()
// 	assert.Equal(t, "var2 = var2", vars)
// 	assert.Equal(t, "2.111 = 2.111", formula)
// }

// func TestUnmarshalJSON(t *testing.T) {
// 	d := &Decimal{}
// 	require.NoError(t, d.UnmarshalJSON([]byte("123.123")))
// 	require.Equal(t, "123.123", d.String())
// }

// func TestMarshalJSON(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 123.123)
// 	b, err := d.MarshalJSON()
// 	require.NoError(t, err)
// 	require.Equal(t, `"123.123"`, string(b))
// }

// func TestBinary(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 123.123)
// 	b, err := d.MarshalBinary()
// 	require.NoError(t, err)

// 	d2 := &Decimal{}
// 	require.NoError(t, d2.UnmarshalBinary(b))
// 	require.Equal(t, "123.123", d.String())
// }

// func TestScan(t *testing.T) {
// 	d := Decimal{}
// 	d2 := NewFromFloatWithName("var1", 54.33)
// 	require.NoError(t, d.Scan(54.33))
// 	require.Equal(t, d.String(), d2.String())
// }

// func TestValue(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 54.33)
// 	SetName, err := d.Value()
// 	require.NoError(t, err)
// 	require.Equal(t, d.String(), SetName.(string))
// }

// func TestUnmarshalText(t *testing.T) {
// 	d := &Decimal{}
// 	require.NoError(t, d.UnmarshalText([]byte("123.123")))
// 	require.Equal(t, "123.123", d.String())
// }

// func TestMarshalText(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 123.123)
// 	b, err := d.MarshalText()
// 	require.NoError(t, err)
// 	require.Equal(t, "123.123", string(b))
// }

// func TestGobEncode(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 123.123)
// 	b, err := d.GobEncode()
// 	require.NoError(t, err)

// 	d2 := &Decimal{}
// 	require.NoError(t, d2.GobDecode(b))
// 	require.Equal(t, "123.123", d.String())
// }

// func TestStringScaled(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 123.123)
// 	assert.Equal(t, "123.1", d.StringScaled(-1))

// }

// func TestUnknownName(t *testing.T) {
// 	d := Decimal{}
// 	require.NoError(t, d.UnmarshalJSON([]byte("123.123")))
// 	require.Equal(t, "123.123", d.String())

// 	vars, formula := d.Math()
// 	assert.Equal(t, "? = ?", vars)
// 	assert.Equal(t, "123.123 = 123.123", formula)

// 	d = d.SetName("var1")
// 	vars, formula = d.Math()
// 	assert.Equal(t, "var1 = var1", vars)
// 	assert.Equal(t, "123.123 = 123.123", formula)
// }

// func TestNullDecimalScan(t *testing.T) {
// 	d := NullDecimal{}
// 	d2 := NewFromFloatWithName("var1", 54.33)
// 	require.NoError(t, d.Scan(54.33))
// 	require.Equal(t, d.decimal.Decimal.String(), d2.String())
// }

// func TestNullDecimalValue(t *testing.T) {
// 	d := NewFromFloatWithName("var1", 54.33)
// 	SetName, err := d.Value()
// 	require.NoError(t, err)
// 	require.Equal(t, d.String(), SetName.(string))
// }

// func TestNullDecimalJSON(t *testing.T) {
// 	d := &NullDecimal{}
// 	require.NoError(t, d.UnmarshalJSON([]byte("123.123")))
// 	require.Equal(t, "123.123", d.Decimal().String())

// 	b, err := d.MarshalJSON()
// 	require.NoError(t, err)
// 	require.Equal(t, `"123.123"`, string(b))

// }

// func TestAtan(t *testing.T) {
// 	d := NewWithName("var1", 1, 0).Atan()
// 	vars, formula := d.Math()
// 	assert.Equal(t, "atan(var1) = ?", vars)
// 	assert.Equal(t, "atan(1) = 0.78539816339744833061616997868383", formula)
// }

// func TestSin(t *testing.T) {
// 	d := NewWithName("var1", 1, 0).Sin()
// 	vars, formula := d.Math()
// 	assert.Equal(t, "sin(var1) = ?", vars)
// 	assert.Equal(t, "sin(1) = 0.841470984807896544828551915928318375739843472469519282898610111931110319333748010828751784005573402229699531838022117989945539661588502120624574802425114599802714611508860519655182175315926637327774878594985045816542706701485174683683726979309922117859910272413672784175028365607893544855897795184024100973080880074046886009375162838756876336134083638363801171409953672944184918309063800980214873465660723218405962257950683415203634506166523593278", formula)
// }

// func TestCos(t *testing.T) {
// 	d := NewWithName("var1", 1, 0).Cos()
// 	vars, formula := d.Math()
// 	assert.Equal(t, "cos(var1) = ?", vars)
// 	assert.Equal(t, "cos(1) = 0.54030230586813965874561515067176071767603141150991567490927772778673118786033739102174242337864109186439207498973007363884202112942385976796862442063752663646870430360736682397798633852405003167527051283327366631405990604840629657123985368031838052877290142895506386796217551784101265975360960112885444847880134909594560331781699767647860744559228420471946006511861233129745921297270844542687374552066388998112901504", formula)
// }

// func TestTan(t *testing.T) {
// 	d := NewWithName("var1", 1, 0).Tan()
// 	vars, formula := d.Math()
// 	assert.Equal(t, "tan(var1) = ?", vars)
// 	assert.Equal(t, "tan(1) = 1.5574077246549025", formula)
// }

func BenchmarkToMath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := NewFromFloatWithName("start", 100)
		for j := 0; j < 100; j++ {
			d = d.Add(NewFromFloatWithName("var"+strconv.Itoa(j), float64(i)))
		}
		d.Math()
	}
}

func BenchmarkDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := decimal.NewFromFloat(100)
		for j := 0; j < 100; j++ {
			d = d.Add(decimal.NewFromFloat(float64(i)))
		}
		d.String()
	}
}
