package postgresql

import (
	"database/sql"

	"github.com/spolabs/affiliate/src/service/db"
)

func AllCryptocurrency(tx *sql.Tx) []db.CryptocurrencyInfo {
	rows, err := tx.Query("SELECT SHORT_NAME,FULL_NAME,RATE,UNIT_POWER,ENABLED FROM ALL_CRYPTOCURRENCY order by SHORT_NAME")
	checkErr(err)
	res := make([]db.CryptocurrencyInfo, 0, 10)
	defer rows.Close()
	for rows.Next() {
		var shortName, fullName, rate string
		var unitPower int32
		var enabled bool
		err = rows.Scan(&shortName, &fullName, &rate, &unitPower, &enabled)
		checkErr(err)
		res = append(res, db.CryptocurrencyInfo{shortName, fullName, rate, unitPower, enabled})
	}
	return res
}

func AddBatchCryptocurrency(tx *sql.Tx, batch []db.CryptocurrencyInfo) {
	stmt, err := tx.Prepare("insert into ALL_CRYPTOCURRENCY(SHORT_NAME,FULL_NAME,RATE,UNIT_POWER,ENABLED) values ($1, $2, $3, $4,$5)")
	defer stmt.Close()
	checkErr(err)
	for _, info := range batch {
		_, err = stmt.Exec(info.ShortName, info.FullName, info.Rate, info.UnitPower, info.Enabled)
		checkErr(err)
	}
}

func UpdateBatchRateAndEnabled(tx *sql.Tx, batch []db.CryptocurrencyInfo) {
	stmt, err := tx.Prepare("update ALL_CRYPTOCURRENCY set RATE=$1,ENABLED=$2 where SHORT_NAME=$3")
	defer stmt.Close()
	checkErr(err)
	for _, info := range batch {
		_, err = stmt.Exec(info.Rate, info.Enabled, info.ShortName)
		checkErr(err)
	}
}

func GetCryptocurrency(tx *sql.Tx, shortName string) *db.CryptocurrencyInfo {
	rows, err := tx.Query("SELECT SHORT_NAME,FULL_NAME,RATE,UNIT_POWER,ENABLED FROM ALL_CRYPTOCURRENCY where SHORT_NAME=$1", shortName)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var shortName, fullName, rate string
		var unitPower int32
		var enabled bool
		err = rows.Scan(&shortName, &fullName, &rate, &unitPower, &enabled)
		checkErr(err)
		return &db.CryptocurrencyInfo{shortName, fullName, rate, unitPower, enabled}
	}
	return nil
}
