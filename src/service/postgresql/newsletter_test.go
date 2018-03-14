package postgresql

import (
	"testing"

	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
)

func TestNewsletter(t *testing.T) {
	config := config.GetServerConfig()
	dbo := db.OpenDb(&config.Db)
	defer dbo.Close()
	tx, _ := dbo.Begin()
	defer tx.Rollback()
	email := "admin@test.com"
	SaveNewsletterEmail(tx, email, true)
	exist, concernMiner := ExistNewsletterEmail(tx, email)
	if !exist || !concernMiner {
		t.Errorf("Failed. ")
	}

}
