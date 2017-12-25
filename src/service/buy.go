package service

import (
	"github.com/spaco/affiliate/src/service/db"
	pg "github.com/spaco/affiliate/src/service/postgresql"
	client "github.com/spaco/affiliate/src/teller_client"
)

func AllCryptocurrencyMap() map[string]*db.CryptocurrencyInfo {
	slice := AllCryptocurrency()
	m := make(map[string]*db.CryptocurrencyInfo, 16)
	for _, info := range slice {
		m[info.ShortName] = info
	}
	return m
}

func AllCryptocurrency() []*db.CryptocurrencyInfo {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	all := pg.AllCryptocurrency(tx)
	checkErr(tx.Commit())
	commit = true
	return all
}

func AddBatchCryptocurrency(batch ...*db.CryptocurrencyInfo) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	pg.AddBatchCryptocurrency(tx, batch...)
	checkErr(tx.Commit())
	commit = true
}

func MappingDepositAddr(address string, currencyType string, ref string) string {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	buyAddrMapping, found := pg.QueryMappingDepositAddr(tx, address, currencyType)
	if found {
		checkErr(tx.Commit())
		commit = true
		return buyAddrMapping.DepositAddr
	}
	depositAddr := client.Bind(currencyType, address)
	pg.SaveDepositAddrMapping(tx, address, currencyType, ref, depositAddr)
	checkErr(tx.Commit())
	commit = true
	return depositAddr
}

func CheckStatus(address string) {

}
