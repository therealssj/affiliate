package service

import (
	"database/sql"
	"fmt"

	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
	pg "github.com/spaco/affiliate/src/service/postgresql"
	"github.com/spaco/affiliate/src/tracking_code"
)

func syncCryptocurrency(newCurrency []db.CryptocurrencyInfo, updateRateCur []db.CryptocurrencyInfo) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	pg.AddBatchCryptocurrency(tx, newCurrency)
	pg.UpdateBatchRate(tx, updateRateCur)
	checkErr(tx.Commit())
	commit = true
}

//const tellerReqName = "teller:req"

//func GetTellerReq() int64 {
//	tx, commit := db.BeginTx()
//	defer db.Rollback(tx, &commit)
//	intVal, _, _ := pg.GetKvStore(tx, tellerReqName)
//	checkErr(tx.Commit())
//	commit = true
//	return intVal
//}

func fillAndGetRewardRemain(tx *sql.Tx, batch []db.DepositRecord) map[string]uint64 {
	addrs := make([]string, 0, 2*len(batch))
	for i, _ := range batch {
		dr := &(batch[i])
		if dr.CurrencyType == "Deprecated" {
			continue
		}
		mapping := pg.QueryMappingDepositAddr(tx, dr.BuyAddr, dr.CurrencyType)
		if mapping == nil {
			panic(fmt.Sprintf("not found BuyAddrMapping for address:%s CurrencyType:%s", dr.BuyAddr, dr.CurrencyType))
		}
		dr.MappingId = mapping.Id
		if len(mapping.Ref) > 0 {
			if id := tracking_code.GetId(mapping.Ref); id > 0 {
				dr.RefAddr, dr.SuperiorRefAddr = pg.GetAddrById(tx, id)
				if len(dr.RefAddr) > 0 {
					addrs = append(addrs, dr.BuyAddr, dr.RefAddr)
					if len(dr.SuperiorRefAddr) > 0 {
						addrs = append(addrs, dr.SuperiorRefAddr)
					}
				}
			}
		}
	}
	if len(addrs) > 0 {
		return pg.QueryRewardRemain(tx, addrs...)
	}
	return make(map[string]uint64, 0)
}

//func SaveTellerReq(req int64) {
//	tx, commit := db.BeginTx()
//	defer db.Rollback(tx, &commit)
//	pg.SaveKvStore(tx, tellerReqName, req, "")
//	checkErr(tx.Commit())
//	commit = true
//}

func ProcessDeposit(batch []db.DepositRecord /*, oldReq int64, req int64*/) {
	rewardConfig := config.GetApiForTellerConfig().RewardConfig
	rewardRecords := make([]db.RewardRecord, 0, 3*len(batch))
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	remainMap := fillAndGetRewardRemain(tx, batch)
	changedRemainMap := make(map[string]uint64, len(remainMap))
	for _, dr := range batch {
		if dr.CurrencyType == "Deprecated" {
			continue
		}
		pg.SaveDepositRecord(tx, &dr)
		if len(dr.RefAddr) > 0 {
			rewardRecords = appendBuyerRewardRecord(tx, rewardRecords, &dr, &rewardConfig, remainMap, changedRemainMap)
			rewardRecords = appendPromoterRewardRecord(tx, rewardRecords, &dr, &rewardConfig, remainMap, changedRemainMap)
			if len(dr.SuperiorRefAddr) > 0 {
				rewardRecords = appendSuperiorPromoterRewardRecord(tx, rewardRecords, &dr, &rewardConfig, remainMap, changedRemainMap)
			}
		}
	}
	if len(rewardRecords) > 0 {
		pg.SaveBatchRewardRecord(tx, rewardRecords)
	}
	if len(changedRemainMap) > 0 {
		pg.UpdateRewardRemain(tx, changedRemainMap)
	}
	//	if oldReq != req {
	//		pg.SaveKvStore(tx, tellerReqName, req, "")
	//	}
	checkErr(tx.Commit())
	commit = true
}

