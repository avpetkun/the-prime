package dbx

import (
	"context"

	"github.com/avpetkun/the-prime/internal/common"
)

func (db *DB) GetAllTasks(ctx context.Context) (tasks []*common.FullTask, err error) {
	const q = `
		SELECT
			id, "type", "name", "desc", icon, premium,
			points, interval, pending, "hidden", "weight",
			max_clicks,
			created_at, updated_at,
			action_link, action_chat_id,
			action_ton_amount,
			action_stars_amount, action_stars_title, action_stars_desc, action_stars_item,
			action_partner_hook, action_partner_match,
			action_tapp_ads_token, action_ads_gram_block_id
		FROM tasks
		WHERE deleted_at is null
		ORDER BY points DESC
	`
	rows, err := db.c.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t common.FullTask
		err = rows.Scan(
			&t.TaskID, &t.Type, &t.Name, &t.Desc, &t.Icon, &t.Premium,
			&t.Points, &t.Interval, &t.Pending, &t.Hidden, &t.Weight,
			&t.MaxClicks,
			&t.CreatedAt, &t.UpdatedAt,
			&t.ActionLink, &t.ActionChatID,
			&t.ActionTonAmount,
			&t.ActionStarsAmount, &t.ActionStarsTitle, &t.ActionStarsDesc, &t.ActionStarsItem,
			&t.ActionPartnerHook, &t.ActionPartnerMatch,
			&t.ActionTappAdsToken, &t.ActionAdsGramBlockID,
		)
		if err != nil {
			break
		}
		tasks = append(tasks, &t)
	}
	if tasks == nil {
		tasks = []*common.FullTask{}
	}
	return tasks, err
}

func (db *DB) CreateTask(ctx context.Context, t *common.FullTask) error {
	const q = `
		INSERT INTO tasks (
			"type", "name", "desc", icon, premium,
			points, interval, pending,
			"hidden", "weight",
			action_link, action_chat_id,
			action_ton_amount,
			action_stars_amount, action_stars_title, action_stars_desc, action_stars_item,
			action_partner_hook, action_partner_match,
			action_tapp_ads_token, action_ads_gram_block_id,
			max_clicks
		)
		VALUES ($1,$2,$3,$4,$5, $6,$7,$8, $9,$10, $11,$12, $13, $14,$15,$16,$17, $18,$19, $20,$21, $22)
		RETURNING id, created_at
	`
	row := db.c.QueryRow(
		ctx, q,
		t.Type, t.Name, t.Desc, t.Icon, t.Premium,
		t.Points, t.Interval, t.Pending,
		t.Hidden, t.Weight,
		t.ActionLink, t.ActionChatID,
		t.ActionTonAmount,
		t.ActionStarsAmount, t.ActionStarsTitle, t.ActionStarsDesc, t.ActionStarsItem,
		t.ActionPartnerHook, t.ActionPartnerMatch,
		t.ActionTappAdsToken, t.ActionAdsGramBlockID,
		t.MaxClicks,
	)
	return row.Scan(&t.TaskID, &t.CreatedAt)
}

func (db *DB) UpdateTask(ctx context.Context, t *common.FullTask) error {
	const q = `
		UPDATE tasks SET
			updated_at = NOW(),
			"type" = $1, "name" = $2, "desc" = $3, icon = $4, premium = $5,
			points = $6, interval = $7, pending = $8,
			"hidden" = $9, "weight" = $10,
			action_link = $11, action_chat_id = $12,
			action_ton_amount = $13,
			action_stars_amount = $14, action_stars_title = $15, action_stars_desc = $16, action_stars_item = $17,
			action_partner_hook = $18, action_partner_match = $19,
			action_tapp_ads_token = $20, action_ads_gram_block_id = $21,
			max_clicks = $22
		WHERE id = $23
		RETURNING updated_at
	`
	row := db.c.QueryRow(
		ctx, q,
		t.Type, t.Name, t.Desc, t.Icon, t.Premium,
		t.Points, t.Interval, t.Pending,
		t.Hidden, t.Weight,
		t.ActionLink, t.ActionChatID,
		t.ActionTonAmount,
		t.ActionStarsAmount, t.ActionStarsTitle, t.ActionStarsDesc, t.ActionStarsItem,
		t.ActionPartnerHook, t.ActionPartnerMatch,
		t.ActionTappAdsToken, t.ActionAdsGramBlockID,
		t.MaxClicks,
		t.TaskID,
	)
	return row.Scan(&t.UpdatedAt)
}

func (db *DB) DeleteTask(ctx context.Context, taskID int64) error {
	const q = `UPDATE tasks SET deleted_at = NOW() WHERE id = $1`
	_, err := db.c.Exec(ctx, q, taskID)
	return err
}
