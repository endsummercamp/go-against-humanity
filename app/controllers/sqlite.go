package controllers

import (
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

type User struct {
	Id int64 `db:"user_id"`
	Username string
	PwHash string
}

func InitDB() {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("sqlite3", path.Join(workdir, "database.sqlite3"))
	if err != nil {
		panic(err)
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
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

/*

func (c *GorpController) Begin() revel.Result {
	panic("e")
	txn, err := DbMap.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

*/

func init() {
	revel.OnAppStart(InitDB)
	/*
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	// revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	// revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
	*/
}