package postgresql

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/spolabs/affiliate/src/service/db"
)

func rewardRecordCheckSum(checksumToken string, reward *db.RewardRecord) string {
	// ID DEPOSIT_ID ADDRESS SENT_AMOUNT SENT
	hash := hmac.New(sha256.New, []byte(checksumToken))
	hash.Write(uint64ToByteArray(reward.Id))
	hash.Write(uint64ToByteArray(reward.DepositId))
	hash.Write([]byte(reward.Address))
	hash.Write(uint64ToByteArray(reward.SentAmount))
	if reward.Sent {
		hash.Write([]byte{byte(1)})
	} else {
		hash.Write([]byte{byte(0)})
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func updatRewardRecordChecksum(tx *sql.Tx, checksumToken string, reward *db.RewardRecord) {
	stmt, err := tx.Prepare("update REWARD_RECORD set CHECKSUM=$1 where ID=$2")
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec(rewardRecordCheckSum(checksumToken, reward), reward.Id)
	checkErr(err)
}

func SaveBatchRewardRecord(tx *sql.Tx, checksumToken string, batch []db.RewardRecord) []uint64 {
	res := make([]uint64, 0, len(batch))
	var lastInsertId uint64
	for _, ar := range batch {
		err := tx.QueryRow("insert into REWARD_RECORD(CREATION,DEPOSIT_ID,ADDRESS,CAL_AMOUNT,SENT_AMOUNT,SENT,REWARD_TYPE) values (now(),$1,$2,$3,$4,false,$5) returning ID", ar.DepositId, ar.Address, ar.CalAmount, ar.SentAmount, ar.RewardType).Scan(&lastInsertId)
		checkErr(err)
		res = append(res, lastInsertId)
		ar.Id = lastInsertId
		updatRewardRecordChecksum(tx, checksumToken, &ar)
	}
	return res
}

func GetUnsentRewardRecord(tx *sql.Tx, checksumToken string) []db.RewardRecord {
	rows, err := tx.Query("SELECT ID,CREATION,DEPOSIT_ID,ADDRESS,CAL_AMOUNT,SENT_AMOUNT,REWARD_TYPE,CHECKSUM FROM REWARD_RECORD where SENT=false and SENT_AMOUNT<>0 order by ID")
	checkErr(err)
	res := make([]db.RewardRecord, 0, 16)
	defer rows.Close()
	for rows.Next() {
		var id, depositId, calAmount, sentAmount uint64
		var address, rewardType, checksum string
		var creation time.Time
		err = rows.Scan(&id, &creation, &depositId, &address, &calAmount, &sentAmount, &rewardType, &checksum)
		checkErr(err)
		reward := db.RewardRecord{Id: id,
			Creation:   creation,
			DepositId:  depositId,
			Address:    address,
			CalAmount:  calAmount,
			SentAmount: sentAmount,
			RewardType: rewardType}
		if rewardRecordCheckSum(checksumToken, &reward) == checksum {
			res = append(res, reward)
		}
	}
	return res
}

func UpdateBatchSentRewardRecord(tx *sql.Tx, checksumToken string, ids ...uint64) {
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		reward, _ := getRewardRecord(tx, id)
		updateRewardRecordSentAndChecksum(tx, checksumToken, reward)
	}
	// stmt, err := tx.Prepare("update REWARD_RECORD set SENT_TIME=now(),SENT=true where ID in " + db.InClause(len(ids), 1))
	// defer stmt.Close()
	// checkErr(err)
	// args := make([]interface{}, 0, len(ids))
	// for _, id := range ids {
	// 	args = append(args, id)
	// }
	// _, err = stmt.Exec(args...)
	// checkErr(err)

}

func updateRewardRecordSentAndChecksum(tx *sql.Tx, checksumToken string, reward *db.RewardRecord) {
	stmt, err := tx.Prepare("update REWARD_RECORD set SENT_TIME=now(),SENT=true,CHECKSUM=$1 where ID=$2")
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec(rewardRecordCheckSum(checksumToken, reward), reward.Id)
	checkErr(err)
}

func getRewardRecord(tx *sql.Tx, id uint64) (*db.RewardRecord, string) {
	rows, err := tx.Query("SELECT ID,CREATION,DEPOSIT_ID,ADDRESS,CAL_AMOUNT,SENT_AMOUNT,SENT_TIME,SENT,REWARD_TYPE,CHECKSUM FROM REWARD_RECORD where ID=$1", id)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		return buildRewardRecord(rows)
	}
	return nil, ""

}

func buildRewardRecord(rows *sql.Rows) (*db.RewardRecord, string) {
	var id, depositId, calAmount, sentAmount uint64
	var address, rewardType, checksum string
	var creation time.Time
	var sentTime db.Time
	var sentTimeNullable db.NullTime
	var sent bool
	err := rows.Scan(&id, &creation, &depositId, &address, &calAmount, &sentAmount, &sentTimeNullable, &sent, &rewardType, &checksum)
	checkErr(err)
	if sentTimeNullable.Valid {
		sentTime = sentTimeNullable.Time
	}
	return &db.RewardRecord{Id: id,
		Creation:   creation,
		DepositId:  depositId,
		Address:    address,
		CalAmount:  calAmount,
		SentAmount: sentAmount,
		SentTime:   sentTime,
		Sent:       sent,
		RewardType: rewardType}, checksum
}

func QueryRewardRecord(tx *sql.Tx, checksumToken string, address string) []db.RewardRecord {
	rows, err := tx.Query("SELECT ID,CREATION,DEPOSIT_ID,ADDRESS,CAL_AMOUNT,SENT_AMOUNT,SENT_TIME,SENT,REWARD_TYPE,CHECKSUM FROM REWARD_RECORD where ADDRESS=$1 order by ID", address)
	checkErr(err)
	res := make([]db.RewardRecord, 0, 16)
	defer rows.Close()
	for rows.Next() {
		reward, checksum := buildRewardRecord(rows)
		if rewardRecordCheckSum(checksumToken, reward) == checksum {
			res = append(res, *reward)
		}
	}
	return nil
}
