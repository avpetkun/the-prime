package dbx

import (
	"context"

	"github.com/avpetkun/the-prime/pkg/tgu"
)

func (db *DB) SaveBotChat(ctx context.Context, chat tgu.Chat) error {
	const q = `
		INSERT INTO bot_chats (chat_id, title, link)
		VALUES ($1, $2, $3)
		ON CONFLICT (chat_id)
		DO UPDATE SET
			title = EXCLUDED.title,
			link = EXCLUDED.link
		WHERE
			bot_chats.chat_id = EXCLUDED.chat_id
	`
	_, err := db.c.Exec(ctx, q, chat.ChatID, chat.Title, chat.Link)
	return err
}

func (db *DB) DeleteBotChat(ctx context.Context, chatID int64) error {
	const q = `DELETE FROM bot_chats WHERE chat_id = $1`
	_, err := db.c.Exec(ctx, q, chatID)
	return err
}

func (db *DB) UpdateBotChatTitle(ctx context.Context, chatID int64, newChatTitle string) error {
	const q = `UPDATE bot_chats SET title = $1 WHERE chat_id = $2`
	_, err := db.c.Exec(ctx, q, newChatTitle, chatID)
	return err
}

func (db *DB) GetBotChatByID(ctx context.Context, chatID int64) (c tgu.Chat, err error) {
	const q = `SELECT chat_id, title, link FROM bot_chats WHERE chat_id = $1`
	err = db.c.QueryRow(ctx, q, chatID).Scan(&c.ChatID, &c.Title, &c.Link)
	return
}

func (db *DB) GetAllBotChats(ctx context.Context) (chats []tgu.Chat, err error) {
	const q = `SELECT chat_id, title, link FROM bot_chats ORDER BY joined DESC`
	rows, err := db.c.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c tgu.Chat
		if err = rows.Scan(&c.ChatID, &c.Title, &c.Link); err != nil {
			return chats, err
		}
		chats = append(chats, c)
	}
	return
}
