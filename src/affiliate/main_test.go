package main

import (

	//	"github.com/shopspring/decimal"

	"testing"

	"github.com/spolabs/affiliate/src/service/db"
)

func TestConvertUnitPower(t *testing.T) {
	//	fmt.Println(decimal.New(int64(100000), 0).String())
	//	fmt.Println(decimal.New(1, 6).String())
	if "0.1" != convertUnitPower(100000, 6) {
		t.Errorf("Failed. Got %s, expected %s.", convertUnitPower(100000, 6), "0.1")
	}
}

func TestNewCryptocurrency(t *testing.T) {
	info := db.CryptocurrencyInfo{"BTC", "BTC", "50000", 8, true}
	cc := newCryptocurrency(&info)
	if "0.00002" != cc.ReverseRate {
		t.Errorf("Failed. %s", cc.ReverseRate)
	}
	info = db.CryptocurrencyInfo{"ETH", "ETH", "4019", 9, true}
	cc = newCryptocurrency(&info)
	if "0.000248818" != cc.ReverseRate {
		t.Errorf("Failed. %s", cc.ReverseRate)
	}
	info = db.CryptocurrencyInfo{"XMR", "XMR", "1492", 12, true}
	cc = newCryptocurrency(&info)
	if "0.000670241287" != cc.ReverseRate {
		t.Errorf("Failed. %s", cc.ReverseRate)
	}
	info = db.CryptocurrencyInfo{"SKY", "SKY", "81", 6, true}
	cc = newCryptocurrency(&info)
	if "0.012346" != cc.ReverseRate {
		t.Errorf("Failed. %s", cc.ReverseRate)
	}
	// json, _ := json.Marshal(cc)
	// t.Errorf(string(json))
}
