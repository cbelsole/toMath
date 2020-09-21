# toMath

toMath is a wrapper library for [shopspring/decimal](https://github.com/shopspring/decimal), an arbitrary-precision fixed-point decimal numbers in go, where you can create a Decimal, run operations on it, and output the math underlying those operations.

## Install
```
go get github.com/cbelsole/tomath
```

## Requirements
Decimal library requires Go version `>=1.7`

## Usage
```go
package main

import (
	"github.com/cbelsole/tomath"
	"github.com/stretchr/testify/assert"
)

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

### Notes

* toMath makes no assertions on Decimal names. Use `()+-*/^%=` characters at your own risk.

## Documentation
[pkg.go.dev/github.com/cbelsole/tomath](https://pkg.go.dev/github.com/cbelsole/tomath)

## License
The MIT License (MIT)

[shopspring/decimal](https://github.com/shopspring/decimal) - The MIT License (MIT)

[fpd.Decimal](https://github.com/oguzbilgic/fpd) - The MIT License (MIT)
