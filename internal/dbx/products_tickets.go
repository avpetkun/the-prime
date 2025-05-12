package dbx

import (
	"context"
	"errors"
	"time"

	"github.com/avpetkun/the-prime/internal/common"
)

func (db *DB) SaveProductTicket(ctx context.Context, userID int64, claimAt time.Time, p *common.Product) (ticketID int64, err error) {
	const q = `
		INSERT INTO products_tickets (
			user_id, product_id, product_type, product_amount, claim_price, claim_at
		) VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id
	`
	err = db.c.QueryRow(ctx, q, userID, p.ID, p.Type, p.Amount, p.Price, claimAt).Scan(&ticketID)
	if isErrUniqueConstraint(err) {
		err = nil
	}
	return
}

func (db *DB) GetProductTicketStatus(ctx context.Context, ticketID int64) (status string, err error) {
	const q = `SELECT status FROM products_tickets WHERE id = $1`
	err = db.c.QueryRow(ctx, q, ticketID).Scan(&status)
	if errors.Is(err, ErrNoRows) {
		err = nil
	}
	return
}

func (db *DB) SetProductTicketStatus(ctx context.Context, ticketID int64, newStatus string) error {
	const q = `UPDATE products_tickets SET status = $1 WHERE id = $2`
	_, err := db.c.Exec(ctx, q, newStatus, ticketID)
	return err
}

func (db *DB) SetProductTicketSent(ctx context.Context, ticketID int64, newStatus string) error {
	const q = `UPDATE products_tickets SET sent_at = NOW(), status = $1 WHERE id = $2`
	_, err := db.c.Exec(ctx, q, newStatus, ticketID)
	return err
}
