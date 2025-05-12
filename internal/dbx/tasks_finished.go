package dbx

import (
	"context"
	"time"
)

func (db *DB) SaveTaskFinish(ctx context.Context, userID, taskID, subID, taskPoints int64, finishAt time.Time) error {
	const q = `INSERT INTO tasks_finished (user_id, task_id, sub_id, points, ts) VALUES ($1,$2,$3,$4,$5)`
	_, err := db.c.Exec(ctx, q, userID, taskID, subID, taskPoints, finishAt)
	return err
}
