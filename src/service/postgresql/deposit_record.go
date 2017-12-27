package postgresql

import (
	"database/sql"
	"github.com/spaco/affiliate/src/service/db"
	"time"
)

func SaveBatchDepositRecord(tx *sql.Tx, batch ...db.DepositRecord) {
	var lastInsertId uint64
	for _, dr := range batch {
		err := tx.QueryRow("insert into DEPOSIT_RECORD(CREATION,MAPPING_ID,SEQ,UPDATED_AT,TRANSACTION_ID,DEPOSIT_AMOUNT,BUY_AMOUNT,RATE,HEIGHT,BUY_ADDR,REF_ADDR,SUPERIOR_REF_ADDR) values (now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id;", dr.MappingId, dr.Seq, time.Unix(dr.UpdatedAt, 0), dr.TransactionId, dr.DepositAmount, dr.BuyAmount, dr.Rate, dr.Height, dr.BuyAddr, dr.RefAddr, dr.SuperiorRefAddr).Scan(&lastInsertId)
		checkErr(err)
		dr.Id = lastInsertId
	}
}

func QueryDepositRecord(tx *sql.Tx, address string) []db.DepositRecord {
	rows, err := tx.Query("SELECT ID,CREATION,MAPPING_ID,BUY_ADDR,REF_ADDR,SUPERIOR_REF_ADDR,SEQ,UPDATED_AT,TRANSACTION_ID,DEPOSIT_AMOUNT,BUY_AMOUNT,RATE,HEIGHT FROM DEPOSIT_RECORD where BUY_ADDR=$1", address)
	checkErr(err)
	res := make([]db.DepositRecord, 0, 8)
	defer rows.Close()
	for rows.Next() {
		var id, mappingId, depositAmount, buyAmount, height uint64
		var seq int64
		var buyAddr, refAddr, superiorRefAddr, transactionId, rate string
		var creation, updatedAt time.Time
		err = rows.Scan(&id, &creation, &mappingId, &buyAddr, &refAddr, &superiorRefAddr, &seq, &updatedAt, &transactionId, &depositAmount, &buyAmount, &rate, &height)
		checkErr(err)
		res = append(res, db.DepositRecord{Id: id,
			Creation:        creation,
			MappingId:       mappingId,
			BuyAddr:         buyAddr,
			RefAddr:         refAddr,
			SuperiorRefAddr: superiorRefAddr,
			Seq:             seq,
			UpdatedAt:       updatedAt.Unix(),
			TransactionId:   transactionId,
			DepositAmount:   depositAmount,
			BuyAmount:       buyAmount,
			Rate:            rate,
			Height:          height})
	}
	return res
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