func appendBuyerRewardRecord(tx *sql.Tx, rewardRecords []db.RewardRecord, dr *db.DepositRecord, rewardConfig *config.RewardConfig, remainMap, changedRemainMap map[string]uint64) []db.RewardRecord {
	rewardAmount := uint64(float64(dr.BuyAmount) * rewardConfig.BuyerRatio)
	if rewardAmount == 0 {
		return rewardRecords
	}
	sentAmount := rewardAmount
	if rm, ok := remainMap[dr.BuyAddr]; ok {
		sentAmount += rm
	}
	remain := sentAmount % uint64(rewardConfig.MinSendAmount)
	remainMap[dr.BuyAddr] = remain
	changedRemainMap[dr.BuyAddr] = remain
	return append(rewardRecords, db.RewardRecord{DepositId: dr.Id,
		Address:    dr.BuyAddr,
		CalAmount:  rewardAmount,
		SentAmount: sentAmount - remain,
		RewardType: db.RewardBuyer})
}

func appendPromoterRewardRecord(tx *sql.Tx, rewardRecords []db.RewardRecord, dr *db.DepositRecord, rewardConfig *config.RewardConfig, remainMap, changedRemainMap map[string]uint64) []db.RewardRecord {
	ratio, _ := getPromoterRatio(tx, rewardConfig, dr.RefAddr)
	rewardAmount := uint64(float64(dr.BuyAmount) * ratio)
	if rewardAmount == 0 {
		return rewardRecords
	}
	sentAmount := rewardAmount
	if rm, ok := remainMap[dr.RefAddr]; ok {
		sentAmount += rm
	}
	remain := sentAmount % uint64(rewardConfig.MinSendAmount)
	remainMap[dr.RefAddr] = remain
	changedRemainMap[dr.RefAddr] = remain
	return append(rewardRecords, db.RewardRecord{DepositId: dr.Id,
		Address:    dr.RefAddr,
		CalAmount:  rewardAmount,
		SentAmount: sentAmount - remain,
		RewardType: db.RewardPromoter})
}

func appendSuperiorPromoterRewardRecord(tx *sql.Tx, rewardRecords []db.RewardRecord, dr *db.DepositRecord, rewardConfig *config.RewardConfig, remainMap, changedRemainMap map[string]uint64) []db.RewardRecord {
	_, ratio := getPromoterRatio(tx, rewardConfig, dr.SuperiorRefAddr)
	rewardAmount := uint64(float64(dr.BuyAmount) * ratio)
	if rewardAmount == 0 {
		return rewardRecords
	}
	sentAmount := rewardAmount
	if rm, ok := remainMap[dr.SuperiorRefAddr]; ok {
		sentAmount += rm
	}
	remain := sentAmount % uint64(rewardConfig.MinSendAmount)
	remainMap[dr.SuperiorRefAddr] = remain
	changedRemainMap[dr.SuperiorRefAddr] = remain
	return append(rewardRecords, db.RewardRecord{DepositId: dr.Id,
		Address:    dr.SuperiorRefAddr,
		CalAmount:  rewardAmount,
		SentAmount: sentAmount - remain,
		RewardType: db.RewardSuperiorPromoter})
}

func getPromoterRatio(tx *sql.Tx, rewardConfig *config.RewardConfig, address string) (float64, float64) {
	sv := pg.SumSalesVolume(tx, address, rewardConfig.SuperiorDiscount)
	return getPromoterRatioBySalesVolume(rewardConfig, sv)
}

func getPromoterRatioBySalesVolume(rewardConfig *config.RewardConfig, sv uint64) (float64, float64) {
	i := len(rewardConfig.LadderLine) - 1
	for ; i >= 0; i-- {
		if sv >= uint64(rewardConfig.LadderLine[i]) {
			break
		}
	}
	if i < 0 {
		i = 0
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

func UpdateBatchRewardRecord(ids ...uint64) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	pg.UpdateBatchRewardRecord(tx, ids...)
	checkErr(tx.Commit())
	commit = true
}
