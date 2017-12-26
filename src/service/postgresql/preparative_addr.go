package postgresql

import (
	"database/sql"
)

func getPreparativeAddr(tx *sql.Tx, currencyType string, size uint16) []string {
	rows, err := tx.Query("SELECT DEPOSIT_ADDR FROM PREPARATIVE_ADDR where CURRENCY_TYPE=$1 limit $2", currencyType, size)
	checkErr(err)
	res := make([]string, 0, size)
	defer rows.Close()
	for rows.Next() {
		var str string
		err = rows.Scan(&str)
		checkErr(err)
		res = append(res, str)
	}
	return res
}

func SaveBatchPreparativeAddr(tx *sql.Tx, currencyType string, depositAddr []string) {
	stmt, err := tx.Prepare("insert into PREPARATIVE_ADDR(CURRENCY_TYPE,DEPOSIT_ADDR) values ($1, $2)")
	checkErr(err)
	defer stmt.Close()
	for _, addr := range depositAddr {
		_, err = stmt.Exec(currencyType, addr)
		checkErr(err)
	}
}

func RemovePreparativeAddr(tx *sql.Tx, currencyType string, depositAddr string) {
	stmt, err := tx.Prepare("delete from PREPARATIVE_ADDR where CURRENCY_TYPE=$1 and  DEPOSIT_ADDR=$2")
	checkErr(err)
	defer stmt.Close()
	_, err = stmt.Exec(currencyType, depositAddr)
	checkErr(err)
}
