package spo_client

import (
	"encoding/json"
	"io/ioutil"
	//	"math/rand"

	"net/http"

	"github.com/spolabs/affiliate/src/config"
	//	"github.com/shopspring/decimal"
)

type balanceItem struct {
	Coins uint64 `json:"coins"`
	Hours uint64 `json:"hours"`
}

type balanceResp struct {
	Confirmed balanceItem `json:"confirmed"`
	Predicted balanceItem `json:"predicted"`
}

func httpGet(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Balance(addr string) (uint64, error) {
	conf := config.GetServerConfig()
	if conf.TestMode {
		return 9527000000, nil
	}
	resp, err := httpGet("http://localhost:8620/balance?addrs=" + addr)
	if err != nil {
		return 0, err
	}
	return balanceRespProcess(resp)
}

func balanceRespProcess(response []byte) (uint64, error) {
	info := new(balanceResp)
	err := json.Unmarshal(response, &info)
	if err != nil {
		return 0, err
	}
	return info.Confirmed.Coins, nil
}
