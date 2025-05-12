package support

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/pkg/tgu"
)

type Service struct {
	cfg Config
	log zerolog.Logger
	bot *tgu.Bot

	userByMessage map[int]int64

	lastUserID int64
}

func NewService(cfg Config, log zerolog.Logger) (*Service, error) {
	bot, err := tgu.CreateBot(log, tgu.BotConfig{
		Token: cfg.BotToken,
	})
	if err != nil {
		return nil, err
	}
	s := &Service{cfg: cfg, log: log, bot: bot, userByMessage: make(map[int]int64)}

	bot.OnUpdate = s.processUpdate
	bot.LoadOffset = func(ctx context.Context) (offset int, err error) {
		// TODO: use db
		return
	}
	bot.SaveOffset = func(ctx context.Context, offset int) error {
		// TODO: use db
		return nil
	}

	return s, nil
}

func (s *Service) Start(ctx context.Context) error {
	return s.bot.Listen(ctx)
}

func (s *Service) processUpdate(ctx context.Context, u telego.Update) error {
	msg := u.Message
	if msg == nil || msg.Text == "/start" {
		return nil
	}
	if msg.Chat.Type == telego.ChatTypePrivate {
		s.log.Info().Any("message", msg).Msg("received user message")

		if s.lastUserID != msg.Chat.ID {
			s.lastUserID = msg.Chat.ID

			s.bot.SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: s.cfg.ChatID},
				Text:   fmt.Sprintf("👤 id %d username @%s", msg.Chat.ID, msg.Chat.Username),
			})
		}

		suportMsg, err := s.bot.ForwardMessage(ctx, &telego.ForwardMessageParams{
			ChatID:     telego.ChatID{ID: s.cfg.ChatID},
			FromChatID: msg.Chat.ChatID(),
			MessageID:  msg.MessageID,
		})
		if err != nil {
			s.log.Error().Err(err).Any("message", msg).Msg("failed to forward message")
			return nil
		}
		s.userByMessage[suportMsg.MessageID] = msg.Chat.ID

		_, err = s.bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: telego.ChatID{ID: msg.Chat.ID},
			Text:   "✅ Your message has been forwarded to our support team.\nWe'll reply soon!",
		})
		if err != nil {
			s.log.Error().Err(err).Any("chat", msg.Chat).Msg("failed to reply to user")
		}
	} else if msg.Chat.ID == s.cfg.ChatID && msg.ReplyToMessage != nil && msg.Text != "" {
		s.log.Info().Any("message", msg).Msg("received support message")
		userID := s.userByMessage[msg.ReplyToMessage.MessageID]
		if userID == 0 {
			if user, ok := msg.ReplyToMessage.ForwardOrigin.(*telego.MessageOriginUser); ok {
				userID = user.SenderUser.ID
			}
		}
		if userID != 0 {
			_, err := s.bot.SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: userID},
				Text:   fmt.Sprintf("💬 Support team reply:\n\n%s", msg.Text),
			})
			if err != nil {
				s.log.Error().Err(err).
					Int64("user_id", userID).
					Any("message", msg).
					Msg("failed to send message to user")
			}
			s.bot.SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: s.cfg.ChatID},
				Text:   "✅ You successfully sent reply",
			})
		}
	}
	return nil
}
