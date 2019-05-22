package money

import (
	"testing"
)

func TestMoneyCalculation_Add(t *testing.T) {
	calc := one.BeginCalc()
	if calc.Error != nil {
		t.Error(one.String(), calc.String(), calc.Error)
	}
	if calc.Add(one).Error != nil {
		t.Error("1+1", calc.String())
	}
	if !calc.ToMoney().Equal(two) {
		t.Error("1+1!=2?", calc.String())
	}
	if calc.Add(one).Error != nil {
		t.Error("2+1=", calc.String())
	}
	if !calc.ToMoney().Equal(three) {
		t.Error("2+1!=3?", calc.String())
	}
}

func TestMoneyCalculation_Sub(t *testing.T) {
	calc := one.BeginCalc()
	if calc.Error != nil {
		t.Error(one.String(), calc.String(), calc.Error)
	}
	if calc.Sub(one).Error != nil {
		t.Error("1-1=", calc.String())
	}
	if !calc.ToMoney().Equal(zero) {
		t.Error("1-1!=0?", calc.String())
	}
}

func TestMoneyCalculation_Mul(t *testing.T) {
	calc := one.BeginCalc()
	if calc.Error != nil {
		t.Error(one.String(), calc.String(), calc.Error)
	}
	if calc.Mul("3").Error != nil {
		t.Error("1*3", calc.String())
	}
	if !calc.ToMoney().Equal(three) {
		t.Error("1*3!=3?", calc.String())
	}
}
