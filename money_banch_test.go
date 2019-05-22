package money

import (
	"testing"
)

func BenchmarkMoneyAdd(b *testing.B) {

	money1, _ := NewMoney("0234123", CurrencyUsd)
	money2, _ := NewMoney("1234234", CurrencyUsd)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		money1, _ = money1.Add(money2)
	}
}
func BenchmarkMoneyCalcAdd(b *testing.B) {

	money1, _ := NewMoney("0234123", CurrencyUsd)
	money2, _ := NewMoney("1234234", CurrencyUsd)
	calc := money1.BeginCalc()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		calc.Add(money2)
	}
}

func BenchmarkMoneyAdd2(b *testing.B) {

	money1, _ := NewMoney("112341232.122330", CurrencyUsd)
	money2, _ := NewMoney("1232341238.908", CurrencyUsd)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		money1.Add2(money2)
	}
}

func BenchmarkMoneyCalcAdd2(b *testing.B) {

	money1, _ := NewMoney("112341232.122330", CurrencyUsd)
	money2, _ := NewMoney("1232341238.908", CurrencyUsd)
	calc := money1.BeginCalc()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		calc.Add(money2)
	}
}

func BenchmarkMoneySub(b *testing.B) {

	money1, _ := NewMoney("1002342.300", CurrencyUsd)
	money2, _ := NewMoney("20392348.1", CurrencyUsd)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		money1, _ = money1.Sub(money2)
	}
}

func BenchmarkMoneyCalcSub(b *testing.B) {

	money1, _ := NewMoney("1002342.300", CurrencyUsd)
	money2, _ := NewMoney("20392348.1", CurrencyUsd)
	calc := money1.BeginCalc()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		calc.Sub(money2)
	}
}
func BenchmarkMoneySub2(b *testing.B) {

	money1, _ := NewMoney("10000.89800", CurrencyUsd)
	money2, _ := NewMoney("123000000000", CurrencyUsd)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		money1.Sub2(money2)
	}
}

func BenchmarkMoneyCalcSub2(b *testing.B) {

	money1, _ := NewMoney("10000.89800", CurrencyUsd)
	money2, _ := NewMoney("123000000000", CurrencyUsd)
	calc := money1.BeginCalc()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		calc.Sub(money2)
	}
}

func BenchmarkMoneyMul(b *testing.B) {

	money1, _ := NewMoney("10000.89800", CurrencyUsd)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		money1.Mul("123456789")
	}
}

func BenchmarkMoneyMul2(b *testing.B) {

	money1, _ := NewMoney("10000.89800", CurrencyUsd)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		money1.Mul2("123456789")
	}
}
func BenchmarkMoneyCalcMul(b *testing.B) {

	money1, _ := NewMoney("10000.89800", CurrencyUsd)
	calc := money1.BeginCalc()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		calc.Mul("123456789")
	}
}

func BenchmarkMoneyCalcMulSimple(b *testing.B) {

	money1, _ := NewMoney("100000", CurrencyUsd)
	calc := money1.BeginCalc()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		calc.Mul("1.23")
	}
}

// goos: darwin
// goarch: amd64
// pkg: git.in.chaitin.net/babysitter/pipeline-crm/server/pkg/money
// BenchmarkMoneyAdd-4                       500000              2498 ns/op             463 B/op         24 allocs/op
// BenchmarkMoneyCalcAdd-4                  1000000              1316 ns/op             272 B/op         13 allocs/op
// BenchmarkMoneyAdd2-4                      500000              4303 ns/op             895 B/op         30 allocs/op
// BenchmarkMoneyCalcAdd2-4                 1000000              1819 ns/op             495 B/op         18 allocs/op
// BenchmarkMoneySub-4                       500000              3932 ns/op             590 B/op         27 allocs/op
// BenchmarkMoneyCalcSub-4                  1000000              1619 ns/op             344 B/op         15 allocs/op
// BenchmarkMoneySub2-4                    100000000               13.8 ns/op             0 B/op          0 allocs/op
// BenchmarkMoneyCalcSub2-4                 1000000              2362 ns/op            5584 B/op         12 allocs/op
// BenchmarkMoneyMul-4                      1000000              2174 ns/op             408 B/op         18 allocs/op
// BenchmarkMoneyMul2-4                     1000000              1100 ns/op             248 B/op         11 allocs/op
// BenchmarkMoneyCalcMul-4                   200000            105510 ns/op          340095 B/op          6 allocs/op
// BenchmarkMoneyCalcMulSimple-4             500000             67576 ns/op          221026 B/op          8 allocs/op
