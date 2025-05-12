package dbx

import (
	"context"
	"time"
)

func (db *DB) SaveTaskStateClaim(
	ctx context.Context, userID, taskID, subID, taskPoints int64, claimAt time.Time,
) (ok bool, err error) {
	const q = `
		INSERT INTO tasks_state (user_id, task_id, sub_id, finish_count, total_points, updated, status)
		VALUES ($1,$2,$3,1,$4,$5,'claim')
		ON CONFLICT (user_id, task_id, sub_id)
		DO UPDATE SET
			finish_count = tasks_state.finish_count + 1,
			total_points = tasks_state.total_points + $4,
			updated = $5,
			status = 'claim'
		WHERE tasks_state.updated <> $5
	`
	cmd, err := db.c.Exec(ctx, q, userID, taskID, subID, taskPoints, claimAt)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}

func (db *DB) SaveTaskStatePending(ctx context.Context, userID, taskID, subID int64, startAt time.Time) error {
	const q = `
		INSERT INTO tasks_state (user_id, task_id, sub_id, updated, status)
		VALUES ($1,$2,$3,$4,'pending')
		ON CONFLICT (user_id, task_id, sub_id)
		DO UPDATE SET updated = $4, status = 'pending'
		WHERE tasks_state.status <> 'claim'
	`
	_, err := db.c.Exec(ctx, q, userID, taskID, subID, startAt)
	return err
}

func (db *DB) SaveTaskStateDone(ctx context.Context, userID, taskID, subID int64, doneAt time.Time) error {
	const q = `
		INSERT INTO tasks_state (user_id, task_id, sub_id, updated, status)
		VALUES ($1,$2,$3,$4,'done')
		ON CONFLICT (user_id, task_id, sub_id)
		DO UPDATE SET updated = $4, status = 'done'
		WHERE tasks_state.status = 'claim'
	`
	_, err := db.c.Exec(ctx, q, userID, taskID, subID, doneAt)
	return err
}
