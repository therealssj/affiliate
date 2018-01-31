package teller_client

import (
	//	"fmt"
	"testing"
)

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
