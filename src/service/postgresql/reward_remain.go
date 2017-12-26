package postgresql

import (
	"database/sql"
)

func UpdateRewardRemain(tx *sql.Tx, address string, amount uint64, version uint64) {

}

func QueryRewardRemain(tx *sql.Tx, address string) (amount uint64, version uint64) {

	return 0, 0
}
