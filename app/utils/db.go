package utils

import (
	"database/sql"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

func InitDB() *gorp.DbMap {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("sqlite3", path.Join(workdir, "database.sqlite3"))
	if err != nil {
		panic(err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(models.User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(models.Match{}, "matches").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}

	return dbmap
}
