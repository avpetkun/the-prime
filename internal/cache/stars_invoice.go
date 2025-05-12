package cache

import (
	"context"
	"strconv"
	"time"
)

func (Cache) keyStarsInvoice(userID, taskID int64) string {
	return "stars_invoice:" + strconv.FormatInt(userID, 10) + ":" + strconv.FormatInt(taskID, 10)
}

func (c Cache) GetStarsInvoice(ctx context.Context, userID, taskID int64) (invoiceURL string, err error) {
	invoiceURL, err = c.c.Get(ctx, c.keyStarsInvoice(userID, taskID)).Result()
	if isNil(err) {
		err = nil
	}
	return
}

func (c Cache) SaveStarsInvoice(ctx context.Context, userID, taskID int64, invoiceURL string, ttl time.Duration) error {
	return c.c.Set(ctx, c.keyStarsInvoice(userID, taskID), invoiceURL, ttl).Err()
}
