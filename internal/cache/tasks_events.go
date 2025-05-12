package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/avpetkun/the-prime/internal/common"
)

func (Cache) keyUserTasksEvents(userID int64) string {
	return "tasks_events:" + strconv.FormatInt(userID, 10)
}

func (c Cache) PopUserTasksEvents(ctx context.Context, userID int64) ([]common.TaskEvent, error) {
	_, list, err := c.c.LMPop(ctx, "LEFT", 10, c.keyUserTasksEvents(userID)).Result()
	if err != nil {
		if isNil(err) {
			return nil, nil
		}
		return nil, err
	}
	events := make([]common.TaskEvent, len(list))
	for i, raw := range list {
		var e common.TaskEvent
		if err = json.Unmarshal([]byte(raw), &e); err != nil {
			return nil, err
		}
		events[i] = e
	}
	return events, nil
}

func (c Cache) PushUserTasksEvent(ctx context.Context, userID int64, e common.TaskEvent) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return c.c.RPush(ctx, c.keyUserTasksEvents(userID), string(data)).Err()
}
