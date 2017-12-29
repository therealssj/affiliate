package postgresql

import (
	"fmt"
	"testing"

	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
)

func TestAddBatchCryptocurrency(t *testing.T) {
	config := config.GetServerConfig()
	dbo := db.OpenDb(&config.Db)
	defer dbo.Close()
	tx, _ := dbo.Begin()
	defer tx.Rollback()
	alls := AllCryptocurrency(tx)
	sli := make([]db.CryptocurrencyInfo, 0, 10)
	types := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		str := randStringRunes(5)
		types = append(types, str)
		sli = append(sli, db.CryptocurrencyInfo{str, str, fmt.Sprintf("%d", i+1), 6})
	}
	AddBatchCryptocurrency(tx, sli)
	if len(AllCryptocurrency(tx)) != len(alls)+5 {
		t.Errorf("Failed. count error")
	}
	stmt, err := dbo.Prepare("DELETE FROM ALL_CRYPTOCURRENCY where SHORT_NAME in ($1, $2, $3, $4, $5)")
	checkErr(err)
	_, err = stmt.Exec(types[0], types[1], types[2], types[3], types[4])
	checkErr(err)
	stmt.Close()
}

func TestAddCryptocurrency(t *testing.T) {
	config := config.GetServerConfig()
	dbo := db.OpenDb(&config.Db)
	defer dbo.Close()
	tx, _ := dbo.Begin()
	defer tx.Rollback()
	alls := AllCryptocurrency(tx)
	sli := make([]db.CryptocurrencyInfo, 0, 10)
	types := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		str := randStringRunes(5)
		types = append(types, str)
		sli = append(sli, db.CryptocurrencyInfo{str, str, fmt.Sprintf("%d", i+1), 6})
	}
	AddBatchCryptocurrency(tx, sli)
	if len(AllCryptocurrency(tx)) != len(alls)+5 {
		t.Errorf("Failed. count error")
	}
	for i := 0; i < 5; i++ {
		rate, ok := GetRate(tx, types[i])
		if !ok || fmt.Sprintf("%d", i+1) != rate {
			t.Errorf("Failed. rate error")
		}
	}
	stmt, err := dbo.Prepare("DELETE FROM ALL_CRYPTOCURRENCY where SHORT_NAME in ($1, $2, $3, $4, $5)")
	checkErr(err)

	_, err = stmt.Exec(types[0], types[1], types[2], types[3], types[4])
	checkErr(err)
	stmt.Close()
}
