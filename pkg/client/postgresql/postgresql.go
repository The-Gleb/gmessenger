package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/The-Gleb/gmessenger/internal/config"
	"github.com/The-Gleb/gmessenger/pkg/utils/retry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc config.Database) (conn *pgx.Conn, err error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.DbName)
	err = retry.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err = pgx.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return conn, nil
}
