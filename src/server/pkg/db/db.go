package db

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tarantool/go-tarantool"

	BaseLogger "github.com/ivankprod/ivankprod.ru/src/server/pkg/logger"
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

	if resp, err := conn.Ping(); resp == nil || err != nil {
		return nil, err
	}

	return conn, nil
}

func KeepAliveTarantool(conn *tarantool.Connection) *tarantool.Connection {
	var (
		connNew *tarantool.Connection
		mu      sync.Mutex
		err     error
	)

	logger := BaseLogger.Get()

	for {
		time.Sleep(time.Second * 3)

		lost := false

		if conn == nil {
			lost = true
		} else if resp, err := conn.Ping(); resp == nil || err != nil {
			lost = true
		}

		if !lost {
			continue
		}

		logger.Println("Lost Tarantool connection. Restoring...")

		if connNew, err = ConnectTarantool(); connNew == nil || err != nil {
			if err == nil {
				logger.Println("Failed connecting to Tarantool database")
			} else {
				logger.Printf("Error connecting to Tarantool database: %v\n", err)
			}

			continue
		}

		logger.Println("Tarantool connection restored")

		mu.Lock()
		conn = connNew
		mu.Unlock()

		return connNew
	}
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
