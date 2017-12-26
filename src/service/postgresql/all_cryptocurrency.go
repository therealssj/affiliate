package postgresql

import (
	"database/sql"

	"github.com/spaco/affiliate/src/service/db"
)

func AllCryptocurrency(tx *sql.Tx) []*db.CryptocurrencyInfo {
	rows, err := tx.Query("SELECT SHORT_NAME,FULL_NAME,RATE FROM ALL_CRYPTOCURRENCY")
	checkErr(err)
	res := make([]*db.CryptocurrencyInfo, 0, 10)
	defer rows.Close()
	for rows.Next() {
		var shortName, fullName, rate string
		err = rows.Scan(&shortName, &fullName, &rate)
		checkErr(err)
		res = append(res, &db.CryptocurrencyInfo{shortName, fullName, rate})
	}
	return res
}

func AddBatchCryptocurrency(tx *sql.Tx, batch ...*db.CryptocurrencyInfo) {
	stmt, err := tx.Prepare("insert into ALL_CRYPTOCURRENCY(SHORT_NAME,FULL_NAME,RATE) values ($1, $2, $3)")
	defer stmt.Close()
	checkErr(err)
	for _, info := range batch {
		_, err = stmt.Exec(info.ShortName, info.FullName, info.Rate)
		checkErr(err)
	}
}

func UpdateBatchRate(tx *sql.Tx, batch ...*db.CryptocurrencyInfo) {
	stmt, err := tx.Prepare("update ALL_CRYPTOCURRENCY set RATE=$1 where SHORT_NAME=$2")
	defer stmt.Close()
	checkErr(err)
	for _, info := range batch {
		_, err = stmt.Exec(info.Rate, info.ShortName)
		checkErr(err)
	}
}

func GetRate(tx *sql.Tx, shortName string) (string, bool) {
	rows, err := tx.Query("SELECT RATE FROM ALL_CRYPTOCURRENCY where SHORT_NAME=$1", shortName)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var rate string
		err = rows.Scan(&rate)
		checkErr(err)
		return rate, true
	}
	return "", false
}
