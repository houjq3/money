package money

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func DecimalFixd(curr Currency) int32 {
	ret, exists := map[Currency]int32{
		CurrencyCny: 2,
		// TODO: just for test
		CurrencyUsd: 3,
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
