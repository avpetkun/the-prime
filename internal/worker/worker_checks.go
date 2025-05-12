package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/pkg/natsu"
)

func (s *Service) startChecksWorkers(ctx context.Context) error {
	return s.ns.Subscribe(ctx, s.cfg.Workers.Checks, map[string]natsu.Handler{
		KeyCheckPartnerTask: s.checkPartnerTaskHandler,
	})
}

func (s *Service) checkPartnerTaskHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	const retryInterval = time.Second * 30

	var msg TaskMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("[checks] can't unmarshal msg task")
		return 0, nil
	}

	task := s.findTaskByID(msg.TaskID)
	if task == nil {
		log.Error().Msg("[checks] task not found")
		return 0, nil
	}

	checkParts := strings.Fields(task.ActionPartnerHook)
	if len(checkParts) != 2 {
		log.Error().Msg("[checks] invalid task.ActionPartnerHook")
		return 0, nil
	}

	req, err := http.NewRequestWithContext(ctx, checkParts[0], checkParts[1], nil)
	if err != nil {
		log.Error().Err(err).Msg("[checks] invalid request for task.ActionPartnerHook")
		return 0, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var success bool
	if task.ActionPartnerMatch == "" {
		io.Copy(io.Discard, resp.Body)
		success = 200 <= resp.StatusCode && resp.StatusCode < 300
	} else {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		success = bytes.Contains(respData, []byte(task.ActionPartnerMatch))
	}

	log.Info().Msgf("[checks] check status %t", success)

	if !success {
		if msg.Time.Unix()+task.Pending < time.Now().Unix() {
			return 0, nil // pop expired
		}
		return retryInterval, nil
	}

	msg.Time = time.Now()
	return 0, s.ns.Publish(ctx, KeyTaskDone, msg)
}
