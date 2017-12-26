package postgresql

import (
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
	"testing"
)

func TestAddBatchCryptocurrency(t *testing.T) {
	config := config.GetServerConfig()
	db := db.OpenDb(&config.Db)
	defer db.Close()
	tx, _ := db.Begin()
	defer tx.Rollback()
	alls := AllCryptocurrency(tx)
	sli := make([]*db.CryptocurrencyInfo, 0, 10)
	types := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		str := randStringRunes(5)
		types = append(types, str)
		sli = append(sli, &db.CryptocurrencyInfo{str, str, float32(1 + i)})
	}
	AddBatchCryptocurrency(tx, sli...)
	if len(AllCryptocurrency(tx)) != len(alls)+5 {
		t.Errorf("Failed. count error")
	}
	stmt, err := db.Prepare("DELETE FROM ALL_CRYPTOCURRENCY where SHORT_NAME in ($1, $2, $3, $4, $5)")
	checkErr(err)
	_, err = stmt.Exec(types[0], types[1], types[2], types[3], types[4])
	checkErr(err)
	stmt.Close()
}

func TestAddCryptocurrency(t *testing.T) {
	config := config.GetServerConfig()
	db := db.OpenDb(&config.Db)
	defer db.Close()
	tx, _ := db.Begin()
	defer tx.Rollback()
	alls := AllCryptocurrency(tx)
	sli := make([]*db.CryptocurrencyInfo, 0, 10)
	types := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		str := randStringRunes(5)
		types = append(types, str)
		sli = append(sli, &db.CryptocurrencyInfo{str, str, float32(1 + i)})
	}
	for _, info := range sli {
		AddCryptocurrency(tx, info.ShortName, info.FullName, info.Rate)
	}
	if len(AllCryptocurrency(tx)) != len(alls)+5 {
		t.Errorf("Failed. count error")
	}
	for i := 0; i < 5; i++ {
		rate, ok := GetRate(tx, types[i])
		if !ok || float32(i+1) != rate {
			t.Errorf("Failed. rate error")
		}
	}
	stmt, err := db.Prepare("DELETE FROM ALL_CRYPTOCURRENCY where SHORT_NAME in ($1, $2, $3, $4, $5)")
	checkErr(err)

	_, err = stmt.Exec(types[0], types[1], types[2], types[3], types[4])
	checkErr(err)
	stmt.Close()
}
