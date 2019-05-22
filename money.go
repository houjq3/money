package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/go-playground/locales/currency"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/zh_Hans_CN"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	bsond "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	locale_zhHans = zh_Hans_CN.New()
	locale_enUS   = en_US.New()

	BsonDecimalZero bsond.Decimal128
	DecimalZero     decimal.Decimal

	ErrCurrencyUnmatched       = errors.New("wrong currency type")
	NullMoney            Money = Money{BsonDecimalZero, ""}
)

func init() {
	BsonDecimalZero, _ = bsond.ParseDecimal128("0")
	DecimalZero, _ = decimal.NewFromString("0")
}

type Money struct {
	Value    bsond.Decimal128 `bson:"value"`
	Currency string           `bson:"currency"`
}

const (
	MinMoney float64 = -1e10
	MaxMoney float64 = 1e10
)

// NewMoney return new money
func NewMoney(value string, currency Currency) (Money, error) {
	// too large value will slow!
	valueF, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return NullMoney, err
	}
	if valueF > MaxMoney || valueF < MinMoney {
		return NullMoney, fmt.Errorf("invalid money range: %s, should in %v %v", value, MinMoney, MaxMoney)
	}
	v, err := bsond.ParseDecimal128(value)
	if err != nil {
		return NullMoney, err
	}
	return newMoney(v, currency)
}

func (m Money) BeginCalc() *MoneyCalculation {
	ret := &MoneyCalculation{}
	ret.money, ret.Error = m.ToDecimal()
	ret.currency = m.GetCurrency()
	return ret
}

// newMoney is internal func
func newMoney(dec bsond.Decimal128, cur Currency) (Money, error) {
	if !cur.IsValid() {
		cur = CurrencyCny
	}
	return Money{dec, cur.String()}, nil
}
func newMoney4Str(val string, cur Currency) (Money, error) {
	if !cur.IsValid() {
		cur = CurrencyCny
	}
	v, err := bsond.ParseDecimal128(val)
	if err != nil {
		return NullMoney, err
	}
	return newMoney(v, cur)
}
func (m Money) IsValid() (ret bool) {
	if Currency(m.Currency).IsValid() {
		ret = true
	}
	return
}

// IsZero
//	true if d == 0
//	false if d > 0
//	false if d < 0
func (m Money) IsZeroM() bool {
	if m.Currency == "" {
		return true
	}
	d, err := decimal.NewFromString(m.Value.String())
	if err != nil {
		return true
	}
	return d.IsZero()
}

// IsNegative return
//
//	true if d < 0
//	false if d == 0
//	false if d > 0
func (m Money) IsNegative() bool {
	if m.Currency == "" {
		return true
	}
	d, err := decimal.NewFromString(m.Value.String())
	if err != nil {
		return true
	}
	return d.IsNegative()
}

// IsPositive return
//
//	true if d > 0
//	false if d == 0
//	false if d < 0
func (m Money) IsPositive() bool {
	if m.Currency == "" {
		return true
	}
	d, err := decimal.NewFromString(m.Value.String())
	if err != nil {
		return true
	}
	return d.IsPositive()
}

func (m Money) ToDecimal() (decimal.Decimal, error) {
	return decimal.NewFromString(m.Value.String())
}

// Equal m == b
func (m Money) Equal(b Money) bool {
	md, bd, err := TwoMoneyToDecimal(m, b)
	if err != nil {
		return false
	}
	return md.Equal(bd) && m.GetCurrency() == b.GetCurrency()
}

func (m Money) LessThan(b Money) bool {
	md, bd, err := TwoMoneyToDecimal(m, b)
	if err != nil {
		return false
	}
	return md.LessThan(bd)
}

func (m Money) LessThanOrEqual(b Money) bool {
	md, bd, err := TwoMoneyToDecimal(m, b)
	if err != nil {
		return false
	}
	return md.LessThanOrEqual(bd)
}

// Copy deep copy
func (m Money) Copy() Money {
	return m
}

// Add return m + t
func (m Money) Add(t Money) (Money, error) {
	md, td, err := TwoMoneyToDecimal(m, t)
	if err != nil {
		return NullMoney, err
	}
	return newMoney4Str(md.Add(td).String(), m.GetCurrency())
}

// Add2 m += t
func (m *Money) Add2(t Money) error {
	md, td, err := TwoMoneyToDecimal(*m, t)
	if err != nil {
		return err
	}
	*m, err = newMoney4Str(md.Add(td).String(), m.GetCurrency())
	return err
}

// Sub return m - t
func (m Money) Sub(t Money) (Money, error) {
	md, td, err := TwoMoneyToDecimal(m, t)
	if err != nil {
		return NullMoney, err
	}
	return newMoney4Str(md.Sub(td).String(), m.GetCurrency())
}

// Sub2 m -= t
func (m *Money) Sub2(t Money) error {
	md, td, err := TwoMoneyToDecimal(*m, t)
	if err != nil {
		return err
	}
	*m, err = newMoney4Str(md.Sub(td).String(), m.GetCurrency())
	return err
}

// Mul return m * b
func (m Money) Mul(b string) (Money, error) {
	d, err := decimal.NewFromString(b)
	if err != nil {
		return NullMoney, err
	}
	md, err := m.ToDecimal()
	return newMoney4Str(md.Mul(d).String(), m.GetCurrency())
}

// Mul2 m *= b
func (m *Money) Mul2(b string) (err error) {
	d, err := decimal.NewFromString(b)
	if err != nil {
		return
	}
	md, err := m.ToDecimal()
	if err != nil {
		return err
	}
	*m, err = newMoney4Str(md.Mul(d).String(), m.GetCurrency())
	return
}

// GetValue returns a rounded fixed-point string
func (m Money) GetValue() string {
	md, err := m.ToDecimal()
	if err != nil {
		log.Warn(err)
	}
	return md.StringFixed(DecimalFixd(m.GetCurrency()))
}

// GetValueFloat64 returns the nearest float64 value for m and a bool indicating whether f represents m exactly.
// For more details, see the documentation for big.Rat.Float64
func (m Money) GetValueFloat64() (float64, bool) {
	md, err := m.ToDecimal()
	if err != nil {
		log.Warn(err)
	}
	return md.Float64()
}

// GetCurrency get currency, if currency is invalid, return CNY
func (m Money) GetCurrency() Currency {
	ret := CurrencyCny
	if Currency(m.Currency).IsValid() {
		ret = Currency(m.Currency)
	}
	return ret
}

// StringShort for export excel
func (m Money) StringShort() string {
	switch m.GetCurrency() {
	case CurrencyCny:
		return "¥ " + m.GetValue()
	case CurrencyUsd:
		return "$ " + m.GetValue()
	case CurrencyOthers:
		return m.GetValue()
	default:
		return m.String()
	}
}

func (m Money) String() string {
	val, exact := m.GetValueFloat64()
	if !exact {
		// FIXME: how to deal with this?
		log.Warn("unexactly?")
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

type innerMoney struct {
	Value    string   `json:"value"`
	Currency Currency `json:"currency"`
}

func (m *Money) UnmarshalGQL(v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data := innerMoney{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return err
	}
	dec, err := bsond.ParseDecimal128(data.Value)
	if err != nil {
		return err
	}
	m.Value = dec
	m.Currency = data.Currency.String()
	return err
}

func (m Money) MarshalGQL(w io.Writer) {
	data := innerMoney{m.GetValue(), m.GetCurrency()}
	buf, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Write(buf)
}
