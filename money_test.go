package money

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

var (
	zero, _    = NewMoney("0", CurrencyCny)
	zeroUSD, _ = NewMoney("0", CurrencyUsd)
	one, _     = NewMoney("1", CurrencyCny)
	oneUSD, _  = NewMoney("1", CurrencyUsd)
	two, _     = NewMoney("2", CurrencyCny)
	three, _   = NewMoney("3", CurrencyCny)
	four, _    = NewMoney("4", CurrencyCny)
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

// func TestM_MarshalAndUnmarshal_Money(t *testing.T) {
// 	money, _ := NewMoney("1.00", CurrencyCny)
// 	buf, err := money.MarshalBSON()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	n, _ := NewMoney("0", CurrencyCny)
// 	err = n.UnmarshalBSON(buf)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if !reflect.DeepEqual(money, n) {
// 		t.Error(money, n, "not equal")
// 	}
// 	if !money.Equal(n) {
// 		t.Error(money, "!=", n)
// 	}
// }

func TestMoney_IsZero(t *testing.T) {
	type fields struct {
		Value    string
		Currency Currency
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"false", fields{"0.11", CurrencyCny}, false},
		{"false", fields{"0.00", CurrencyCny}, true},
		{"true", fields{"0.00", Currency("")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := newMoney4Str(tt.fields.Value, tt.fields.Currency)
			if got := m.IsZeroM(); got != tt.want {
				t.Errorf("Money.IsZero() = %v, want %v", got, tt.want)
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
			}
			got, err := m.Add(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Money.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
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
	if !target.Equal(two) {
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
			}
			got, err := m.Sub(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Money.Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
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
	if !target.Equal(NullMoney) {
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
	if !target.Equal(two) {
		t.Errorf("Money.Sub2() = %v, want %v", target, two)
	}
}

func TestMoney_GetValue(t *testing.T) {
	type fields struct {
		Value    string
		Currency Currency
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"valid money", fields{"1", CurrencyCny}, "1.00"},
		{"valid positive money", fields{"-1", CurrencyCny}, "-1.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := newMoney4Str(tt.fields.Value, tt.fields.Currency)
			if got := m.GetValue(); got != tt.want {
				t.Errorf("Money.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_GetValueFloat64(t *testing.T) {
	a, _ := newMoney4Str("1", CurrencyCny)
	tests := []struct {
		name   string
		fields Money
		want   float64
		want1  bool
	}{
		{"one", a, 1.0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
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
		{"has", Money{Currency: CurrencyUsd.String()}, CurrencyUsd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
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
		{"usd money string", oneUSD, "USD1.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				Value:    tt.fields.Value,
				Currency: tt.fields.Currency,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Money.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_BeginCalc(t *testing.T) {
	calc := one.BeginCalc()
	if calc.Error != nil {
		t.Error(one.String(), calc.String(), calc.Error)
	}
}

func TestMoney_IsPositive(t *testing.T) {
	type fields struct {
		Value    string
		Currency Currency
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"true", fields{"0.11", CurrencyCny}, true},
		{"false", fields{"0.00", CurrencyCny}, false},
		{"false", fields{"-1.00", Currency("")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := newMoney4Str(tt.fields.Value, tt.fields.Currency)
			if got := m.IsPositive(); got != tt.want {
				t.Errorf("Money.IsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_IsNegative(t *testing.T) {
	type fields struct {
		Value    string
		Currency Currency
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"false", fields{"0.11", CurrencyCny}, false},
		{"false", fields{"0.00", CurrencyCny}, false},
		{"true", fields{"-1.00", Currency("")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := newMoney4Str(tt.fields.Value, tt.fields.Currency)
			if got := m.IsNegative(); got != tt.want {
				t.Errorf("Money.IsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}
