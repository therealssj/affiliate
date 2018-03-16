package postgresql

import (
	"testing"

	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
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
	UpdateRewardRemain(tx, config.Db.ChecksumToken, m)
	m2 := QueryRewardRemain(tx, config.Db.ChecksumToken, mapKeySlice(m)...)
	for k, v := range m {
		if m2[k] != v {
			t.Errorf("Failed. key %s, expect:%d, actual:%d", k, v, m2[k])
		}
	}
	m["testaddress3"] = 395
	m["testaddress4"] = 428
	UpdateRewardRemain(tx, config.Db.ChecksumToken, m)
	m2 = QueryRewardRemain(tx, config.Db.ChecksumToken, mapKeySlice(m)...)
	for k, v := range m {
		if m2[k] != v {
			t.Errorf("Failed. key %s, expect:%d, actual:%d", k, v, m2[k])
		}
	}
	if len(m2) != 4 {
		t.Errorf("Failed.")
	}
	stmt, err := tx.Prepare("update REWARD_REMAIN set CHECKSUM=$1 where ADDRESS=$2")
	checkErr(err)

	_, err = stmt.Exec("wrong-checksum", "testaddress3")
	checkErr(err)
	stmt.Close()
	m2 = QueryRewardRemain(tx, config.Db.ChecksumToken, mapKeySlice(m)...)
	if len(m2) != 3 {
		t.Errorf("Failed.")
	}
	if _, ok := m2["testaddress3"]; ok {
		t.Errorf("Failed.")
	}
	m = make(map[string]uint64, 4)
	m["k9QgadMDxisLfj2CLgNrwgzZZSEmZMeUpK"] = 878000
	UpdateRewardRemain(tx, "!@#$%^&*(){}|_+<>?", m)
	m2 = QueryRewardRemain(tx, "!@#$%^&*(){}|_+<>?", mapKeySlice(m)...)
	if m2["k9QgadMDxisLfj2CLgNrwgzZZSEmZMeUpK"] != 878000 {
		t.Errorf("Failed.")
	}
}

func mapKeySlice(m map[string]uint64) []string {
	res := make([]string, 0, len(m))
	for k, _ := range m {
		res = append(res, k)
	}
	return res
}
