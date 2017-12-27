package postgresql

import (
	"database/sql"
	"github.com/spaco/affiliate/src/service/db"
)

func UpdateRewardRemain(tx *sql.Tx, data map[string]uint64) {
	stmt, err := tx.Prepare("INSERT INTO REWARD_REMAIN (ADDRESS, CREATION, LAST_MODIFIED, AMOUNT) VALUES ($1,now(),now(),$2) ON CONFLICT (ADDRESS) DO UPDATE SET AMOUNT=$2")
	defer stmt.Close()
	checkErr(err)
	for k, v := range data {
		_, err = stmt.Exec(k, v)
		checkErr(err)
	}
}

func QueryRewardRemain(tx *sql.Tx, addr ...string) map[string]uint64 {
	m := make(map[string]uint64, len(addr))
	args := make([]interface{}, 0, len(addr))
	for _, ad := range addr {
		args = append(args, ad)
	}
	rows, err := tx.Query("select ADDRESS,AMOUNT from REWARD_REMAIN where ADDRESS in "+db.InClause(len(addr), 1), args...)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var address string
		var amount uint64
		err = rows.Scan(&address, amount)
		checkErr(err)
		m[address] = amount
	}
	return m
}
