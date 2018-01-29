package db

import (
	"sort"
	"testing"
)

func TestCryptocurrencyInfoSliceSort(t *testing.T) {
	slice := make([]CryptocurrencyInfo, 0, 4)
	slice = append(slice, CryptocurrencyInfo{"BTC", "BTC", "10000", 8})
	slice = append(slice, CryptocurrencyInfo{"ETH", "ETH", "5000", 9})
	slice = append(slice, CryptocurrencyInfo{"SKY", "SKY", "127", 6})
	sort.Sort(CryptocurrencyInfoSlice(slice))
	if slice[0].ShortName != "BTC" {
		t.Errorf("Sort Failed.")
	}
	if slice[1].ShortName != "ETH" {
		t.Errorf("Sort Failed.")
	}
	if slice[2].ShortName != "SKY" {
		t.Errorf("Sort Failed.")
	}
}
