package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/avpetkun/the-prime/internal/common"
)

func (Cache) keyUserTasks(userID int64) string {
	return "user_tasks:" + strconv.FormatInt(userID, 10)
}

func (c Cache) GetUserTasks(ctx context.Context, userID int64) (map[common.TaskKey]common.TaskFlow, error) {
	results, err := c.c.HGetAll(ctx, c.keyUserTasks(userID)).Result()
	if err != nil {
		if isNil(err) {
			return nil, nil
		}
		return nil, err
	}
	tasks := make(map[common.TaskKey]common.TaskFlow, len(results))
	for _, res := range results {
		var f common.TaskFlow
		if err = json.Unmarshal([]byte(res), &f); err != nil {
			return nil, err
		}
		tasks[f.TaskKey] = f
	}
	return tasks, nil
}

func (c Cache) GetUserTask(ctx context.Context, userID, taskID, taskSubID int64) (*common.TaskFlow, error) {
	field := strconv.FormatInt(taskID, 10) + ":" + strconv.FormatInt(taskSubID, 10)
	rawFlow, err := c.c.HGet(ctx, c.keyUserTasks(userID), field).Result()
	if err != nil {
		if isNil(err) {
			return nil, nil
		}
		return nil, err
	}
	var f common.TaskFlow
	if err = json.Unmarshal([]byte(rawFlow), &f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (c Cache) SetUserTask(ctx context.Context, userID int64, f *common.TaskFlow) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	field := strconv.FormatInt(f.TaskID, 10) + ":" + strconv.FormatInt(f.SubID, 10)
	return c.c.HSet(ctx, c.keyUserTasks(userID), field, string(data)).Err()
}

func (c Cache) DelUserTask(ctx context.Context, userID, taskID, taskSubID int64) error {
	field := strconv.FormatInt(taskID, 10) + ":" + strconv.FormatInt(taskSubID, 10)
	return c.c.HDel(ctx, c.keyUserTasks(userID), field).Err()
}
