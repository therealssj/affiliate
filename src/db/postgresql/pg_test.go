package postgresql

import (
	"math/rand"
	"testing"
	"time"
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
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("ERROR: %s", err)
		}
	}()
	db := open()
	defer db.Close()
	addr := randStringRunes(34)
	id := GenerateTrackingCode(addr, emptyStr)
	if id < 1 {
		t.Errorf("Failed. Got %d <1.", id)
	}
	id2, refAddr := GetTrackingCode(addr)
	if id != id2 {
		t.Errorf("Failed. Got %d, expected %d.", id2, id)
	}
	if len(refAddr) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr)
	}
	addr2, refAddr2 := GetAddrById(id)
	if addr != addr2 {
		t.Errorf("Failed. Got %s, expected %s.", addr, addr2)
	}
	if len(refAddr2) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr2)
	}
	//更新数据
	stmt, err := db.Prepare("DELETE FROM TRACKING_CODE where ADDRESS=$1")
	checkErr(err)

	_, err = stmt.Exec(addr)
	checkErr(err)
	id2, refAddr = GetTrackingCode(addr)
	if 0 != id2 {
		t.Errorf("Failed. Got %d, expected %d.", id2, 0)
	}
	if len(refAddr) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr)
	}
	addr2, refAddr2 = GetAddrById(id)
	if len(addr2) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", addr2)
	}
	if len(refAddr2) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr2)
	}
}

func TestTrackingCodeWithRef(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("ERROR: %s", err)
		}
	}()
	db := open()
	defer db.Close()
	addr := randStringRunes(34)
	refAddr0 := randStringRunes(34)
	GenerateTrackingCode(refAddr0, emptyStr)
	refAddr := refAddr0
	id := GenerateTrackingCode(addr, refAddr)
	if id < 1 {
		t.Errorf("Failed. Got %d <1.", id)
	}
	id2, refAddr := GetTrackingCode(addr)
	if id != id2 {
		t.Errorf("Failed. Got %d, expected %d.", id2, id)
	}
	if len(refAddr) == 0 {
		t.Errorf("Failed. Got blank string: %s.", refAddr)
	}

	addr2, refAddr2 := GetAddrById(id)
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
	id2, refAddr = GetTrackingCode(addr)
	if 0 != id2 {
		t.Errorf("Failed. Got %d, expected %d.", id2, 0)
	}
	if len(refAddr) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr)
	}
	addr2, refAddr2 = GetAddrById(id)
	if len(addr2) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", addr2)
	}
	if len(refAddr2) != 0 {
		t.Errorf("Failed. Got not blank string: %s.", refAddr2)
	}
	stmt, err = db.Prepare("DELETE FROM TRACKING_CODE where ADDRESS=$1")
	checkErr(err)

	_, err = stmt.Exec(refAddr0)
	checkErr(err)
}
