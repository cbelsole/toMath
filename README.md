# toMath

**This library is currently in an alpha state. Do not use in production.**

toMath is a wrapper library for [shopspring/decimal](https://github.com/shopspring/decimal) where you can create a decimal, run operations on it, and output the math underlying those operations. For example:

```
	d := NewFromFloat("var1", 1.1).
		Round(1).
		Add(NewFromFloat("var2", 1)).
		Add(NewFromFloat("var2", 1)).
		Div(NewFromFloat("var3", 2)).
		Mul(NewFromFloat("var4", 2)).
		Result("var5")

	vars, formula := d.Math()
	assert.Equal(t, `((round(1)(var1) + var2 + var2) / var3) * var4 = var5`, vars)
	assert.Equal(t, `((round(1)(1.1) + 1 + 1) / 2) * 2 = 3.1`, formula)
```
