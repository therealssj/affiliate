package service

import (
	"database/sql"
	"fmt"

	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
	pg "github.com/spaco/affiliate/src/service/postgresql"
	"github.com/spaco/affiliate/src/tracking_code"
)

func SyncCryptocurrency(newCurrency []db.CryptocurrencyInfo, updateRateCur []db.CryptocurrencyInfo) {
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

func getRewardRemain(tx *sql.Tx, batch []db.DepositRecord) map[string]uint64 {
	addrs = make([]string, 0, 2*len(batch))
	for _, dr := range batch {
		if len(dr.RefAddr) > 0 {
			addrs = append(addrs, dr.BuyAddr, dr.RefAddr)
			if len(dr.SuperiorRefAddr) > 0 {
				addrs = append(addrs, dr.SuperiorRefAddr)
	return pg.QueryRewardRemain(tx, addrs...)
}

func ProcessDeposit(batch []db.DepositRecord, req int64) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	rewardConfig := config.GetDaemonConfig().RewardConfig
	rewardRecords := make([]db.RewardRecord, 0, 3*len(batch))
	remainMap := getRewardRemain(tx, batch)
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
		pg.SaveBatchDepositRecord(tx, dr)
		if len(dr.RefAddr) > 0 {
			//reward buyer
			rewardAmount := uint64(float64(dr.BuyAmount) * rewardConfig.BuyerRate)
			if rm, ok := remainMap[dr.BuyAddr]; ok {
				rewardAmount += rm
			}
			remain := rewardAmount % uint64(rewardConfig.MinSendAmount)
			remainMap[dr.BuyAddr] = remain
			rewardRecords = append(rewardRecords, db.RewardRecord{DepositId: dr.Id,
				Address:    dr.BuyAddr,
				CalAmount:  rewardAmount,
				SentAmount: rewardAmount - remain,
				RewardType: db.RewardBuyer})
			//reward promoter
			ratio, _ := getPromoterRatio(tx, &rewardConfig, dr.RefAddr)
			rewardAmount = uint64(float64(dr.BuyAmount) * ratio)
			if rm, ok := remainMap[dr.RefAddr]; ok {
				rewardAmount += rm
			}
			remain = rewardAmount % uint64(rewardConfig.MinSendAmount)
			remainMap[dr.RefAddr] = remain
			rewardRecords = append(rewardRecords, db.RewardRecord{DepositId: dr.Id,
				Address:    dr.RefAddr,
				CalAmount:  rewardAmount,
				SentAmount: rewardAmount - remain,
				RewardType: db.RewardPromoter})
			if len(dr.SuperiorRefAddr) > 0 {
				//reward promoter
				_, ratio = getPromoterRatio(tx, &rewardConfig, dr.SuperiorRefAddr)
				rewardAmount = uint64(float64(dr.BuyAmount) * ratio)
				if rm, ok := remainMap[dr.SuperiorRefAddr]; ok {
					rewardAmount += rm
				}
				remain = rewardAmount % uint64(rewardConfig.MinSendAmount)
				remainMap[dr.SuperiorRefAddr] = remain
				rewardRecords = append(rewardRecords, db.RewardRecord{DepositId: dr.Id,
					Address:    dr.SuperiorRefAddr,
					CalAmount:  rewardAmount,
					SentAmount: rewardAmount - remain,
					RewardType: db.RewardSuperiorPromoter})
			}
		}
	}
	pg.SaveBatchRewardRecord(tx, rewardRecords...)
	pg.UpdateRewardRemain(tx, remainMap)
	pg.SaveKvStore(tx, tellerReqName, req, "")
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

func GetUnsentRewardRecord() []db.RewardRecord {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	rrs := pg.GetUnsentRewardRecord(tx)
	checkErr(tx.Commit())
	commit = true
	return rrs
}

func UpdateBatchRewardRecord(ids []uint64) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	pg.UpdateBatchRewardRecord(tx, ids...)
	commit = true
}
