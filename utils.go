package money

import (
	"fmt"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func DecimalFixd(curr Currency) int32 {
	ret, exists := map[Currency]int32{
		CurrencyCny: 2,
		CurrencyUsd: 2,
	}[curr]
	if exists {
		return ret
	}
	return 2
}
func CopyDecimal(d decimal.Decimal) decimal.Decimal {
	// FIXME: ugly copy
	ret := decimal.New(0, 0)
	buf, _ := d.MarshalBinary()
	ret.UnmarshalBinary(buf)
	return ret
}

func ToString(input interface{}) string {
	return fmt.Sprintf("%v", input)
}

func TwoMoneyToDecimal(first Money, second Money) (decimal.Decimal, decimal.Decimal, error) {
	dd, err := first.ToDecimal()
	if err != nil {
		return DecimalZero, DecimalZero, err
	}
	bd, err := second.ToDecimal()
	if err != nil {
		return DecimalZero, DecimalZero, err
	}
	return dd, bd, nil
}

func AplusBeqC(a, b, c string) (ret bool) {
	ad, err := decimal.NewFromString(a)
	if err != nil {
		return
	}
	bd, err := decimal.NewFromString(b)
	if err != nil {
		return
	}
	cd, err := decimal.NewFromString(c)
	if err != nil {
		return
	}
	if ad.Add(bd).Equal(cd) {
		ret = true
	}
	return
}

func Add2Str(a, b string) string {
	ad, _ := decimal.NewFromString(a)
	bd, _ := decimal.NewFromString(b)
	return ad.Add(bd).String()
}

func EqualDecimal(a, b string) bool {
	ad, _ := decimal.NewFromString(a)
	bd, _ := decimal.NewFromString(b)
	return ad.Equal(bd)
}

func IsZero(a interface{}) bool {
	switch tt := a.(type) {
	case string:
		ad, _ := decimal.NewFromString(tt)
		return ad.IsZero()
	}
	return false
}

func NormalizeDecimal(val string) string {
	ret, err := decimal.NewFromString(val)
	if err != nil {
		return "0"
	}
	return ret.String()
}

func Div2Money(a, b Money) float64 {
	if b.IsZeroM() {
		log.Warn("b is 0")
		return 0
	}
	aF, _ := a.GetValueFloat64()
	bF, _ := b.GetValueFloat64()
	return aF / bF
}
