package postgresql

import (
	"database/sql"
)

func GetTrackingCode(tx *sql.Tx, address string) (uint64, string) {
	rows, err := tx.Query("SELECT ID,REF_ADDRESS FROM TRACKING_CODE where ADDRESS=$1", address)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var id uint64
		var refAddress sql.NullString
		err = rows.Scan(&id, &refAddress)
		checkErr(err)
		if refAddress.Valid {
			return id, refAddress.String
		} else {
			return id, ""
		}
	}
	return 0, ""
}

func GetAddrById(tx *sql.Tx, id uint64) (string, string) {
	rows, err := tx.Query("SELECT ADDRESS,REF_ADDRESS FROM TRACKING_CODE where ID=$1", id)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var address string
		var refAddress sql.NullString
		err = rows.Scan(&address, &refAddress)
		checkErr(err)
		if refAddress.Valid {
			return address, refAddress.String
		} else {
			return address, ""
		}
	}
	return "", ""
}

func GenerateTrackingCode(tx *sql.Tx, address string, refAddress string) uint64 {
	var lastInsertId uint64
	var ra interface{} = nil
	if len(refAddress) > 0 {
		ra = refAddress
	}
	err := tx.QueryRow("insert into TRACKING_CODE(ADDRESS,REF_ADDRESS,CREATION) values ($1, $2, now()) returning id;", address, ra).Scan(&lastInsertId)
	checkErr(err)
	return lastInsertId
}
