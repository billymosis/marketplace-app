package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func Connection(driver, host, database, username, password string, port, maxOpenConnections int) (*pgxpool.Pool, error) {
	dsn, err := parseDSN(driver, host, database, username, password, port, maxOpenConnections)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		//debug.PrintStack()
		return nil, err
	}

	if err := pingDatabase(db); err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), db)

	return pool, err
}

func pingDatabase(db *pgxpool.Config) error {
	pool, err := pgxpool.NewWithConfig(context.Background(), db)
	if err != nil {
		log.Fatal(err)
		pool.Close()
		return errPingDatabase
	}
	r := 3
	for i := 0; i < r; i++ {
		err := pool.Ping(context.Background())
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	pool.Close()
	return errPingDatabase
}

func parseDSN(driver, host, database, username, password string, port int, maxconn int) (string, error) {

	switch driver {
	case "postgres":
		return postgreParseDSN(host, database, username, password, port, maxconn), nil
	default:
		return "", errUnSupportedDriver
	}
}

func postgreParseDSN(host, database, username, password string, port int, maxconn int) string {
	if os.Getenv("ENV") == "production" {
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%s sslmode=verify-full sslrootcert=ap-southeast-1-bundle.pem TimeZone=UTC",
			host, port, username, password, database, strconv.Itoa(maxconn))
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%s sslmode=disable",
		host, port, username, password, database, strconv.Itoa(maxconn))
}
