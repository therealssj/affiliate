package postgresql

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const emptyStr = ""

func open() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=lijt password=lijtlijt dbname=affiliate sslmode=disable")
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetTrackingCode(address string) (uint64, string) {
	db := open()
	defer db.Close()
	rows, err := db.Query("SELECT ID,REF_ADDRESS FROM TRACKING_CODE where ADDRESS=$1", address)
	checkErr(err)
	for rows.Next() {
		var id uint64
		var refAddress sql.NullString
		err = rows.Scan(&id, &refAddress)
		checkErr(err)
		if refAddress.Valid {
			return id, refAddress.String
		} else {
			return id, emptyStr
		}
	}
	return 0, emptyStr
}

func GetAddrById(id uint64) (string, string) {
	db := open()
	defer db.Close()
	rows, err := db.Query("SELECT ADDRESS,REF_ADDRESS FROM TRACKING_CODE where ID=$1", id)
	checkErr(err)
	for rows.Next() {
		var address string
		var refAddress sql.NullString
		err = rows.Scan(&address, &refAddress)
		checkErr(err)
		if refAddress.Valid {
			return address, refAddress.String
		} else {
			return address, emptyStr
		}
	}
	return emptyStr, emptyStr
}

func GenerateTrackingCode(address string, refAddress string) uint64 {
	db := open()
	defer db.Close()
	var lastInsertId uint64
	var ra interface{} = nil
	if len(refAddress) > 0 {
		ra = refAddress
	}
	err := db.QueryRow("insert into TRACKING_CODE(ADDRESS,REF_ADDRESS,CREATION) values ($1, $2, now()) returning id;", address, ra).Scan(&lastInsertId)
	checkErr(err)
	return lastInsertId
}
