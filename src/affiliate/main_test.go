package main

import (
	//	"fmt"
	//	"github.com/shopspring/decimal"
	"testing"
)

func TestConvertUnitPower(t *testing.T) {
	//	fmt.Println(decimal.New(int64(100000), 0).String())
	//	fmt.Println(decimal.New(1, 6).String())
	if "0.1" != convertUnitPower(100000, 6) {
		t.Errorf("Failed. Got %s, expected %s.", convertUnitPower(100000, 6), "0.1")
	}
}
