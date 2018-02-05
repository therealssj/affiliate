package postgresql

import (
	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
	"testing"
	"time"
)

func TestDepositRecord(t *testing.T) {
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
	mappingId := SaveDepositAddrMapping(tx, address, currencyType, ref, depositAddr)
	if mappingId < 1 {
		t.Errorf("Failed. SaveDepositAddrMapping error")
	}
	dr := db.DepositRecord{MappingId: mappingId, Seq: 100000000000, UpdatedAt: time.Now().Unix(), TransactionId: "testTransactionId", DepositAmount: 100, BuyAmount: 10000, Rate: "100", Height: 101, CurrencyType: currencyType, DepositAddr: "testDepositAddr", BuyAddr: "testBuyAddr", SuperiorRefAddr: "testSuperiorRefAddr", RefAddr: "testRefAddr"}
	SaveDepositRecord(tx, &dr)
	depositId := dr.Id
	if depositId < 1 {
		t.Errorf("Failed. SaveDepositRecord error")
	}
	if len(QueryDepositRecord(tx, dr.BuyAddr, dr.CurrencyType)) != 1 {
		t.Errorf("Failed. QueryDepositRecord error")
	}
	if len(QueryDepositRecordByAddr(tx, dr.BuyAddr)) != 1 {
		t.Errorf("Failed. QueryDepositRecord error")
	}
	if SumSalesVolume(tx, dr.RefAddr, 0.5) != dr.BuyAmount {
		t.Errorf("Failed. SumSalesVolume error")
	}
	if SumSalesVolume(tx, dr.SuperiorRefAddr, 0.5) != uint64(float64(dr.BuyAmount)*0.5) {
		t.Errorf("Failed. SumSalesVolume error")
	}
}
