package postgresql

import (
	"math/rand"
	"testing"
	"time"

	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestTrackingCodeWithNilRef(t *testing.T) {
	config := config.GetServerConfig()
	db := db.OpenDb(&config.Db)
	defer db.Close()
	tx, _ := db.Begin()
	defer tx.Rollback()
	addr := randStringRunes(34)
	id := GenerateTrackingCode(tx, addr, "")
	if id < 1 {
		t.Errorf("Failed. Got %d <1.", id)
	}
	id2, refAddr := GetTrackingCode(tx, addr)
	if id != id2 {
		t.Errorf("Failed. Got %d, expected %d.", id2, id)
	}
	if len(refAddr) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr)
	}
	addr2, refAddr2 := GetAddrById(tx, id)
	if addr != addr2 {
		t.Errorf("Failed. Got %s, expected %s.", addr, addr2)
	}
	if len(refAddr2) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr2)
	}
	//更新数据
	stmt, err := tx.Prepare("DELETE FROM TRACKING_CODE where ADDRESS=$1")
	checkErr(err)

	_, err = stmt.Exec(addr)
	checkErr(err)
	stmt.Close()
	id2, refAddr = GetTrackingCode(tx, addr)
	if 0 != id2 {
		t.Errorf("Failed. Got %d, expected %d.", id2, 0)
	}
	if len(refAddr) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr)
	}
	addr2, refAddr2 = GetAddrById(tx, id)
	if len(addr2) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", addr2)
	}
	if len(refAddr2) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr2)
	}
}

func TestTrackingCodeWithRef(t *testing.T) {
	config := config.GetServerConfig()
	db := db.OpenDb(&config.Db)
	defer db.Close()
	tx, _ := db.Begin()
	defer tx.Rollback()
	addr := randStringRunes(34)
	refAddr0 := randStringRunes(34)
	GenerateTrackingCode(tx, refAddr0, "")
	refAddr := refAddr0
	id := GenerateTrackingCode(tx, addr, refAddr)
	if id < 1 {
		t.Errorf("Failed. Got %d <1.", id)
	}
	id2, refAddr := GetTrackingCode(tx, addr)
	if id != id2 {
		t.Errorf("Failed. Got %d, expected %d.", id2, id)
	}
	if len(refAddr) == 0 {
		t.Errorf("Failed. Got blank string: %s.", refAddr)
	}

	addr2, refAddr2 := GetAddrById(tx, id)
	if addr != addr2 {
		t.Errorf("Failed. Got %s, expected %s.", addr, addr2)
	}
	if len(refAddr2) == 0 {
		t.Errorf("Failed. Got blank string: %s.", refAddr2)
	}
	//更新数据
	stmt, err := db.Prepare("DELETE FROM TRACKING_CODE where ADDRESS=$1")
	checkErr(err)

	_, err = stmt.Exec(addr)
	checkErr(err)

	_, err = stmt.Exec(refAddr0)
	checkErr(err)
	stmt.Close()
}
