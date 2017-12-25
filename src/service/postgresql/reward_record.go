package postgresql

import (
	"database/sql"

	"github.com/spaco/affiliate/src/service/db"
)

func SaveBatchRewardRecord(tx *sql.Tx, batch ...*db.RewardRecord) []uint64 {
	res := make([]uint64, 0, len(batch))
	var lastInsertId uint64
	for _, ar := range batch {
		err := tx.QueryRow("insert into REWARD_RECORD(VERSION,CREATION,DEPOSIT_ID,ADDRESS,CAL_AMOUNT,SENT_AMOUNT,SENT,REWARD_TYPE) values (0,now(),$1,$2,$3,$4,false,$5) returning id;", ar.DepositId, ar.Address, ar.CalAmount, ar.SentAmount, ar.RewardType).Scan(&lastInsertId)
		checkErr(err)
		res = append(res, lastInsertId)
	}
	return res
}
func UpdateBatchRewardRecord(tx *sql.Tx, []uint64){

	
}
