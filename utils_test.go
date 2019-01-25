package money

import (
	"testing"
)

func TestDecimalFixd(t *testing.T) {
	if DecimalFixd(CurrencyUsd) != 3 {
		t.Error("error")
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want string
	}{
		{"", 1, "1"},
		{"", 1.1, "1.1"},
		{"", -0, "0"},
		{"", -1, "-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
