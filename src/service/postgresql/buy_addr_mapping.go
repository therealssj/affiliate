package postgresql

import (
	"database/sql"
	"time"

	"github.com/spaco/affiliate/src/service/db"
)

func QueryMappingDepositAddr(tx *sql.Tx, address string, currencyType string) *db.BuyAddrMapping {
	rows, err := tx.Query("select ID,CREATION,DEPOSIT_ADDR,REF from BUY_ADDR_MAPPING where CURRENCY_TYPE=$1 and ADDRESS=$2", currencyType, address)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var id uint64
		var creation time.Time
		var depositAddr string
		var ref string
		var refNullable sql.NullString
		err = rows.Scan(&id, &creation, &depositAddr, &refNullable)
		checkErr(err)
		if refNullable.Valid {
			ref = refNullable.String
		}
		return &db.BuyAddrMapping{id, creation, address, currencyType, depositAddr, ref}
	}
	return nil
}

func SaveDepositAddrMapping(tx *sql.Tx, address string, currencyType string, ref string, depositAddr string) uint64 {
	var lastInsertId uint64
	err := tx.QueryRow("insert into BUY_ADDR_MAPPING(CREATION,ADDRESS,CURRENCY_TYPE,DEPOSIT_ADDR,REF) values (now(),$1, $2, $3, $4) returning ID", address, currencyType, depositAddr, ref).Scan(&lastInsertId)
	checkErr(err)
	return lastInsertId

}
