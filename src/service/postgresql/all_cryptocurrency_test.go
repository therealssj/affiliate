package postgresql

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
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
		sli = append(sli, db.CryptocurrencyInfo{str, str, decimal.NewFromFloat(float64(i + 1)).String(), 6, true})
	}
	AddBatchCryptocurrency(tx, sli)
	if len(AllCryptocurrency(tx)) != len(alls)+5 {
		t.Errorf("Failed. count error")
	}
	stmt, err := tx.Prepare("DELETE FROM ALL_CRYPTOCURRENCY where SHORT_NAME in ($1, $2, $3, $4, $5)")
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
		sli = append(sli, db.CryptocurrencyInfo{str, str, decimal.NewFromFloat(float64(i + 1)).String(), 6, true})
	}
	AddBatchCryptocurrency(tx, sli)
	if len(AllCryptocurrency(tx)) != len(alls)+5 {
		t.Errorf("Failed. count error")
	}
	for i := 0; i < 5; i++ {
		info := GetCryptocurrency(tx, types[i])
		if info == nil || decimal.NewFromFloat(float64(i+1)).String() != info.Rate {
			t.Errorf("Failed. rate error, expect:[%s], actual:[%s]", decimal.NewFromFloat(float64(i+1)).String(), info.Rate)
		}
	}

	stmt, err := tx.Prepare("DELETE FROM ALL_CRYPTOCURRENCY where SHORT_NAME in ($1, $2, $3, $4, $5)")
	checkErr(err)

	_, err = stmt.Exec(types[0], types[1], types[2], types[3], types[4])
	checkErr(err)
	stmt.Close()
}

func TestUpdateBatchRateAndEnabled(t *testing.T) {
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
		sli = append(sli, db.CryptocurrencyInfo{str, str, decimal.NewFromFloat(float64(i + 1)).String(), 6, true})
	}
	AddBatchCryptocurrency(tx, sli)
	if len(AllCryptocurrency(tx)) != len(alls)+5 {
		t.Errorf("Failed. count error")
	}
	for i, _ := range sli {
		sli[i].Rate = decimal.NewFromFloat(float64(10 * (i + 1))).String()
		sli[i].Enabled = false
	}
	UpdateBatchRateAndEnabled(tx, sli)
	for i := 0; i < 5; i++ {
		info := GetCryptocurrency(tx, types[i])
		if info == nil || decimal.NewFromFloat(float64(10*(i+1))).String() != info.Rate || info.Enabled {
			t.Errorf("Error")
		}
	}

}
