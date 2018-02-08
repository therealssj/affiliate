package postgresql

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"github.com/spolabs/affiliate/src/service/db"
)

func rewardRemainCheckSum(checksumToken string, address string, amount uint64) string {
	// ADDRESS AMOUNT
	hash := hmac.New(sha256.New, []byte(checksumToken))
	hash.Write([]byte(address))
	hash.Write(uint64ToByteArray(amount))
	return hex.EncodeToString(hash.Sum(nil))
}

func UpdateRewardRemain(tx *sql.Tx, checksumToken string, data map[string]uint64) {
	stmt, err := tx.Prepare("INSERT INTO REWARD_REMAIN (ADDRESS, CREATION, LAST_MODIFIED, AMOUNT, CHECKSUM) VALUES ($1,now(),now(),$2, $3) ON CONFLICT (ADDRESS) DO UPDATE SET AMOUNT=$2, CHECKSUM=$3")
	defer stmt.Close()
	checkErr(err)
	for k, v := range data {
		_, err = stmt.Exec(k, v, rewardRemainCheckSum(checksumToken, k, v))
		checkErr(err)
	}
}

func QueryRewardRemain(tx *sql.Tx, checksumToken string, addr ...string) map[string]uint64 {
	m := make(map[string]uint64, len(addr))
	args := make([]interface{}, 0, len(addr))
	for _, ad := range addr {
		args = append(args, ad)
	}
	rows, err := tx.Query("select ADDRESS,AMOUNT,CHECKSUM from REWARD_REMAIN where ADDRESS in "+db.InClause(len(addr), 1), args...)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var address, checksum string
		var amount uint64
		err = rows.Scan(&address, &amount, &checksum)
		checkErr(err)
		if rewardRemainCheckSum(checksumToken, address, amount) == checksum {
			m[address] = amount
		}
	}
	return m
}
