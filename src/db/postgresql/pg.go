package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spaco/affiliate/src/config"
)

var db *sql.DB

const emptyStr = ""

func OpenDb(dbConfig *config.Db) {
	conn_str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SslMode)
	//	fmt.Println(conn_str)
	var err error
	db, err = sql.Open("postgres", conn_str)
	checkErr(err)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(50)
	db.Ping()
}

func CloseDb() {
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetTrackingCodeOrGenerate(address string, refAddress string) uint64 {
	id, _ := GetTrackingCode(address)
	if id == 0 {
		id = GenerateTrackingCode(address, refAddress)
	}
	return id
}

func GetTrackingCode(address string) (uint64, string) {
	rows, err := db.Query("SELECT ID,REF_ADDRESS FROM TRACKING_CODE where ADDRESS=$1", address)
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
			return id, emptyStr
		}
	}
	return 0, emptyStr
}

func GetAddrById(id uint64) (string, string) {
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
	var lastInsertId uint64
	var ra interface{} = nil
	if len(refAddress) > 0 {
		ra = refAddress
	}
	tx, err := db.Begin()
	checkErr(err)
	defer tx.Rollback()
	err = tx.QueryRow("insert into TRACKING_CODE(ADDRESS,REF_ADDRESS,CREATION) values ($1, $2, now()) returning id;", address, ra).Scan(&lastInsertId)
	checkErr(err)
	checkErr(tx.Commit())
	return lastInsertId
}
