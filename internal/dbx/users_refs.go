package dbx

import "context"

func (db *DB) SaveUserRefferal(ctx context.Context, fromUserID, toUserID, points int64, level int) error {
	const q = `
		INSERT INTO users_refs (from_id, to_id, points, level)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (from_id, to_id)
		DO UPDATE SET points = users_refs.points + $3
	`
	_, err := db.c.Exec(ctx, q, fromUserID, toUserID, points, level)
	return err
}
