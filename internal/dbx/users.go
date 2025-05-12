package dbx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avpetkun/the-prime/pkg/tgu"
)

func (db *DB) CreateUser(ctx context.Context, joinAt time.Time, u tgu.User, refID int64, ipAddress, userAgent string) error {
	const q = `
		INSERT INTO users (
			joined,
			id, ref_id, ip_address, user_agent,
			first_name, last_name, username,
			lang_code, is_premium, photo_url, allow_messages
		)
		VALUES ($1, $2,$3,$4,$5, $6,$7,$8, $9,$10,$11,$12)
		ON CONFLICT (id)
		DO UPDATE SET
			ip_address = EXCLUDED.ip_address,
			user_agent = EXCLUDED.user_agent,
			first_name = EXCLUDED.first_name,
			last_name  = EXCLUDED.last_name,
			username   = EXCLUDED.username,
			lang_code  = EXCLUDED.lang_code,
			is_premium = EXCLUDED.is_premium
		WHERE
			users.id = EXCLUDED.id
	`
	_, err := db.c.Exec(
		ctx, q,
		joinAt,
		u.ID, refID, ipAddress, userAgent,
		u.FirstName, u.LastName, u.Username,
		u.LanguageCode, u.IsPremium, u.PhotoUrl, u.AllowsWriteToPm,
	)
	return err
}

func (db *DB) SetUserInit(ctx context.Context, userID int64) error {
	const q = `UPDATE users SET inited = NOW() WHERE id = $1 AND inited is null`
	_, err := db.c.Exec(ctx, q, userID)
	return err
}

func (db *DB) SetUserUsername(ctx context.Context, userID int64, newUsername string) error {
	const q = `UPDATE users SET username = $1 WHERE id = $2`
	_, err := db.c.Exec(ctx, q, newUsername, userID)
	return err
}

func (db *DB) GetAllUsersPoints(ctx context.Context, walk func(userID, points int64) error) error {
	const q = `SELECT id, points FROM users`
	rows, err := db.c.Query(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var userID, points int64
		if err = rows.Scan(&userID, &points); err != nil {
			return err
		}
		if err = walk(userID, points); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) GetUserPoints(ctx context.Context, userID int64) (points int64, err error) {
	const q = `SELECT points FROM users WHERE id = $1`
	err = db.c.QueryRow(ctx, q, userID).Scan(&points)
	return
}

func (db *DB) GetUserRefs(ctx context.Context, userID int64) (refPoints, refCount int64, err error) {
	const q = `SELECT ref_points, ref_count FROM users WHERE id = $1`
	err = db.c.QueryRow(ctx, q, userID).Scan(&refPoints, &refCount)
	return
}

func (db *DB) IncUserPoints(ctx context.Context, userID, incPoints int64) error {
	const q = `UPDATE users SET points = points + $1 WHERE id = $2`
	_, err := db.c.Exec(ctx, q, incPoints, userID)
	return err
}

func (db *DB) IncUserPointsFromRef(ctx context.Context, userID, incPoints int64) error {
	const q = `UPDATE users SET points = points + $1, ref_points = ref_points + $1 WHERE id = $2`
	_, err := db.c.Exec(ctx, q, incPoints, userID)
	return err
}

func (db *DB) IncUserRefCount(ctx context.Context, userID int64) error {
	const q = `UPDATE users SET ref_count = ref_count + 1 WHERE id = $1`
	_, err := db.c.Exec(ctx, q, userID)
	return err
}

func (db *DB) SetUserPoints(ctx context.Context, userID, newPoints int64) error {
	const q = `UPDATE users SET points = $1 WHERE id = $2`
	_, err := db.c.Exec(ctx, q, newPoints, userID)
	return err
}

func (db *DB) WithdrawUserPoints(ctx context.Context, userID, decPoints int64) error {
	const q = `UPDATE users SET points = points - $2 WHERE id = $1 AND points >= $2`
	cmd, err := db.c.Exec(ctx, q, userID, decPoints)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("insufficient balance")
	}
	return nil
}

func (db *DB) GetUserRefID(ctx context.Context, userID int64) (refID int64, err error) {
	const q = `SELECT ref_id FROM users WHERE id = $1`
	err = db.c.QueryRow(ctx, q, userID).Scan(&refID)
	if err != nil && errors.Is(err, ErrNoRows) {
		err = nil
	}
	return
}
