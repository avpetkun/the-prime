package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/tgu"
)

func (s *Service) startTelegramWorkers(ctx context.Context) error {
	return s.ns.Subscribe(ctx, s.cfg.Workers.Chats, map[string]natsu.Handler{
		KeyChatCheckUser: s.chatsCheckUserHandler,
	})
}

func (s *Service) chatsCheckUserHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	const retryInterval = time.Minute
	const errorInterval = time.Second * 5

	var msg TaskMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("[chats] can't unmarshal msg task")
		return 0, nil
	}

	task := s.findTaskByID(msg.TaskID)
	if task == nil {
		log.Error().Msg("[chats] task not found")
		return 0, nil
	}
	if task.ActionChatID == 0 {
		log.Error().Msg("[chats] task chat_id is empty")
		return 0, nil
	}

	chatFound, memberExist, err := s.bot.CheckChatMember(ctx, task.ActionChatID, msg.UserID)
	if err != nil {
		if ok, retry := tgu.ParseErrTooManyRequests(err); ok {
			return retry, nil
		}
		return errorInterval, err
	}
	if !chatFound {
		return 0, nil // skip no-chat tasks
	}

	log.Info().Msgf("[chats] check status %t", memberExist)

	if !memberExist {
		if msg.Time.Unix()+task.Pending < time.Now().Unix() {
			return 0, nil // pop expired
		}
		return retryInterval, nil
	}

	msg.Time = time.Now()
	return 0, s.ns.Publish(ctx, KeyTaskDone, msg)
}
