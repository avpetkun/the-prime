package worker

import (
	"context"
	"strings"
	"time"

	"github.com/mymmrac/telego"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/pkg/tgu"
)

func (s *Service) OnTelegramStarsTx(ctx context.Context, txID, payload string, amount int) (ok bool, err error) {
	userID, taskID, ok := common.ParseStarsInvoicePayload(payload)
	if !ok {
		return false, nil
	}

	if err = s.db.SaveStarsTransaction(ctx, txID, userID, taskID, amount); err != nil {
		return false, err
	}

	return true, s.ns.Publish(ctx, KeyTaskDone, TaskMessage{
		Time:   time.Now(),
		UserID: userID,
		TaskID: taskID,
		SubID:  0,
	})
}

func (s *Service) OnTelegramSelfMember(ctx context.Context, joined bool, chat tgu.Chat) error {
	if joined {
		return s.db.SaveBotChat(ctx, chat)
	}
	return s.db.DeleteBotChat(ctx, chat.ChatID)
}

func (s *Service) OnTelegramPrivateMessage(ctx context.Context, user tgu.User, msg *telego.Message) error {
	msgParts := strings.Fields(msg.Text)
	if msgParts[0] != "/start" {
		return nil
	}
	if len(msgParts) > 1 {
		if exist, _ := s.ch.CheckUser(ctx, user.ID); !exist {
			err := s.ns.Publish(ctx, KeyUserNew, UserNewMessage{
				User:       user,
				JoinAt:     time.Now(),
				IPAddress:  "",
				UserAgent:  "bot",
				StartParam: msgParts[1],
			})
			if err != nil {
				s.log.Error().Err(err).
					Int64("user_id", user.ID).
					Str("start_param", msgParts[1]).
					Str("message_text", msg.Text).
					Msg("[telegram] can't publish user start")
				return err
			}
		}
	}

	loc := s.cfg.Miniapp.En
	if msg.From.LanguageCode == "ru" {
		loc = s.cfg.Miniapp.Ru
	}

	_, err := s.bot.SendPhoto(ctx, &telego.SendPhotoParams{
		ChatID:    telego.ChatID{ID: user.ID},
		ParseMode: telego.ModeHTML,
		Caption:   loc.HelloText,
		Photo:     telego.InputFile{URL: s.cfg.Miniapp.HelloImage},
		ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{
				Text: loc.WebappButton,
				WebApp: &telego.WebAppInfo{
					URL: s.cfg.Miniapp.WebappURL,
				},
			}},
			{{
				Text: loc.SupportButton,
				URL:  s.cfg.Miniapp.SupportURL,
			}},
		}},
	})
	if err != nil {
		s.log.Warn().Err(err).
			Int64("user_id", user.ID).
			Msg("[telegram] failed to send hello message")
	}
	return nil
}
