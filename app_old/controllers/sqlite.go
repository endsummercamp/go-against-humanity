package controllers

import (
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"

	"github.com/revel/revel"
	"os"
	"path"
)

var (
	DbMap *gorp.DbMap
)

func InitDB() {
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
	DbMap = dbmap
}

type GorpController struct {
	*revel.Controller
	Txn *gorp.Transaction
}

func init() {
	revel.OnAppStart(InitDB)
}