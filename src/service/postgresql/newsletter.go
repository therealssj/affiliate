package postgresql

import (
	"database/sql"
)

func ExistNewsletterEmail(tx *sql.Tx, email string) (bool, bool) {
	rows, err := tx.Query("SELECT CONCERN_MINER FROM NEWSLETTER where EMAIL=$1", email)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var concernMiner bool
		err = rows.Scan(&concernMiner)
		checkErr(err)
		return true, concernMiner
	}
	return false, false
}

func SaveNewsletterEmail(tx *sql.Tx, email string, concernMiner bool) {
	stmt, err := tx.Prepare("insert into NEWSLETTER(EMAIL,CONCERN_MINER,CREATION) values ($1, $2,now())")
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec(email, concernMiner)
	checkErr(err)
}
