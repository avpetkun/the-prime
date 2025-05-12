package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/pkg/natsu"
)

type WebhookRewardMessage struct {
	Received time.Time `json:"received"`
	Token    string    `json:"token"`
	UserID   int64     `json:"userID"`
	TaskID   int64     `json:"taskID,omitempty"`
	SubID    int64     `json:"subID,omitempty"`
	Payout   float64   `json:"payout,omitempty"`
}

func (s *Service) startWebhookWorkers(ctx context.Context) error {
	return s.ns.Subscribe(ctx, s.cfg.Workers.Webhook, map[string]natsu.Handler{
		KeyWebhookReward: s.webhookRewardHandler,
	})
}

func (s *Service) webhookRewardHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	var reward WebhookRewardMessage
	if err := json.Unmarshal(data, &reward); err != nil {
		log.Error().Err(err).Msg("[webhooks] can't unmarshal reward")
		return 0, nil
	}

	fullTask := s.findTaskByTokenAndID(reward.TaskID, reward.Token)
	if fullTask == nil {
		log.Warn().Msg("[webhooks] can't find task for reward")
		return 0, nil
	}
	task := fullTask.UserTask
	task.SubID = reward.SubID

	err := s.db.SaveTaskPartnerEvent(ctx, reward.UserID, task.TaskID, task.SubID, reward.Payout, reward.Received)
	if err != nil {
		return 0, err
	}

	err = s.ns.Publish(ctx, KeyTaskDone, TaskMessage{
		Time:   reward.Received,
		UserID: reward.UserID,
		TaskID: task.TaskID,
		SubID:  task.SubID,
		Strict: true,
	})
	if err == nil {
		log.Info().
			Int64("user_id", reward.UserID).
			Int64("task_id", task.TaskID).
			Int64("sub_id", task.SubID).
			Float64("payout", reward.Payout).
			Int64("task_points", task.Points).
			Str("task_type", string(task.Type)).
			Msg("[webhooks] processed partner task")
	}
	return 0, err
}
