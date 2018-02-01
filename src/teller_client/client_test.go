package teller_client

import (
	//	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/spaco/affiliate/src/config"
	"net/http/httptest"
	"testing"
)

func TestSetAuthHeaders(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	token := "test-token"
	teller := &config.Teller{"", token, false}
	setAuthHeaders(r, teller)
	timestamp := r.Header.Get("timestamp")
	auth := r.Header.Get("auth")
	hash := hmac.New(sha256.New, []byte(token))
	hash.Write([]byte(timestamp))
	if hex.EncodeToString(hash.Sum(nil)) != auth {
		t.Errorf("Failed. check error")
	}
}

func TestStatusRespProcess1(t *testing.T) {
	json := `{
    "errmsg": "",
    "code": 0,
    "data": {
        "statuses": [
            {
                "seq": 0,
                "updated_at": 1517369203,
                "status": "done",
                "coin_type": "ETH"
            },
            {
                "seq": 1,
                "updated_at": 1517369203,
                "status": "done",
                "coin_type": "SKY"
            },
            {
                "seq": 2,
                "updated_at": 1517380670,
                "status": "done",
                "coin_type": "BTC"
            },
            {
                "seq": 3,
                "updated_at": 1517382139,
                "status": "done",
                "coin_type": "SKY"
            }
        ]
    }
}`
	arr, err := statusRespProcess([]byte(json), "SKY")
	if err != nil {
		t.Errorf("Failed.")
	}
	if len(arr) != 2 {
		t.Errorf("Failed.")
	}
	for _, s := range arr {
		if s.Status != StatusDone {
			t.Errorf("Failed.")
		}
	}
	arr, err = statusRespProcess([]byte(json), "ETH")
	if err != nil {
		t.Errorf("Failed.")
	}
	if len(arr) != 1 {
		t.Errorf("Failed.")
	}
	for _, s := range arr {
		if s.Status != StatusDone {
			t.Errorf("Failed.")
		}
	}
	arr, err = statusRespProcess([]byte(json), "BTC")
	if err != nil {
		t.Errorf("Failed.")
	}
	if len(arr) != 1 {
		t.Errorf("Failed.")
	}
	for _, s := range arr {
		if s.Status != StatusDone {
			t.Errorf("Failed.")
		}
	}
}

func TestStatusRespProcess2(t *testing.T) {
	json := `{
    "errmsg": "",
    "code": 0,
    "data": {
        "statuses": [
            {
                "seq": 0,
                "updated_at": 1517305932,
                "status": "done",
                "coin_type": "SKY"
            },
            {
                "seq": 1,
                "updated_at": 1517387087,
                "status": "waiting_deposit",
                "coin_type": "SKY"
            }
        ]
    }
}`
	arr, err := statusRespProcess([]byte(json), "SKY")
	if err != nil {
		t.Errorf("Failed.")
	}
	if len(arr) != 2 {
		t.Errorf("Failed.")
	}
	for _, s := range arr {
		if s.Seq == 0 && s.Status != StatusDone {
			t.Errorf("Failed.")
		}
		if s.Seq == 1 && s.Status != StatusWaitingDeposit {
			t.Errorf("Failed.")
		}
	}
}

func TestConfigRespProcess(t *testing.T) {
	json := `{
    "errmsg": "",
    "code": 0,
    "data": {
        "enabled": true,
        "max_bound_addrs": 4,
        "max_decimals": 3,
        "allcoins": {
            "BTC": {
                "coin_name": "BTC",
                "rate": "55556",
                "enabled": true,
                "unit": 100000000,
                "confirmations_required": 1
            },
            "ETH": {
                "coin_name": "ETH",
                "rate": "5912",
                "enabled": true,
                "unit": 1000000000,
                "confirmations_required": 1
            },
            "SKY": {
                "coin_name": "SKY",
                "rate": "126",
                "enabled": true,
                "unit": 1000000,
                "confirmations_required": 0
            },
            "XMR": {
                "coin_name": "XMR",
                "rate": "1601",
                "enabled": true,
                "unit": 1000000000000,
                "confirmations_required": 1
            }
        }
    }
}`
	confResp, err := configRespProcess([]byte(json))
	if err != nil {
		t.Errorf("Failed.")
	}
	if len(confResp.AllCoins) != 4 {
		t.Errorf("Failed.")
	}
	for _, coin := range confResp.AllCoins {
		if !coin.Enabled {
			t.Errorf("Failed.")
		}
	}
}

func TestRateWithErrProcess(t *testing.T) {
	json := `{
    "errmsg": "",
    "code": 0,
    "data": {
        "enabled": true,
        "max_bound_addrs": 4,
        "max_decimals": 3,
        "allcoins": {
            "BTC": {
                "coin_name": "BTC",
                "rate": "55556",
                "enabled": true,
                "unit": 100000000,
                "confirmations_required": 1
            },
            "ETH": {
                "coin_name": "ETH",
                "rate": "5912",
                "enabled": true,
                "unit": 1000000000,
                "confirmations_required": 1
            },
            "SKY": {
                "coin_name": "SKY",
                "rate": "126",
                "enabled": true,
                "unit": 1000000,
                "confirmations_required": 0
            },
            "XMR": {
                "coin_name": "XMR",
                "rate": "1601",
                "enabled": true,
                "unit": 1000000000000,
                "confirmations_required": 1
            }
        }
    }
}`
	confResp, err := configRespProcess([]byte(json))
	if err != nil {
		t.Errorf("Failed.")
	}
	arr, err := rateWithErrProcess(confResp)
	if len(arr) != 4 {
		t.Errorf("Failed.")
	}
	for _, currency := range arr {
		if !currency.Enabled || len(currency.ShortName) < 3 || len(currency.FullName) < 3 || len(currency.Rate) < 2 || currency.UnitPower < 5 {
			t.Errorf("Failed.")
		}
	}
}

func TestBindRespProcess(t *testing.T) {
	json := `{
    "errmsg": "",
    "code": 0,
    "data": {
        "tokenAddress": "2do3K1YLMy3Aq6EcPMdncEurP5BfAUdFPJj",
        "tokenType": "SKY"
    }
}`
	res, err := bindRespProcess([]byte(json), "SKY")
	if err != nil {
		t.Errorf("Failed.")
	}
	if res != "2do3K1YLMy3Aq6EcPMdncEurP5BfAUdFPJj" {
		t.Errorf("Failed.")
	}
}
