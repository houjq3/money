package money

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type MoneyCalculation struct {
	currency Currency
	money    decimal.Decimal
	Error    error
}

func (m *MoneyCalculation) ToMoney() (ret Money) {
	ret, m.Error = NewMoney(m.money.StringFixed(DecimalFixd(m.currency)), m.currency)
	return ret
}

func (m *MoneyCalculation) Add(t Money) *MoneyCalculation {
	if m.Error != nil {
		return m
	}
	var tDec decimal.Decimal
	tDec, m.Error = t.ToDecimal()
	m.money = m.money.Add(tDec)
	return m
}

func (m *MoneyCalculation) Sub(t Money) *MoneyCalculation {
	if m.Error != nil {
		return m
	}
	var tDec decimal.Decimal
	tDec, m.Error = t.ToDecimal()
	m.money = m.money.Sub(tDec)
	return m
}

func (m *MoneyCalculation) Mul(b string) *MoneyCalculation {
	if m.Error != nil {
		return m
	}
	var bDec decimal.Decimal
	bDec, m.Error = decimal.NewFromString(b)
	m.money = m.money.Mul(bDec)
	return m
}

func (m *MoneyCalculation) Div(b string) *MoneyCalculation {
	if m.Error != nil {
		return m
	}
	var bDec decimal.Decimal
	bDec, m.Error = decimal.NewFromString(b)
	m.money = m.money.Div(bDec)
	return m
}
func (m *MoneyCalculation) String() string {
	return fmt.Sprintf("%v%v Err:%v", m.currency.String(), m.money.String(), m.Error)
}
