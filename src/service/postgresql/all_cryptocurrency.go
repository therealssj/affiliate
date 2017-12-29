package postgresql

import (
	"database/sql"

	"github.com/spaco/affiliate/src/service/db"
)

func AllCryptocurrency(tx *sql.Tx) []db.CryptocurrencyInfo {
	rows, err := tx.Query("SELECT SHORT_NAME,FULL_NAME,RATE,UNIT_POWER FROM ALL_CRYPTOCURRENCY order by SHORT_NAME")
	checkErr(err)
	res := make([]db.CryptocurrencyInfo, 0, 10)
	defer rows.Close()
	for rows.Next() {
		var shortName, fullName, rate string
		var unitPower int32
		err = rows.Scan(&shortName, &fullName, &rate, &unitPower)
		checkErr(err)
		res = append(res, db.CryptocurrencyInfo{shortName, fullName, rate, unitPower})
	}
	return res
}

func AddBatchCryptocurrency(tx *sql.Tx, batch []db.CryptocurrencyInfo) {
	stmt, err := tx.Prepare("insert into ALL_CRYPTOCURRENCY(SHORT_NAME,FULL_NAME,RATE,UNIT_POWER) values ($1, $2, $3, $4)")
	defer stmt.Close()
	checkErr(err)
	for _, info := range batch {
		_, err = stmt.Exec(info.ShortName, info.FullName, info.Rate, info.UnitPower)
		checkErr(err)
	}
}

func UpdateBatchRate(tx *sql.Tx, batch []db.CryptocurrencyInfo) {
	stmt, err := tx.Prepare("update ALL_CRYPTOCURRENCY set RATE=$1 where SHORT_NAME=$2")
	defer stmt.Close()
	checkErr(err)
	for _, info := range batch {
		_, err = stmt.Exec(info.Rate, info.ShortName)
		checkErr(err)
	}
}

func GetCryptocurrency(tx *sql.Tx, shortName string) *db.CryptocurrencyInfo {
	rows, err := tx.Query("SELECT SHORT_NAME,FULL_NAME,RATE,UNIT_POWER FROM ALL_CRYPTOCURRENCY where SHORT_NAME=$1", shortName)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var shortName, fullName, rate string
		var unitPower int32
		err = rows.Scan(&shortName, &fullName, &rate, &unitPower)
		checkErr(err)
		return &db.CryptocurrencyInfo{shortName, fullName, rate, unitPower}
	}
	return nil
}
