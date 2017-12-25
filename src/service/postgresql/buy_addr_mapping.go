package postgresql

import (
	"database/sql"
	"github.com/spaco/affiliate/src/service/db"
	"time"
)

func QueryMappingDepositAddr(tx *sql.Tx, address string, currencyType string) (*db.BuyAddrMapping, bool) {
	rows, err := tx.Query("select ID,VERSION,CREATION,LAST_MODIFIED,DEPOSIT_ADDR,REF,DEPOSIT_AMOUNT,BUY_AMOUNT,LAST_UPDATED,TRANSACTION_IDS,SENT_COIN from BUY_ADDR_MAPPING where CURRENCY_TYPE=$1 and DEPOSIT_ADDR=$2", currencyType, address)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var id uint64
		var version uint64
		var creation time.Time
		var lastModified time.Time
		var depositAddr string
		var ref string
		var refNullable sql.NullString
		var depositAmount float64
		var depositAmountNullable sql.NullFloat64
		var buyAmount uint64
		var buyAmountNullable sql.NullInt64
		var lastUpdated time.Time
		var transactionIds string
		var transactionIdsNullable sql.NullString
		var sentCoin bool
		err = rows.Scan(&id, &version, &creation, &lastModified, &depositAddr, &refNullable, &depositAmountNullable, &buyAmountNullable, &lastUpdated, &transactionIdsNullable, &sentCoin)
		checkErr(err)
		if refNullable.Valid {
			ref = refNullable.String
		}
		if depositAmountNullable.Valid {
			depositAmount = depositAmountNullable.Float64
		}
		if buyAmountNullable.Valid {
			buyAmount = uint64(buyAmountNullable.Int64)
		}
		if transactionIdsNullable.Valid {
			transactionIds = transactionIdsNullable.String
		}
		return &db.BuyAddrMapping{id, version, creation, lastModified, address, currencyType, depositAddr, ref, depositAmount, buyAmount, lastUpdated, transactionIds, sentCoin}, true
	}
	return &db.BuyAddrMapping{}, false
}

func SaveDepositAddrMapping(tx *sql.Tx, address string, currencyType string, ref string, depositAddr string) uint64 {
	var lastInsertId uint64
	err := tx.QueryRow("insert into BUY_ADDR_MAPPING(VERSION,CREATION,LAST_MODIFIED,ADDRESS,CURRENCY_TYPE,DEPOSIT_ADDR,REF,SENT_COIN) values (0,now(),now(),$1, $2, $3,$4,false) returning id;", address, currencyType, depositAddr, ref).Scan(&lastInsertId)
	checkErr(err)
	return lastInsertId

}
