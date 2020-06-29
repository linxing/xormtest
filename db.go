package xormtest

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"xorm.io/xorm"
)

type DB struct {
	engine *xorm.Engine
	name   string
	tables []interface{}
}

// Close drop databases and close database connection
func (db *DB) Close() error {
	// close database connection
	defer db.engine.Close()
	// drop tables
	return dropTables(db.engine, db.name, db.tables)
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

	if err = createDatabase(driver, dataSourceName, dbName); err != nil {
		return nil, err
	}

	switch driver {
	case "postgres":
		engine, err = xorm.NewEngine("postgres", dataSourceName+" dbname="+dbName)
		if err != nil {
			return nil, err
		}
	default:
		engine, err = xorm.NewEngine("mysql", dataSourceName+dbName)
		if err != nil {
			return nil, err
		}
	}

	db := &DB{
		engine: engine,
		name:   dbName,
		tables: beans,
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

	switch driver {
	case "postgres":
		engine, err := xorm.NewEngine("postgres", dataSourceName+" dbname=postgres")
		if err != nil {
			return err
		}

		defer engine.Close()

		_, err = engine.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			if pqerr, ok := err.(*pq.Error); ok && pqerr.Code != "42P04" {
				fmt.Println(pqerr.Code)
				return err
			}
		}

		return nil

	default:
		engine, err := xorm.NewEngine("mysql", dataSourceName)
		if err != nil {
			return err
		}

		defer engine.Close()

		_, err = engine.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
		return err
	}
}

func dropTables(dbEngine *xorm.Engine, dbName string, tables []interface{}) error {
	return dbEngine.DropTables(tables...)
}
