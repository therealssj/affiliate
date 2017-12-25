package postgresql

import (
	"database/sql"

	"github.com/spaco/affiliate/src/service/db"
)

func SaveBatchDepositRecord(tx *sql.Tx, batch ...db.DepositRecord) []uint64 {
	res := make([]uint64, 0, len(batch))
	var lastInsertId uint64
	for _, dr := range batch {
		err := tx.QueryRow("insert into DEPOSIT_RECORD(CREATION,MAPPING_ID,SEQ,UPDATED_AT,TRANSACTION_ID,DEPOSIT_AMOUNT,BUY_AMOUNT,RATE,HEIGHT) values (now(), $1, $2, $3, $4, $5, $6, $7, $8) returning id;", dr.MappingId, dr.Seq, dr.UpdatedAt, dr.TransactionId, dr.DepositAmount, dr.BuyAmount, dr.Rate, dr.Height).Scan(&lastInsertId)
		checkErr(err)
		res = append(res, lastInsertId)
	}
	return res
}
