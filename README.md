# toMath

**This library is currently in an alpha state. Do not use in production.**

toMath is a wrapper library for [shopspring/decimal](https://github.com/shopspring/decimal) where you can create a decimal, run operations on it, and output the math underlying those operations. For example:

```go
package main

import "github.com/cbelsole/tomath"

func main() {
	d := tomath.NewFromFloatWithName("var1", 1.1).
		Round(1).
		Add(tomath.NewFromFloatWithName("var2", 1)).
		Add(tomath.NewFromFloatWithName("var2", 1)).
		Div(tomath.NewFromFloatWithName("var3", 2)).
		Mul(tomath.NewFromFloatWithName("var4", 2)).
		SetName("var5")

	vars, formula := d.Math()
	assert.Equal(t, "(round(1)(var1) + var2 + var2) / var3 * var4 = var5", vars)
	assert.Equal(t, "(round(1)(1.1) + 1 + 1) / 2 * 2 = 3.1", formula)

	d = d.Resolve().Add(timesOneHundred(tomath.NewFromFloatWithName("var6", 3))).SetName("var7")
	vars, formula = d.Math()
	assert.Equal(t, "var5 + var6TimesOneHundred = var7", vars)
	assert.Equal(t, "3.1 + 300 = 303.1", formula)
}

func timesOneHundred(input tomath.Decimal) tomath.Decimal {
	oneHundred := tomath.NewFromFloatWithName("OneHundred", 100)
	return input.Mul(oneHundred).ResolveTo(input.GetName() + "Times" + oneHundred.GetName())
}
```
