package service

import (
	"github.com/spaco/affiliate/src/service/db"
	pg "github.com/spaco/affiliate/src/service/postgresql"
	client "github.com/spaco/affiliate/src/teller_client"
)

func AllCryptocurrencyMap() map[string]db.CryptocurrencyInfo {
	slice := AllCryptocurrency()
	m := make(map[string]db.CryptocurrencyInfo, 16)
	for _, info := range slice {
		m[info.ShortName] = info
	}
	return m
}

func AllCryptocurrency() []db.CryptocurrencyInfo {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	all := pg.AllCryptocurrency(tx)
	checkErr(tx.Commit())
	commit = true
	return all
}

func AddBatchCryptocurrency(batch []db.CryptocurrencyInfo) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	pg.AddBatchCryptocurrency(tx, batch)
	checkErr(tx.Commit())
	commit = true
}

func MappingDepositAddr(address string, currencyType string, ref string) (string, error) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	buyAddrMapping, found := pg.QueryMappingDepositAddr(tx, address, currencyType)
	if found {
		checkErr(tx.Commit())
		commit = true
		return buyAddrMapping.DepositAddr, nil
	}
	depositAddr, err := client.Bind(currencyType, address)
	if err != nil {
		return "", err
	}
	pg.SaveDepositAddrMapping(tx, address, currencyType, ref, depositAddr)
	checkErr(tx.Commit())
	commit = true
	return depositAddr, nil
}

func CheckMappingAddr(address string, currencyType string) bool {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	_, found := pg.QueryMappingDepositAddr(tx, address, currencyType)
	checkErr(tx.Commit())
	commit = true
	return found
}

func CheckCryptocurrency(shortName string) bool {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	_, found := pg.GetRate(tx, shortName)
	checkErr(tx.Commit())
	commit = true
	return found
}

func QueryDepositRecord(address string) []db.DepositRecord {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	res := pg.QueryDepositRecord(tx, address)
	checkErr(tx.Commit())
	commit = true
	return res
}
