package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quincy/scout-events-app/src/config"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"log"
)

type PgxClient interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Close()
}

func CreateDbConnection(ctx context.Context, config config.AppConfig) (PgxClient, error) {
	connString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		config.DatabaseUsername(),
		config.DatabasePassword(),
		config.DatabaseHostname(),
		config.DatabasePort(),
		config.DatabaseName(),
	)
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
