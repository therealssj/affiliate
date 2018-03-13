package spo_client

import (
	"testing"
)

func TestBalanceRespProcess(t *testing.T) {
	res, err := balanceRespProcess([]byte(`{
		"confirmed": {
			"coins": 76000000,
			"hours": 149785
		},
		"predicted": {
			"coins": 76000000,
			"hours": 149785
		}
	}`))
	if err != nil || res != 76000000 {
		t.Errorf("Failed. %d", res)
	}
	res, err = balanceRespProcess([]byte(`{
		"confirmed": {
			"coins": 13000000,
			"hours": 39768
		},
		"predicted": {
			"coins": 13000000,
			"hours": 39768
		}
	}`))
	if err != nil || res != 13000000 {
		t.Errorf("Failed. %d", res)
	}
}
