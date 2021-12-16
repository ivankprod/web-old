package db

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tarantool/go-tarantool"
)

// Connect to Tarantool
func ConnectTarantool() (*tarantool.Connection, error) {
	conn, err := tarantool.Connect(fmt.Sprintf("%s:%s", os.Getenv("DB_TARANTOOL_HOST"), os.Getenv("DB_TARANTOOL_PORT")), tarantool.Opts{
		User: os.Getenv("DB_TARANTOOL_USER"),
		Pass: os.Getenv("DB_TARANTOOL_PASSWORD"),
	})

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Connect to MySQL
func ConnectMySQL() (*sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))
	dbCredentials := fmt.Sprintf("%s:%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	dbConnStr := fmt.Sprintf("tcp(%s:%s)/%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s@%s", dbCredentials, dbConnStr))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, err
	}

	return db, nil
}
