package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/go-playground/locales/currency"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/zh_Hans_CN"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/shopspring/decimal"
)

// FIXME: Implement marshal and unmarshal for json and bson like decimal.Decimal, Can avoid erroneous use struct

var (
	locale_zhHans = zh_Hans_CN.New()
	locale_enUS   = en_US.New()

	ErrCurrencyUnmatched = errors.New("wrong currency type")

	NullMoney Money = Money{"0", CurrencyCny, decimal.New(0, 0)}
)

// Money Do not set or get values by point operation
type Money struct {
	Value    string   `bson:"value"`
	Currency Currency `bson:"currency"`
	money    decimal.Decimal
}

// NewMoney return new money
// Don't modify internal variables in this form: **m.Value = "100"**
// TODO: NewMoney("-1e2000000000",    CurrencyCny) ? sloooooooooooow ~
func NewMoney(value string, currency Currency) (Money, error) {
	v, err := decimal.NewFromString(value)
	if err != nil {
		return NullMoney, err
	}
	return newMoney(v, currency)
}

// newMoney is internal func
// TODO: obj pool?
func newMoney(dec decimal.Decimal, cur Currency) (Money, error) {
	if !cur.IsValid() {
		cur = CurrencyCny
	}
	return Money{dec.StringFixed(DecimalFixd(cur)), cur, dec}, nil
}
func (m *Money) maintain() {
	m.Value = m.GetValue()
}
func (m Money) IsValid() (ret bool) {
	if m.money.StringFixed(DecimalFixd(m.GetCurrency())) == m.Value && m.Currency.IsValid() {
		ret = true
	}
	return
}

// IsZero
//	true if d == 0
//	false if d > 0
//	false if d < 0
func (m Money) IsZero() bool {
	return m.money.IsZero()
}

// Equal m == b
func (m Money) Equal(b Money) bool {
	return m.money.Equal(b.money)
}

// Copy deep copy
func (m Money) Copy() (Money, error) {
	return newMoney(CopyDecimal(m.money), m.GetCurrency())
}

// Add return m + t
func (m Money) Add(t Money) (Money, error) {
	if m.GetCurrency() != t.GetCurrency() {
		return NullMoney, ErrCurrencyUnmatched
	}
	ret := CopyDecimal(m.money)
	return newMoney(ret.Add(t.money), m.GetCurrency())
}

// Add2 m += t
func (m *Money) Add2(t Money) error {
	if m.GetCurrency() != t.GetCurrency() {
		return ErrCurrencyUnmatched
	}
	m.money = m.money.Add(t.money)
	m.maintain()
	return nil
}

// Sub return m - t
func (m Money) Sub(t Money) (Money, error) {
	if m.GetCurrency() != t.GetCurrency() {
		return NullMoney, ErrCurrencyUnmatched
	}
	ret := CopyDecimal(m.money)
	return newMoney(ret.Sub(t.money), m.GetCurrency())
}

// Sub2 m -= t
func (m *Money) Sub2(t Money) error {
	if m.GetCurrency() != t.GetCurrency() {
		return ErrCurrencyUnmatched
	}
	m.money = m.money.Sub(t.money)
	m.maintain()
	return nil
}

// Mul return m * b
func (m Money) Mul(b string) (ret Money, err error) {
	d, err := decimal.NewFromString(b)
	if err != nil {
		return
	}
	retd := CopyDecimal(m.money)
	return newMoney(retd.Mul(d), m.GetCurrency())
}

// Mul2 m *= b
func (m *Money) Mul2(b string) (err error) {
	d, err := decimal.NewFromString(b)
	if err != nil {
		return
	}
	m.money = m.money.Mul(d)
	m.maintain()
	return
}

// GetValue returns a rounded fixed-point string
func (m Money) GetValue() string {
	return m.money.StringFixed(DecimalFixd(m.GetCurrency()))
}

// GetValueFloat64 returns the nearest float64 value for m and a bool indicating whether f represents m exactly.
// For more details, see the documentation for big.Rat.Float64
func (m *Money) GetValueFloat64() (float64, bool) {
	return m.money.Float64()
}

// GetCurrency get currency, if currency is invalid, return CNY
func (m Money) GetCurrency() Currency {
	ret := CurrencyCny
	if m.Currency.IsValid() {
		ret = m.Currency
	}
	return ret
}

func (m Money) String() string {
	val, exact := m.GetValueFloat64()
	if !exact {
		// FIXME: how to deal with this?
		log.Println("unexactly?")
	}
	switch m.GetCurrency() {
	case CurrencyCny:
		return locale_zhHans.FmtCurrency(val, uint64(DecimalFixd(m.GetCurrency())), currency.CNY)
	case CurrencyUsd:
		return locale_enUS.FmtCurrency(val, uint64(DecimalFixd(m.GetCurrency())), currency.USD)
	default:
		return fmt.Sprintf("%v %v", m.GetValue(), m.GetCurrency())
	}
}
func (m *Money) UnmarshalBSON(payload []byte) error {
	var out bson.Document
	err := bson.Unmarshal(payload, &out)
	if err != nil {
		return err
	}
	i := out.Iterator()
	for i.Next() {
		elem := i.Element()
		switch elem.Key() {
		case "value":
			m.Value, _ = elem.Value().StringValueOK()
		case "currency":
			cur, ok := elem.Value().StringValueOK()
			if ok {
				m.Currency = Currency(cur)
			}
		}
	}
	v, err := decimal.NewFromString(m.Value)
	if err != nil {
		return err
	}
	m.money = v
	return nil
}
func (m Money) MarshalBSON() (ret []byte, err error) {
	d := bson.NewDocument()
	d.Append(bson.EC.String("value", m.Value))
	d.Append(bson.EC.String("currency", m.Currency.String()))
	return d.MarshalBSON()
}

func (m *Money) UnmarshalGQL(v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, m)
	if err != nil {
		return err
	}
	dec, err := decimal.NewFromString(m.Value)
	if err != nil {
		return err
	}
	m.money = dec
	return err
}

func (m Money) MarshalGQL(w io.Writer) {
	data := struct {
		Value    string   `json:"value"`
		Currency Currency `json:"currency"`
	}{m.GetValue(), m.GetCurrency()}

	buf, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Write(buf)
}
