package postgresql

import (
	"database/sql"
	"strings"
	"time"

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

func GetUnsentRewardRecord(tx *sql.Tx) []db.RewardRecord {
	rows, err := tx.Query("SELECT ID,VERSION,CREATION,DEPOSIT_ID,ADDRESS,CAL_AMOUNT,SENT_AMOUNT,REWARD_TYPE FROM REWARD_RECORD where SENT=false and SentAmount<>0")
	checkErr(err)
	res := make([]db.RewardRecord, 0, 16)
	defer rows.Close()
	for rows.Next() {
		var id, version, depositId, calAmount, sentAmount uint64
		var address, rewardType string
		var creation time.Time
		err = rows.Scan(&id, &version, &creation, &depositId, &address, &calAmount, &sentAmount, &rewardType)
		checkErr(err)
		res = append(res, db.RewardRecord{Id: id,
			Version:    version,
			Creation:   creation,
			DepositId:  depositId,
			Address:    address,
			CalAmount:  calAmount,
			SentAmount: sentAmount,
			RewardType: rewardType})
	}
	return res
}

func UpdateBatchRewardRecord(tx *sql.Tx, ids ...uint64) {
	stmt, err := tx.Prepare("update REWARD_RECORD set SENT_TIME=now(),SENT=true where ID in (?" + strings.Repeat(",?", len(ids)-1) + ")")
	defer stmt.Close()
	checkErr(err)
	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	_, err = stmt.Exec(args...)
	checkErr(err)
}
