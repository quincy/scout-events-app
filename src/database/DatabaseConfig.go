package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"log"
	"os"
	"strconv"
)

type PgxClient interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Close()
}

type DbConfig struct {
	Username string
	Password string
	Hostname string
	Port     int
	Dbname   string
}

func CreateDbConnection(ctx context.Context, config *DbConfig) (PgxClient, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.Username, config.Password, config.Hostname, config.Port, config.Dbname)
	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse database connection string: %v", err)
	}
	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	conn, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return conn, err
}

func MakeDatabaseConfig() (*DbConfig, error) {
	dbUsername, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		dbUsername = "admin"
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		dbPassword = "admin"
	}
	dbHostname, ok := os.LookupEnv("DB_HOSTNAME")
	if !ok {
		dbHostname = "localhost"
	}
	var dbPort int
	dbPortStr, ok := os.LookupEnv("DB_PORT")
	if !ok {
		dbPort = 26257
	} else {
		var err error
		dbPort, err = strconv.Atoi(dbPortStr)
		if err != nil {
			return nil, err
		}
	}
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		dbName = "scouting"
	}

	return &DbConfig{
		Username: dbUsername,
		Password: dbPassword,
		Hostname: dbHostname,
		Port:     dbPort,
		Dbname:   dbName,
	}, nil
}
