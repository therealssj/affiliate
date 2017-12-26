package postgresql

import (
	"database/sql"

	"github.com/spaco/affiliate/src/service/db"
)

func SaveBatchDepositRecord(tx *sql.Tx, batch ...db.DepositRecord) {
	var lastInsertId uint64
	for _, dr := range batch {
		err := tx.QueryRow("insert into DEPOSIT_RECORD(CREATION,MAPPING_ID,SEQ,UPDATED_AT,TRANSACTION_ID,DEPOSIT_AMOUNT,BUY_AMOUNT,RATE,HEIGHT,BUY_ADDR,REF_ADDR,SUPERIOR_REF_ADDR) values (now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id;", dr.MappingId, dr.Seq, dr.UpdatedAt, dr.TransactionId, dr.DepositAmount, dr.BuyAmount, dr.Rate, dr.Height, dr.BuyAddr, dr.RefAddr, dr.SuperiorRefAddr).Scan(&lastInsertId)
		checkErr(err)
		dr.Id = lastInsertId
	}
}

func sumL1PromoteSalesVolume(tx *sql.Tx, address string) uint64 {
	rows, err := tx.Query("SELECT sum(BUY_AMOUNT) FROM DEPOSIT_RECORD where REF_ADDR=$1", address)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var resNullable sql.NullInt64
		err = rows.Scan(&resNullable)
		checkErr(err)
		if resNullable.Valid {
			return uint64(resNullable.Int64)
		}
		return 0
	}
	return 0
}

func sumL2PromoteSalesVolume(tx *sql.Tx, address string) uint64 {
	rows, err := tx.Query("SELECT sum(BUY_AMOUNT) FROM DEPOSIT_RECORD where SUPERIOR_REF_ADDR=$1", address)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var resNullable sql.NullInt64
		err = rows.Scan(&resNullable)
		checkErr(err)
		if resNullable.Valid {
			return uint64(resNullable.Int64)
		}
		return 0
	}
	return 0
}

func SumSalesVolume(tx *sql.Tx, address string, superiorRatio float64) uint64 {
	return sumL1PromoteSalesVolume(tx, address) + uint64(float64(sumL2PromoteSalesVolume(tx, address))*superiorRatio)
}
