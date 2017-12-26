package service

import (
	"database/sql"
	"fmt"
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
	pg "github.com/spaco/affiliate/src/service/postgresql"
	"github.com/spaco/affiliate/src/tracking_code"
)

func SyncCryptocurrency(newCurrency []*db.CryptocurrencyInfo, updateRateCur []*db.CryptocurrencyInfo) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	pg.AddBatchCryptocurrency(tx, newCurrency...)
	pg.UpdateBatchRate(tx, updateRateCur...)
	checkErr(tx.Commit())
	commit = true
}

const tellerReqName = "teller:req"

func GetTellerReq() int64 {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	intVal, _, _ := pg.GetKvStore(tx, tellerReqName)
	checkErr(tx.Commit())
	commit = true
	return intVal
}

func UpdateTellerReq(val int64) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	pg.SaveKvStore(tx, tellerReqName, val, "")
	checkErr(tx.Commit())
	commit = true
}

func ProcessDeposit(batch []db.DepositRecord, req int64) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	for _, dr := range batch {
		mapping, found := pg.QueryMappingDepositAddr(tx, dr.BuyAddr, dr.CurrencyType)
		if !found {
			panic(fmt.Sprintf("not found BuyAddrMapping for address:%s CurrencyType:%s", dr.BuyAddr, dr.CurrencyType))
		}
		dr.MappingId = mapping.Id
		if len(mapping.Ref) > 0 {
			if id := tracking_code.GetId(mapping.Ref); id > 0 {
				dr.RefAddr, dr.SuperiorRefAddr = pg.GetAddrById(tx, id)
			}
		}
	}
	pg.SaveBatchDepositRecord(tx, batch...)
	//	rewardConfig := config.GetDaemonConfig().RewardConfig
	rewardRecords := make([]db.RewardRecord, 0, 3*len(batch))
	var rewardRecord db.RewardRecord
	for _, dr := range batch {
		if len(dr.RefAddr) > 0 {
			rewardRecords = append(rewardRecords, rewardRecord)
			//			SumSalesVolume()
		}
	}

	// calculator reward
	commit = true
}

func getPromoterRatio(tx *sql.Tx, rewardConfig *config.RewardConfig, address string) (float64, float64) {
	sv := pg.SumSalesVolume(tx, address, rewardConfig.SuperiorDiscount)
	i := len(rewardConfig.LadderLine) - 1
	for ; i >= 0; i-- {
		if sv >= uint64(rewardConfig.LadderLine[i]) {
			break
		}
	}
	return rewardConfig.PromoterRatio[i], rewardConfig.SuperiorPromoterRatio[i]
}

func SendReward() {

}
