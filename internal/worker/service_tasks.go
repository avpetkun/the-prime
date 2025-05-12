package worker

import (
	"context"
	"time"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/pkg/timeu"
)

func (s *Service) findTaskByID(taskID int64) *common.FullTask {
	for _, t := range s.allTasks {
		if t.TaskID == taskID {
			return t
		}
	}
	return nil
}

func (s *Service) findTaskByTokenAndID(taskID int64, token string) *common.FullTask {
	for _, t := range s.allTasks {
		if t.ActionPartnerHook == token && (taskID == 0 || taskID == t.TaskID) {
			return t
		}
	}
	return nil
}

func (s *Service) getInviteTask() *common.FullTask {
	for _, t := range s.allTasks {
		if t.Type == common.TaskInvite {
			return t
		}
	}
	return nil
}

func (s *Service) startLoadTasksLoop(ctx context.Context) error {
	tasks, err := s.db.GetAllTasks(ctx)
	if err != nil {
		return err
	}
	s.allTasks = tasks

	go func() {
		for !timeu.SleepContext(ctx, time.Second*10) {
			tasks, err = s.db.GetAllTasks(ctx)
			if err != nil {
				s.log.Error().Err(err).Msg("failed to load tasks")
			} else {
				s.allTasks = tasks
			}
		}
	}()
	return nil
}
