# toMath

**This library is currently in an alpha state. Do not use in production.**

toMath is a wrapper library for [shopspring/decimal](https://github.com/shopspring/decimal) where you can create a decimal, run operations on it, and output the math underlying those operations. For example:

```go
	d := NewFromFloatWithName("var1", 1.1).
		Round(1).
		Add(NewFromFloatWithName("var2", 1)).
		Add(NewFromFloatWithName("var2", 1)).
		Div(NewFromFloatWithName("var3", 2)).
		Mul(NewFromFloatWithName("var4", 2)).
		SetName("var5")

	vars, formula := d.Math()
	assert.Equal(t, "(round(1)(var1) + var2 + var2) / var3 * var4 = var5", vars)
	assert.Equal(t, "(round(1)(1.1) + 1 + 1) / 2 * 2 = 3.1", formula)

	d = d.Resolve().Add(NewFromFloatWithName("var6", 3)).SetName("var7")
	vars, formula = d.Math()
	assert.Equal(t, "var5 + var6 = var7", vars)
	assert.Equal(t, "3.1 + 3 = 6.1", formula)
```
