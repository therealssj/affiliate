package postgresql

import (
	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
	"testing"
)

func TestRewardRemain(t *testing.T) {
	config := config.GetServerConfig()
	dbo := db.OpenDb(&config.Db)
	defer dbo.Close()
	tx, _ := dbo.Begin()
	defer tx.Rollback()
	m := make(map[string]uint64, 4)
	m["testaddress1"] = 101
	m["testaddress2"] = 234
	m["testaddress3"] = 382
	UpdateRewardRemain(tx, m)
	m2 := QueryRewardRemain(tx, mapKeySlice(m)...)
	for k, v := range m {
		if m2[k] != v {
			t.Errorf("Failed. key %s, expect:%d, actual:%d", k, v, m2[k])
		}
	}
	m["testaddress3"] = 395
	m["testaddress4"] = 428
	UpdateRewardRemain(tx, m)
	m2 = QueryRewardRemain(tx, mapKeySlice(m)...)
	for k, v := range m {
		if m2[k] != v {
			t.Errorf("Failed. key %s, expect:%d, actual:%d", k, v, m2[k])
		}
	}
}

func mapKeySlice(m map[string]uint64) []string {
	res := make([]string, 0, len(m))
	for k, _ := range m {
		res = append(res, k)
	}
	return res
}
