package worker

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/tgu"
)

type UserNewMessage struct {
	tgu.User
	JoinAt     time.Time `json:"joined"`
	IPAddress  string    `json:"ip,omitempty"`
	UserAgent  string    `json:"ua,omitempty"`
	StartParam string    `json:"start,omitempty"`
}

func (s *Service) startUsersWorkers(ctx context.Context) error {
	return s.ns.Subscribe(ctx, s.cfg.Workers.Users, map[string]natsu.Handler{
		KeyUserNew:  s.usersNewHandler,
		KeyUserInit: s.usersInitHandler,
	})
}

func (s *Service) usersNewHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	var user UserNewMessage
	if err := json.Unmarshal(data, &user); err != nil {
		log.Error().Err(err).Msg("[users] can't unmarshal msg new")
		return 0, nil
	}

	exist, err := s.ch.CheckUser(ctx, user.ID)
	if exist || err != nil {
		return 0, err
	}

	refID, _ := strconv.ParseInt(user.StartParam, 10, 64)
	if refID > 0 {
		exist, err = s.ch.CheckUser(ctx, refID)
		if err != nil {
			return 0, err
		}
		inviteTask := s.getInviteTask()
		if exist && inviteTask != nil {
			err = s.ns.Publish(ctx, KeyTaskDone, TaskMessage{
				Time:   time.Now(),
				UserID: refID,
				TaskID: inviteTask.TaskID,
				SubID:  0,
				Strict: true,
			})
			if err != nil {
				return 0, err
			}
		}
	}

	err = s.db.Tx(ctx, func(tx *dbx.DB) error {
		if refID > 0 {
			if err = tx.IncUserRefCount(ctx, refID); err != nil {
				return err
			}
		}
		return tx.CreateUser(ctx, user.JoinAt, user.User, refID, user.IPAddress, user.UserAgent)
	})
	if err != nil {
		return 0, err
	}
	if err = s.ch.AddUser(ctx, user.ID); err != nil {
		return 0, err
	}
	if refID > 0 {
		if err = s.ch.IncUserRefCount(ctx, refID); err != nil {
			return 0, err
		}
	}

	log.Info().Msg("[users] new user joined")
	return 0, nil
}

type UserInitMessage struct {
	Time time.Time `json:"time"`
	User int64     `json:"user"`
}

func (s *Service) usersInitHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	var msg UserInitMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("[users] can't unmarshal msg init")
		return 0, nil
	}

	inited, err := s.ch.CheckUserInit(ctx, msg.User)
	if inited || err != nil {
		return 0, err
	}
	if err = s.db.SetUserInit(ctx, msg.User); err != nil {
		return 0, err
	}
	if err = s.ch.AddUserInit(ctx, msg.User); err != nil {
		return 0, err
	}

	log.Info().Msg("[users] new user inited")
	return 0, nil
}
