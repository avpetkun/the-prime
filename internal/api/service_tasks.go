package api

import (
	"context"
	"fmt"
	"time"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/pkg/timeu"
)

func (s *Service) getTask(taskID int64) (*common.FullTask, error) {
	for _, t := range s.allTasks {
		if t.TaskID == taskID {
			return t, nil
		}
	}
	return nil, fmt.Errorf("task %d not found", taskID)
}

func (s *Service) startLoadTasksLoop(ctx context.Context) error {
	tasks, err := s.GetAllTasksWithClicks(ctx)
	if err != nil {
		return err
	}
	s.allTasks = tasks

	go func() {
		for !timeu.SleepContext(ctx, time.Second*10) {
			tasks, err = s.GetAllTasksWithClicks(ctx)
			if err != nil {
				s.log.Error().Err(err).Msg("failed to load tasks")
			} else {
				s.allTasks = tasks
			}
		}
	}()
	return nil
}
