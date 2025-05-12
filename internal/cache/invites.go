package cache

import (
	"context"
	"strconv"
	"time"
)

func (Cache) inviteMessageKey(userID int64) string {
	return "invite_message:" + strconv.FormatInt(userID, 10)
}

func (c Cache) GetInviteMessage(ctx context.Context, userID int64) (msgID string, err error) {
	msgID, err = c.c.Get(ctx, c.inviteMessageKey(userID)).Result()
	if isNil(err) {
		return "", nil
	}
	return
}

func (c Cache) SaveInviteMessage(ctx context.Context, userID int64, msgID string, ttl time.Duration) error {
	return c.c.Set(ctx, c.inviteMessageKey(userID), msgID, ttl).Err()
}
