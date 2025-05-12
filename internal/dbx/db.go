package dbx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	p *pgxpool.Pool
	c dbx
}

var ErrNoRows = pgx.ErrNoRows

func isErrUniqueConstraint(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unique constraint")
}

type dbx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func Connect(ctx context.Context, cfg Config) (*DB, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s", cfg.User, cfg.Pass, cfg.Addr, cfg.Name)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return &DB{p: pool,
		c: pool,
	}, nil
}

func (db *DB) Close() {
	db.p.Close()
}

func (db *DB) Tx(ctx context.Context, callback func(tx *DB) error) error {
	if db.p == nil {
		return callback(db)
	}
	tx, err := db.p.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	db = dbPool.Get().(*DB)
	db.c = tx
	defer func() {
		db.c = nil
		dbPool.Put(db)
	}()
	if err = callback(db); err == nil {
		return tx.Commit(ctx)
	}
	return errors.Join(err, tx.Rollback(ctx))
}

var dbPool = sync.Pool{New: func() any { return new(DB) }}
