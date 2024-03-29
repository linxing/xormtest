package xormtest

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

type DB struct {
	engine *xorm.Engine
	name   string
}

// Close drop database and close database connection
func (db *DB) Close() (err error) {
	// drop database
	defer func() {
		if e := dropDatabase(db.engine.DriverName(), db.engine.DataSourceName(), db.name); e != nil {
			err = e
		}
	}()
	// close database connection
	return db.engine.Close()
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

func NewDB(driver string, dataSourceName string, dbName string, beans ...interface{}) (*DB, error) {

	var engine *xorm.Engine
	var err error

	dbName = getDBNameWithPrefix(dbName)
	if err := createDatabase(driver, dataSourceName, dbName); err != nil {
		return nil, err
	}

	switch driver {
	case "postgres":
		engine, err = xorm.NewEngine(driver, dataSourceName+" dbname="+dbName)
	case "sqlite3":
		engine, err = xorm.NewEngine(driver, dataSourceName)
	default:
		engine, err = xorm.NewEngine(driver, dataSourceName+dbName)
	}
	if err != nil {
		return nil, err
	}

	db := &DB{
		engine: engine,
		name:   dbName,
	}

	if err := db.initDB(beans); err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return db, nil
}

func NewCustomizePrefixDB(driver string, dataSourceName string, dbName string, databasePrefix string, beans ...interface{}) (*DB, error) {

	var engine *xorm.Engine
	var err error

	dbName = databasePrefix + dbName

	if err := createDatabase(driver, dataSourceName, dbName); err != nil {
		return nil, err
	}

	switch driver {
	case "postgres":
		engine, err = xorm.NewEngine(driver, dataSourceName+" dbname="+dbName)
	case "sqlite3":
		engine, err = xorm.NewEngine(driver, dataSourceName)
	default:
		engine, err = xorm.NewEngine(driver, dataSourceName+dbName)
	}
	if err != nil {
		return nil, err
	}

	db := &DB{
		engine: engine,
		name:   dbName,
	}

	if err := db.initDB(beans); err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return db, nil
}

func createDatabase(driver string, dataSourceName string, dbName string) error {

	var err error
	var db *sql.DB

	if driver == "postgres" {
		db, err = sql.Open(driver, dataSourceName+" dbname=postgres")
	} else {
		db, err = sql.Open(driver, dataSourceName)
	}
	if err != nil {
		return err
	}

	defer db.Close()

	if driver == "postgres" {
		if _, err = db.Exec("CREATE DATABASE " + dbName); err != nil {
			if pqerr, ok := err.(*pq.Error); ok && pqerr.Code == "42P04" {
				return nil
			}
		}
	} else if driver != "sqlite3" {
		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
		if err != nil {
			return err
		}
	}

	return nil
}

func dropDatabase(driver string, dataSourceName string, dbName string) error {

	var err error
	var db *sql.DB

	if driver == "postgres" {
		db, err = sql.Open(driver, dataSourceName+" dbname=postgres")
	} else {
		db, err = sql.Open(driver, dataSourceName)
	}

	if err != nil {
		return err
	}

	defer db.Close()

	if driver == "sqlite3" {
		err = os.Remove(dataSourceName)
	} else {
		_, err = db.Exec("DROP DATABASE " + dbName)
	}

	return err
}

func getDBNameWithPrefix(dbName string) string {
	ts := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d", dbName, ts)
}
