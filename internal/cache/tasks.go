package cache

import (
	"context"
	"strconv"
)

func (Cache) keyTaskClicks(taskID int64) string {
	return "task_clicks:" + strconv.FormatInt(taskID, 10)
}

func (c Cache) GetTaskClicks(ctx context.Context, taskID int64) (int64, error) {
	s, err := c.c.Get(ctx, c.keyTaskClicks(taskID)).Result()
	if err != nil {
		if isNil(err) {
			return 0, nil
		}
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}

func (c Cache) SetTaskClicks(ctx context.Context, taskID, clicksCount int64) error {
	return c.c.Set(ctx, c.keyTaskClicks(taskID), clicksCount, 0).Err()
}

func (c Cache) IncTaskClicks(ctx context.Context, taskID int64) error {
	return c.c.Incr(ctx, c.keyTaskClicks(taskID)).Err()
}

func (c Cache) DecTaskClicks(ctx context.Context, taskID int64) error {
	return c.c.Decr(ctx, c.keyTaskClicks(taskID)).Err()
}
