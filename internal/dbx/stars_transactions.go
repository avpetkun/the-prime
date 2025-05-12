package dbx

import "context"

func (db *DB) SaveStarsTransaction(ctx context.Context, txID string, userID, taskID int64, starsAmount int) error {
	const q = `
		INSERT INTO stars_transactions (
			tx_id, user_id, task_id, amount
		) VALUES ($1,$2,$3,$4,$5)
	 	ON CONFLICT DO NOTHING
	`
	_, err := db.c.Exec(ctx, q, txID, userID, taskID, starsAmount)
	return err
}
