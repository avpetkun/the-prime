package worker

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/internal/cache"
	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/pkg/natsu"
)

type TaskMessage struct {
	Time   time.Time `json:"ts"`
	UserID int64     `json:"uid"`
	TaskID int64     `json:"tid"`
	SubID  int64     `json:"sid,omitempty"`
	Strict bool      `json:"strict,omitempty"`
}

func (s *Service) startTasksWorkers(ctx context.Context) error {
	return s.ns.Subscribe(ctx, s.cfg.Workers.Tasks, map[string]natsu.Handler{
		KeyTaskStart: s.tasksStartHandler,
		KeyTaskDone:  s.tasksDoneHandler,
		KeyTaskClaim: s.tasksClaimHandler,
	})
}

func (s *Service) tasksStartHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	var msg TaskMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("[tasks] can't unmarshal msg task")
		return 0, nil
	}

	task := s.findTaskByID(msg.TaskID)
	if task == nil {
		if time.Since(msg.Time) > time.Minute {
			log.Warn().Msg("[tasks] can't find task for msg")
			return 0, nil
		}
		return time.Second * 10, nil
	}

	err := s.db.SaveTaskStatePending(ctx, msg.UserID, msg.TaskID, msg.SubID, msg.Time)
	if err != nil {
		return 0, err
	}

	switch task.Type {
	case common.TaskFree, common.TaskFreeLink, common.TaskTonConnect, common.TaskTonDisconnect:
		msg.Time = time.Now()
		err = s.ns.Publish(ctx, KeyTaskDone, msg)
	case common.TaskJoin:
		err = s.ns.Publish(ctx, KeyChatCheckUser, msg)
	case common.TaskPartnerCheck:
		err = s.ns.Publish(ctx, KeyCheckPartnerTask, msg)
	}
	if err == nil {
		log.Info().Msg("[tasks] task started")
	}
	return 0, err
}

func (s *Service) tasksDoneHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	var msg TaskMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("[tasks] can't unmarshal msg task")
		return 0, nil
	}

	task := s.findTaskByID(msg.TaskID)
	if task == nil {
		log.Warn().Msg("[tasks] can't find task for msg")
		return 0, nil
	}
	if task.Hidden {
		return 0, nil
	}

	exist, err := s.ch.CheckUser(ctx, msg.UserID)
	if err != nil {
		return 0, err
	}
	if !exist && msg.Strict {
		return 0, nil
	}

	var flow *common.TaskFlow
	if exist {
		flow, err = s.ch.GetUserTask(ctx, msg.UserID, msg.TaskID, msg.SubID)
		if err != nil {
			return 0, err
		}
		if flow == nil || flow.Status != common.TaskPending {
			if time.Since(msg.Time) < time.Minute {
				return time.Second, nil
			}
			if msg.Strict {
				return 0, nil
			}
		}
	}
	if flow == nil {
		flow = &common.TaskFlow{
			TaskKey: common.TaskKey{TaskID: msg.TaskID, SubID: msg.SubID},
			Start:   msg.Time.Unix(),
		}
	}

	referrals := s.calculateRefRewards(ctx, msg.UserID, task.Points)

	err = s.db.Tx(ctx, func(tx *dbx.DB) error {
		ok, err := tx.SaveTaskStateClaim(ctx, msg.UserID, msg.TaskID, msg.SubID, task.Points, msg.Time)
		if err != nil {
			return err
		}
		if !ok {
			return io.EOF
		}
		for _, ref := range referrals {
			if err = tx.IncUserPointsFromRef(ctx, ref.UserID, ref.Points); err != nil {
				return err
			}
			if err = tx.SaveUserRefferal(ctx, msg.UserID, ref.UserID, ref.Points, ref.Level); err != nil {
				return err
			}
		}
		if err = tx.IncUserPoints(ctx, msg.UserID, task.Points); err != nil {
			return err
		}
		return tx.SaveTaskFinish(ctx, msg.UserID, msg.TaskID, msg.SubID, task.Points, msg.Time)
	})
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}

	err = s.ch.Tx(ctx, func(c cache.Cache) error {
		flow.Status = common.TaskClaim
		if err = c.SetUserTask(ctx, msg.UserID, flow); err != nil {
			return err
		}
		if err = c.IncUserPoints(ctx, msg.UserID, task.Points); err != nil {
			return err
		}
		for _, ref := range referrals {
			if err = c.IncUserPoints(ctx, ref.UserID, ref.Points); err != nil {
				return err
			}
			if err = c.IncUserRefPoints(ctx, ref.UserID, ref.Points); err != nil {
				return err
			}
		}
		return c.PushUserTasksEvent(ctx, msg.UserID, common.TaskEvent{
			TaskKey: flow.TaskKey,
			Time:    time.Now().Unix(),
			Status:  common.TaskClaim,
		})
	})
	if err == nil {
		log.Info().Msg("[tasks] task done")
	}
	return 0, err
}

func (s *Service) tasksClaimHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	var msg TaskMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("[tasks] can't unmarshal msg task")
		return 0, nil
	}
	if err := s.db.SaveTaskStateDone(ctx, msg.UserID, msg.TaskID, msg.SubID, msg.Time); err != nil {
		return 0, err
	}

	log.Info().Msg("[tasks] task claimed")
	return 0, nil
}
