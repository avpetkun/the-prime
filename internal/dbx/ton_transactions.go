package dbx

import (
	"context"

	"github.com/avpetkun/the-prime/pkg/tonu"
)

func (db *DB) SaveTonTx(ctx context.Context, userID, taskID int64, tx *tonu.Transaction) error {
	const q = `INSERT INTO ton_transactions (
			user_id, task_id, ts, hash,
			ton_amount, is_deposit,
			src_addr, dst_addr, comment
		)
		VALUES ($1,$2,$3,$4, $5,$6, $7,$8,$9)
		ON CONFLICT DO NOTHING
	`
	_, err := db.c.Exec(
		ctx, q,
		userID, taskID, tx.Timestamp, tx.Hash,
		tx.TonAmount.Int64(), tx.IsDeposit,
		tx.SrcAddr, tx.DstAddr, tx.Comment,
	)
	return err
}

func (db *DB) CheckTonTxByHash(ctx context.Context, txHash string) (exist bool, err error) {
	const q = `SELECT EXISTS(
		SELECT 1 FROM ton_transactions WHERE "hash" = $1
	)`
	err = db.c.QueryRow(ctx, q, txHash).Scan(&exist)
	return
}
