package xormtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type User struct {
	ID          int64     `xorm:"'id' not null BIGINT(22) pk autoincr"`
	UserName    string    `xorm:"'user_name' not null VARCHAR(255) unique(unique_user_name)"`
	Password    string    `xorm:"'password' not null VARCHAR(255)"`
	PhoneNumber string    `xorm:"'phone_number' not null VARCHAR(255) unique(unique_phone_number)"`
	CreatedAt   time.Time `xorm:"'created_at' not null created DATETIME"`
	UpdatedAt   time.Time `xorm:"'updated_at' not null updated DATETIME"`
}

var (
	beans = []interface{}{
		new(User),
	}
)

func TestNewMySQLDB(t *testing.T) {
	mockdb, err := NewDB("mysql", "root:root@tcp(localhost:3306)/", "xormtest", beans...)
	require.NoError(t, err)

	defer mockdb.Close()

	tables, err := mockdb.engine.DBMetas()
	require.NoError(t, err)

	assert.Len(t, tables, len(beans))
}

func TestNewPostgresDB(t *testing.T) {
	mockdb, err := NewDB("postgres", "user=postgres password=postgres host=localhost port=5432 sslmode=disable", "xormtest", beans...)
	require.NoError(t, err)

	defer mockdb.Close()

	tables, err := mockdb.engine.DBMetas()
	require.NoError(t, err)

	assert.Len(t, tables, len(beans))
}

/*func TestNewSqlite3DB(t *testing.T) {
	mockdb, err := NewDB("sqlite3", "test.db", "xormtest", beans...)
	require.NoError(t, err)

	defer mockdb.Close()

	tables, err := mockdb.engine.DBMetas()
	require.NoError(t, err)

	assert.Len(t, tables, len(beans))
}*/
