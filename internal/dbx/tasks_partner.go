package dbx

import (
	"context"
	"time"
)

func (db *DB) SaveTaskPartnerEvent(ctx context.Context, userID, taskID, subID int64, payout float64, receivedAt time.Time) error {
	const q = `
		INSERT INTO tasks_partner_events (user_id, task_id, sub_id, payout, ts)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT DO NOTHING
	`
	_, err := db.c.Exec(ctx, q, userID, taskID, subID, payout, receivedAt)
	return err
}
