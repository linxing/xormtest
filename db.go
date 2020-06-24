package xormtest

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type DB struct {
	engine *xorm.Engine
	name   string
}

// Close drop databases and close database connection
func (db *DB) Close() error {
	// close database connection
	defer db.engine.Close()
	// drop database
	return dropDatabase(db.engine.DB().DB, db.name)
}

func (db *DB) initDB(beans []interface{}) error {

	if testing.Verbose() {
		// db.engine.ShowExecTime(true)
		db.engine.ShowSQL(true)
	}

	return db.engine.Sync2(beans...)
}

func (db *DB) Engine() *xorm.Engine {
	return db.engine
}

func NewDB(name string, beans ...interface{}) (*DB, error) {
	suffix := genSuffix()
	dbName := fmt.Sprintf("%s_%s", name, suffix)
	if err := createDatabase(getDSN(""), dbName); err != nil {
		return nil, err
	}

	engine, err := xorm.NewEngine("mysql", getDSN(dbName))
	if err != nil {
		return nil, err
	}

	db := &DB{engine, dbName}
	if err := db.initDB(beans); err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return db, nil
}

func createDatabase(dsn string, name string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	return err
}

func dropDatabase(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE " + name)
	return err
}

func genSuffix() string {
	ts := time.Now().UnixNano()
	return fmt.Sprintf("%d", rand.New(rand.NewSource(ts)).Int())
}

func getDSN(name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		optional(os.Getenv("DB_USERNAME"), "root"),
		optional(os.Getenv("DB_PASSWORD"), "root"),
		optional(os.Getenv("DB_ADDR"), "localhost:3306"),
		name)
}

func optional(value string, option string) string {
	if value != "" {
		return value
	}
	return option
}
