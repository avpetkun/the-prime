package dbx

import (
	"context"
	"errors"
)

func (db *DB) GetKeyNum(ctx context.Context, key string) (n int, err error) {
	const q = `SELECT num FROM key_val WHERE key = $1`
	err = db.c.QueryRow(ctx, q, key).Scan(&n)
	if errors.Is(err, ErrNoRows) {
		return 0, nil
	}
	return
}

func (db *DB) GetKeyStr(ctx context.Context, key string) (s string, err error) {
	const q = `SELECT str FROM key_val WHERE key = $1`
	err = db.c.QueryRow(ctx, q, key).Scan(&s)
	if errors.Is(err, ErrNoRows) {
		return "", nil
	}
	return
}

func (db *DB) SetKeyNum(ctx context.Context, key string, newNum int) error {
	const q = `
		INSERT INTO key_val (key, num)
		VALUES ($1,$2)
		ON CONFLICT (key)
		DO UPDATE SET
			num = EXCLUDED.num
		WHERE
			key_val.key = EXCLUDED.key
	`
	_, err := db.c.Exec(ctx, q, key, newNum)
	return err
}

func (db *DB) SetKeyStr(ctx context.Context, key, newStr string) error {
	const q = `
		INSERT INTO key_val (key, str)
		VALUES ($1,$2)
		ON CONFLICT (key)
		DO UPDATE SET
			str = EXCLUDED.str
		WHERE
			key_val.key = EXCLUDED.key
	`
	_, err := db.c.Exec(ctx, q, key, newStr)
	return err
}
