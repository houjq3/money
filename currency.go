package money

import (
	"fmt"
	"io"
	"strconv"
)

type Currency string

const (
	CurrencyCny Currency = "CNY"
	CurrencyUsd Currency = "USD"
)

func (e Currency) IsValid() bool {
	switch e {
	case CurrencyCny, CurrencyUsd:
		return true
	}
	return false
}

func (e Currency) String() string {
	return string(e)
}

func (e *Currency) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Currency(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Currency", str)
	}
	return nil
}

func (e Currency) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
