package postgresql

import (
	"testing"
	"time"

	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
)

func TestRewardRecord(t *testing.T) {
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
	dr := db.DepositRecord{MappingId: mappingId, Seq: 100000000000, UpdatedAt: time.Now().Unix(), TransactionId: "testTransactionId", DepositAmount: 100, BuyAmount: 10000, Rate: "100", Height: 101, CurrencyType: currencyType, DepositAddr: "testDepositAddr", BuyAddr: "testBuyAddr"}
	SaveDepositRecord(tx, &dr)
	depositId := dr.Id
	if depositId < 1 {
		t.Errorf("Failed. SaveDepositRecord error")
	}
	data := make([]db.RewardRecord, 0, 4)
	data = append(data, db.RewardRecord{DepositId: depositId, Address: "testBuyAddr", CalAmount: 1100000, SentAmount: 1000000, RewardType: db.RewardBuyer})
	data = append(data, db.RewardRecord{DepositId: depositId, Address: "testBuyAddr", CalAmount: 2100000, SentAmount: 2000000, RewardType: db.RewardPromoter})
	ids := SaveBatchRewardRecord(tx, config.Db.ChecksumToken, data)
	if len(ids) != len(data) || ids[0] < 1 || ids[1] < 1 {
		t.Errorf("Failed. SaveBatchRewardRecord error")
	}
	if len(QueryRewardRecord(tx, config.Db.ChecksumToken, data[0].Address)) != 2 {
		t.Errorf("Failed. QueryRewardRecord error")
	}
	if len(GetUnsentRewardRecord(tx, config.Db.ChecksumToken)) != 2 {
		t.Errorf("Failed. GetUnsentRewardRecord error")
	}
	UpdateBatchSentRewardRecord(tx, config.Db.ChecksumToken, ids...)
	if len(GetUnsentRewardRecord(tx, config.Db.ChecksumToken)) != 0 {
		t.Errorf("Failed. GetUnsentRewardRecord error")
	}
	data = QueryRewardRecord(tx, config.Db.ChecksumToken, data[0].Address)
	if len(data) != 2 || !data[0].Sent || !data[1].Sent {
		t.Errorf("Failed. QueryRewardRecord error")
	}
	updatRewardRecordChecksum(tx, "wrong-checksum-token", &(data[0]))
	data = QueryRewardRecord(tx, config.Db.ChecksumToken, data[0].Address)
	if len(data) != 1 {
		t.Errorf("Failed. QueryRewardRecord error")
	}
}
