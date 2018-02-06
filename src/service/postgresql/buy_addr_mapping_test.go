package postgresql

import (
	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
	"testing"
)

func TestBuyAddrMapping(t *testing.T) {
	config := config.GetServerConfig()
	dbo := db.OpenDb(&config.Db)
	defer dbo.Close()
	tx, _ := dbo.Begin()
	defer tx.Rollback()
	var (
		address      = "testaddress"
		currencyType = "testcoin"
		ref          = "testref"
		depositAddr  = "testdepositAddr"
	)
	if GetCryptocurrency(tx, currencyType) == nil {
		AddBatchCryptocurrency(tx, []db.CryptocurrencyInfo{db.CryptocurrencyInfo{currencyType, currencyType, "100", 6, true}})
	}
	id := SaveDepositAddrMapping(tx, address, currencyType, ref, depositAddr)
	if id < 1 {
		t.Errorf("Failed. SaveDepositAddrMapping error")
	}
	if QueryMappingDepositAddr(tx, address, currencyType) == nil {
		t.Errorf("Failed. QueryMappingDepositAddr error")
	}
	stmt, err := tx.Prepare("DELETE FROM BUY_ADDR_MAPPING where id=$1")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)
	stmt.Close()
	id = SaveDepositAddrMapping(tx, address, currencyType, "", depositAddr)
	if id < 1 {
		t.Errorf("Failed. SaveDepositAddrMapping error")
	}
	if QueryMappingDepositAddr(tx, address, currencyType) == nil {
		t.Errorf("Failed. QueryMappingDepositAddr error")
	}
	stmt, err = tx.Prepare("DELETE FROM BUY_ADDR_MAPPING where id=$1")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)
	stmt.Close()

}
