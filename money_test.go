package money

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

var (
	one, _    = NewMoney("1", CurrencyCny)
	oneUSD, _ = NewMoney("1", CurrencyUsd)
	two, _    = NewMoney("2", CurrencyCny)
	three, _  = NewMoney("3", CurrencyCny)
	four, _   = NewMoney("4", CurrencyCny)
)

func TestNewMoney(t *testing.T) {
	t.Log(NewMoney("0.11", CurrencyCny))
	t.Log(NewMoney("+0.11", CurrencyCny))
	t.Log(NewMoney("-0.11", CurrencyCny))
	t.Log(NewMoney("1e200", CurrencyCny))
	t.Log(NewMoney("1e-1", CurrencyCny))
}

func Test_Copy(t *testing.T) {
	if !CopyDecimal(decimal.New(1, 1)).Equal(decimal.New(1, 1)) {
		t.Error("copy decimal faild")
	}
}

func TestM_MarshalAndUnmarshal_Money(t *testing.T) {
	money, _ := NewMoney("1.00", CurrencyCny)
	buf, err := money.MarshalBSON()
	if err != nil {
		t.Error(err)
	}
	n, _ := NewMoney("0", CurrencyCny)
	err = n.UnmarshalBSON(buf)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(money, n) {
		t.Error(money, n, "not equal")
	}
	if !money.Equal(n) {
		t.Error(money, "!=", n)
	}
}

// TODO: reflect.DeepEqual?
func Test_newMoney(t *testing.T) {
	d := decimal.New(1, 1)
	money1, _ := newMoney(d, CurrencyUsd)
	t.Log(money1)
	want := Money{"1.000", CurrencyUsd, d}
	t.Log(want)
	if want.GetValue() != money1.GetValue() {
		t.Error(money1, want)
	}

	money2, _ := newMoney(d, CurrencyCny)
	t.Log(money2)
	if want.GetValue() == money2.GetValue() {
		t.Error(money2, want)
	}
}

func TestMoney_maintain(t *testing.T) {
	type fields struct {
		Value    string
		Currency Currency
		money    decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"invalid money to valid", fields{"1.0", CurrencyCny, decimal.New(1, 1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			m.maintain()
			t.Log(m.Value, m.Currency, m.money)
			if m.GetValue() != m.Value {
				t.Error("maintain error", m)
			}
		})
	}
}

func TestMoney_IsValid(t *testing.T) {
	one := decimal.New(1, 1)
	money1 := Money{"1", CurrencyCny, one}
	if money1.IsValid() {
		t.Error(money1)
	}
	money1.maintain()
	if !money1.IsValid() {
		t.Error(money1)
	}
}

func TestMoney_IsZero(t *testing.T) {
	type fields struct {
		Value    string
		Currency Currency
		money    decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"true", fields{"0.11", CurrencyCny, decimal.New(11, -1)}, false},
		{"false", fields{"0.00", CurrencyCny, decimal.New(0, 0)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			if got := m.IsZero(); got != tt.want {
				t.Errorf("Money.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Copy(t *testing.T) {
	type fields struct {
		Value    string
		Currency Currency
		money    decimal.Decimal
	}
	tests := []struct {
		name    string
		fields  fields
		want    Money
		wantErr bool
	}{
		{"valid money", fields{"0.11", CurrencyCny, decimal.New(11, -2)}, Money{"0.11", CurrencyCny, decimal.New(11, -2)}, false},
		{"invalid money", fields{"0.1", CurrencyCny, decimal.New(11, -2)}, Money{"0.11", CurrencyCny, decimal.New(11, -2)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			got, err := m.Copy()
			if (err != nil) != tt.wantErr {
				t.Errorf("Money.Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Money.Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	type args struct {
		t Money
	}
	tests := []struct {
		name    string
		fields  Money
		args    args
		want    Money
		wantErr bool
	}{
		{"1+1", one, args{one}, two, false},
		{"3+1", three, args{one}, four, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			got, err := m.Add(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Money.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.money.Equal(tt.want.money) {
				t.Errorf("Money.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Add2(t *testing.T) {
	target, err := NewMoney("1", CurrencyCny)
	if err != nil {
		t.Error(err)
	}
	target.Add2(one)
	if !target.money.Equal(two.money) {
		t.Errorf("Money.Add2() = %v, want %v", target, two)
	}
}

func TestMoney_Sub(t *testing.T) {
	type args struct {
		t Money
	}
	tests := []struct {
		name    string
		fields  Money
		args    args
		want    Money
		wantErr bool
	}{
		{"1-1", one, args{one}, NullMoney, false},
		{"3-1", three, args{one}, two, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			got, err := m.Sub(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Money.Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.money.Equal(tt.want.money) {
				t.Errorf("Money.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Sub2(t *testing.T) {
	target, err := NewMoney("1", CurrencyCny)
	if err != nil {
		t.Error(err)
	}
	target.Sub2(one)
	if !target.money.Equal(NullMoney.money) {
		t.Errorf("Money.Sub2() = %v, want %v", target, two)
	}
}

func TestMoney_Mul(t *testing.T) {
	type args struct {
		b string
	}
	tests := []struct {
		name    string
		fields  Money
		args    args
		wantRet Money
		wantErr bool
	}{
		{"1*1", one, args{"1"}, one, false},
		{"2*2", two, args{"2"}, four, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			gotRet, err := m.Mul(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Money.Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("Money.Mul() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestMoney_Mul2(t *testing.T) {
	target, err := NewMoney("1", CurrencyCny)
	if err != nil {
		t.Error(err)
	}
	target.Mul2("2")
	if !target.money.Equal(two.money) {
		t.Errorf("Money.Sub2() = %v, want %v", target, two)
	}
}

func TestMoney_GetValue(t *testing.T) {
	type fields struct {
		Value    string
		Currency Currency
		money    decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"valid money", fields{"1", CurrencyCny, one.money}, "1.00"},
		{"valid positive money", fields{"-1", CurrencyCny, one.money.Mul(decimal.New(-1, 1))}, "-10.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			if got := m.GetValue(); got != tt.want {
				t.Errorf("Money.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_GetValueFloat64(t *testing.T) {
	tests := []struct {
		name   string
		fields Money
		want   float64
		want1  bool
	}{
		{"one", Money{"1", CurrencyCny, one.money}, 1.0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			got, got1 := m.GetValueFloat64()
			if got != tt.want {
				t.Errorf("Money.GetValueFloat64() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Money.GetValueFloat64() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMoney_GetCurrency(t *testing.T) {
	tests := []struct {
		name   string
		fields Money
		want   Currency
	}{
		{"nil", Money{}, CurrencyCny},
		{"has", Money{Currency: CurrencyUsd}, CurrencyUsd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			if got := m.GetCurrency(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Money.GetCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_String(t *testing.T) {
	tests := []struct {
		name   string
		fields Money
		want   string
	}{
		{"cny money string", one, "CNY1.00"},
		{"usd money string", oneUSD, "USD1.000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
				money:    tt.fields.money,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Money.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
