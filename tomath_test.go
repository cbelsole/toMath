package tomath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	d := New("var1", 0, 0)
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestNewFromInt(t *testing.T) {
	d := NewFromInt("var1", 0)
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestNewFromInt32(t *testing.T) {
	d := NewFromInt32("var1", 0)
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestNewFromBigInt(t *testing.T) {
	d := NewFromBigInt("var1", big.NewInt(0), 0)
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestNewFromString(t *testing.T) {
	d, err := NewFromString("var1", "0")
	assert.NoError(t, err)
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestRequireFromString(t *testing.T) {
	d := RequireFromString("var1", "0")
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestNewFromFloat(t *testing.T) {
	d := NewFromFloat("var1", 0)
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestNewFromFloat32(t *testing.T) {
	d := NewFromFloat32("var1", 0)
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestNewFromFloatWithExponent(t *testing.T) {
	d := NewFromFloatWithExponent("var1", 0, 0)
	vars, formula := d.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `0 = 0`, formula)
}

func TestAbs(t *testing.T) {
	d := New("var1", -1, 0).Abs()
	vars, formula := d.Math()
	assert.Equal(t, `abs(var1) = ?`, vars)
	assert.Equal(t, `abs(-1) = 1`, formula)
}

func TestAdd(t *testing.T) {
	d := New("var1", -1, 0).Add(New("var2", 0, 0))
	vars, formula := d.Math()
	assert.Equal(t, `var1 + var2 = ?`, vars)
	assert.Equal(t, `-1 + 0 = -1`, formula)
}

func TestSub(t *testing.T) {
	d := New("var1", -1, 0).Sub(New("var2", 0, 0))
	vars, formula := d.Math()
	assert.Equal(t, `var1 - var2 = ?`, vars)
	assert.Equal(t, `-1 - 0 = -1`, formula)
}

func TestNeg(t *testing.T) {
	d := New("var1", 1, 0).Neg()
	vars, formula := d.Math()
	assert.Equal(t, `neg(var1) = ?`, vars)
	assert.Equal(t, `neg(1) = -1`, formula)
}

func TestMul(t *testing.T) {
	d := New("var1", 1, 0).Mul(New("var2", 2, 0))
	vars, formula := d.Math()
	assert.Equal(t, `var1 * var2 = ?`, vars)
	assert.Equal(t, `1 * 2 = 2`, formula)
}

func TestShift(t *testing.T) {
	d := New("var1", 1, 0).Shift(1)
	vars, formula := d.Math()
	assert.Equal(t, `shift(1)(var1) = ?`, vars)
	assert.Equal(t, `shift(1)(1) = 10`, formula)
}

func TestDiv(t *testing.T) {
	d := New("var1", 4, 0).Div(New("var2", 2, 0))
	vars, formula := d.Math()
	assert.Equal(t, `var1 / var2 = ?`, vars)
	assert.Equal(t, `4 / 2 = 2`, formula)
}

func TestDivRound(t *testing.T) {
	d := NewFromFloat("var1", 4.333).DivRound(NewFromFloat("var2", 2.7), 3)
	vars, formula := d.Math()
	assert.Equal(t, `divRound(3)(var1 / var2) = ?`, vars)
	assert.Equal(t, `divRound(3)(4.333 / 2.7) = 1.605`, formula)
}

func TestQuoRem(t *testing.T) {
	d1, d2 := NewFromFloat("var1", 4.333).QuoRem(NewFromFloat("var2", 2.7), 3)

	vars, formula := d1.Math()
	assert.Equal(t, `quoRem(3)(var1 / var2) = var1var2Quotient`, vars)
	assert.Equal(t, `quoRem(3)(4.333 / 2.7) = 1.604`, formula)

	vars, formula = d2.Math()
	assert.Equal(t, `quoRem(3)(var1 / var2) = var1var2Remainder`, vars)
	assert.Equal(t, `quoRem(3)(4.333 / 2.7) = 0.0022`, formula)
}

func TestMod(t *testing.T) {
	d := New("var1", 4, 0).Mod(New("var2", 2, 0))
	vars, formula := d.Math()
	assert.Equal(t, `var1 % var2 = ?`, vars)
	assert.Equal(t, `4 % 2 = 0`, formula)
}

func TestPow(t *testing.T) {
	d := New("var1", 4, 0).Pow(New("var2", 2, 0))
	vars, formula := d.Math()
	assert.Equal(t, `var1^var2 = ?`, vars)
	assert.Equal(t, `4^2 = 16`, formula)
}

func TestRound(t *testing.T) {
	d := NewFromFloat("var1", 4.333).Round(2)
	vars, formula := d.Math()
	assert.Equal(t, `round(2)(var1) = ?`, vars)
	assert.Equal(t, `round(2)(4.333) = 4.33`, formula)
}

func TestRoundBank(t *testing.T) {
	d := NewFromFloat("var1", 4.333).RoundBank(2)
	vars, formula := d.Math()
	assert.Equal(t, `roundBank(2)(var1) = ?`, vars)
	assert.Equal(t, `roundBank(2)(4.333) = 4.33`, formula)
}

func TestRoundCash(t *testing.T) {
	d := NewFromFloat("var1", 4.333).RoundCash(5)
	vars, formula := d.Math()
	assert.Equal(t, `roundCash(5)(var1) = ?`, vars)
	assert.Equal(t, `roundCash(5)(4.333) = 4.35`, formula)
}

func TestFloor(t *testing.T) {
	d := NewFromFloat("var1", 4.333).Floor()
	vars, formula := d.Math()
	assert.Equal(t, `floor(var1) = ?`, vars)
	assert.Equal(t, `floor(4.333) = 4`, formula)
}

func TestCeil(t *testing.T) {
	d := NewFromFloat("var1", 4.333).Ceil()
	vars, formula := d.Math()
	assert.Equal(t, `ceil(var1) = ?`, vars)
	assert.Equal(t, `ceil(4.333) = 5`, formula)
}

func TestTruncate(t *testing.T) {
	d := NewFromFloat("var1", 4.333).Truncate(0)
	vars, formula := d.Math()
	assert.Equal(t, `truncate(0)(var1) = ?`, vars)
	assert.Equal(t, `truncate(0)(4.333) = 4`, formula)
}

func TestSetName(t *testing.T) {
	d := NewFromFloat("var1", 1).Add(NewFromFloat("var2", 1)).SetName("var3")
	vars, formula := d.Math()
	assert.Equal(t, `var1 + var2 = var3`, vars)
	assert.Equal(t, `1 + 1 = 2`, formula)
}

func TestComplexExample(t *testing.T) {
	d := NewFromFloat("var1", 1.1).
		Round(1).
		Add(NewFromFloat("var2", 1)).
		Add(NewFromFloat("var2", 1)).
		Div(NewFromFloat("var3", 2)).
		Mul(NewFromFloat("var4", 2)).
		SetName("var5")

	vars, formula := d.Math()
	assert.Equal(t, `((round(1)(var1) + var2 + var2) / var3) * var4 = var5`, vars)
	assert.Equal(t, `((round(1)(1.1) + 1 + 1) / 2) * 2 = 3.1`, formula)

	// the Math() function is idempotent
	vars, formula = d.Math()
	assert.Equal(t, `((round(1)(var1) + var2 + var2) / var3) * var4 = var5`, vars)
	assert.Equal(t, `((round(1)(1.1) + 1 + 1) / 2) * 2 = 3.1`, formula)
}

func TestMin(t *testing.T) {
	d := Min(NewFromFloat("var1", 1), NewFromFloat("var2", 2), NewFromFloat("var3", 100))

	vars, formula := d.Math()
	assert.Equal(t, `min(var1, var2, var3) = ?`, vars)
	assert.Equal(t, `min(1, 2, 100) = 1`, formula)
}

func TestMax(t *testing.T) {
	d := Max(NewFromFloat("var1", 1), NewFromFloat("var2", 2), NewFromFloat("var3", 100))

	vars, formula := d.Math()
	assert.Equal(t, `max(var1, var2, var3) = ?`, vars)
	assert.Equal(t, `max(1, 2, 100) = 100`, formula)
}

func TestSum(t *testing.T) {
	d := Sum(NewFromFloat("var1", 1), NewFromFloat("var2", 2), NewFromFloat("var3", 100))

	vars, formula := d.Math()
	assert.Equal(t, `sum(var1, var2, var3) = ?`, vars)
	assert.Equal(t, `sum(1, 2, 100) = 103`, formula)
}

func TestAvg(t *testing.T) {
	d := Avg(NewFromFloat("var1", 1), NewFromFloat("var2", 2), NewFromFloat("var3", 100))

	vars, formula := d.Math()
	assert.Equal(t, `avg(var1, var2, var3) = ?`, vars)
	assert.Equal(t, `avg(1, 2, 100) = 34.3333333333333333`, formula)
}

func TestRescalePair(t *testing.T) {
	d1, d2 := RescalePair(New("var1", 111111, -5), New("var2", 2111, -3))

	vars, formula := d1.Math()
	assert.Equal(t, `var1 = var1`, vars)
	assert.Equal(t, `1.11111 = 1.11111`, formula)

	vars, formula = d2.Math()
	assert.Equal(t, `var2 = var2`, vars)
	assert.Equal(t, `2.111 = 2.111`, formula)
}

func TestUnmarshalJSON(t *testing.T) {
	d := &Decimal{}
	require.NoError(t, d.UnmarshalJSON([]byte(`123.123`)))
	require.Equal(t, "123.123", d.String())
}

func TestMarshalJSON(t *testing.T) {
	d := NewFromFloat("var1", 123.123)
	b, err := d.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"123.123"`, string(b))
}

func TestBinary(t *testing.T) {
	d := NewFromFloat("var1", 123.123)
	b, err := d.MarshalBinary()
	require.NoError(t, err)

	d2 := &Decimal{}
	require.NoError(t, d2.UnmarshalBinary(b))
	require.Equal(t, "123.123", d.String())
}

func TestScan(t *testing.T) {
	d := Decimal{}
	d2 := NewFromFloat("var1", 54.33)
	require.NoError(t, d.Scan(54.33))
	require.Equal(t, d.String(), d2.String())
}

func TestValue(t *testing.T) {
	d := NewFromFloat("var1", 54.33)
	SetName, err := d.Value()
	require.NoError(t, err)
	require.Equal(t, d.String(), SetName.(string))
}

func TestUnmarshalText(t *testing.T) {
	d := &Decimal{}
	require.NoError(t, d.UnmarshalText([]byte(`123.123`)))
	require.Equal(t, "123.123", d.String())
}

func TestMarshalText(t *testing.T) {
	d := NewFromFloat("var1", 123.123)
	b, err := d.MarshalText()
	require.NoError(t, err)
	require.Equal(t, `123.123`, string(b))
}

func TestGobEncode(t *testing.T) {
	d := NewFromFloat("var1", 123.123)
	b, err := d.GobEncode()
	require.NoError(t, err)

	d2 := &Decimal{}
	require.NoError(t, d2.GobDecode(b))
	require.Equal(t, "123.123", d.String())
}

func TestStringScaled(t *testing.T) {
	d := NewFromFloat("var1", 123.123)
	assert.Equal(t, "123.1", d.StringScaled(-1))

}

func TestUnknownName(t *testing.T) {
	d := Decimal{}
	require.NoError(t, d.UnmarshalJSON([]byte(`123.123`)))
	require.Equal(t, "123.123", d.String())

	vars, formula := d.Math()
	assert.Equal(t, "? = ?", vars)
	assert.Equal(t, "123.123 = 123.123", formula)

	d = d.SetName("var1")
	vars, formula = d.Math()
	assert.Equal(t, "var1 = var1", vars)
	assert.Equal(t, "123.123 = 123.123", formula)
}
