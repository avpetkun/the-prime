package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"time"

	"github.com/mymmrac/telego"
	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/pkg/fragment"
	"github.com/avpetkun/the-prime/pkg/natsu"
)

type FragmentMessage struct {
	ProductClaimMessage
	TicketID int64 `json:"ticketID"`
}

func (s *Service) startFragmentWorkers(ctx context.Context) error {
	return s.ns.Subscribe(ctx, 1, map[string]natsu.Handler{
		KeyFragmentSend: s.fragmentSendHandler,
	})
}

func (s *Service) fragmentSendHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	var msg FragmentMessage

	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Error().Err(err).Msg("[fragment] can't unmarshal msg")
		return 0, nil
	}

	status, err := s.db.GetProductTicketStatus(ctx, msg.TicketID)
	if err != nil {
		return 0, err
	}
	if status == "scam" {
		if s.cfg.Fragment.ChatGifts != 0 {
			s.bot.SendMessage(ctx, &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: s.cfg.Fragment.ChatGifts},
				ParseMode: telego.ModeHTML,
				Text: fmt.Sprintf(
					"#%d: @%s (id %d)\n%d %s (price %d)\n💩 it was fucking scam!",
					msg.TicketID,
					html.EscapeString(msg.User.Username), msg.User.ID,
					msg.Product.Amount, msg.Product.Type, msg.Product.Price,
				),
			})
		}
		return 0, nil
	}

	balance, err := s.frag.GetBalance(ctx)
	if err != nil {
		return time.Second, err
	}
	if balance < s.cfg.Fragment.MinBalance {
		if s.cfg.Fragment.ChatGifts != 0 && s.wasInsufficientBalance.CompareAndSwap(0, 1) {
			_, err = s.bot.SendMessage(ctx, &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: s.cfg.Fragment.ChatGifts},
				ParseMode: telego.ModeHTML,
				Text: fmt.Sprintf(
					"❗️ Insufficient fragment balance\n%.9F &lt; %.9F\n<code>%s</code>",
					balance, s.cfg.Fragment.MinBalance, s.frag.WalletAddress(),
				),
			})
			if err != nil {
				s.wasInsufficientBalance.Store(0)
			}
		}
		err = s.db.SetProductTicketStatus(ctx, msg.TicketID, "waiting")
		if err != nil {
			log.Error().Err(err).Msg("[fragment] failed to set ticket status")
		}
		log.Warn().
			Float64("curr_balance", balance).
			Float64("min_balance", s.cfg.Fragment.MinBalance).
			Msg("[fragment] check balance error")
		return time.Minute, nil
	}
	s.wasInsufficientBalance.Store(0)

	var tx *fragment.SentTx
	switch msg.Product.Type {
	case common.ProductTgStars:
		tx, err = s.frag.SendTelegramStars(ctx, log, msg.User.Username, msg.Product.Amount, false)
	case common.ProductTgPremium:
		tx, err = s.frag.SendTelegramPremium(ctx, log, msg.User.Username, msg.Product.Amount, false)
	}
	if err != nil {
		if errors.Is(err, fragment.ErrUserNotFound) {
			err = s.db.SetProductTicketStatus(ctx, msg.TicketID, "user not found")
			if err != nil {
				log.Error().Err(err).Msg("[fragment] failed to set ticket status")
			}
			if s.cfg.Fragment.ChatGifts != 0 {
				_, err = s.bot.SendMessage(ctx, &telego.SendMessageParams{
					ChatID:    telego.ChatID{ID: s.cfg.Fragment.ChatGifts},
					ParseMode: telego.ModeHTML,
					Text: fmt.Sprintf(
						"#%d: @%s (id %d)\n%d %s (price %d)\n❗️❗️❗️ User not found!",
						msg.TicketID,
						html.EscapeString(msg.User.Username), msg.User.ID,
						msg.Product.Amount, msg.Product.Type, msg.Product.Price,
					),
				})
				if err != nil {
					log.Error().Err(err).Msg("[fragment] failed to send admin withdraw err push")
				}
			}
			log.Warn().Err(err).Msg("[fragment] user not found")
			return 0, nil
		}
		err = s.db.SetProductTicketStatus(ctx, msg.TicketID, err.Error())
		if err != nil {
			log.Error().Err(err).Msg("[fragment] failed to set ticket err status")
		}
		log.Error().Err(err).Msg("[fragment] send error")
		return time.Minute, err
	}

	log.Info().Str("tx_hash", tx.Hash).Uint64("tx_lt", tx.LT).Msg("[fragment] gift was successfully sent")

	err = s.db.SetProductTicketSent(ctx, msg.TicketID, fmt.Sprintf("sent tx %s %d", tx.Hash, tx.LT))
	if err != nil {
		log.Error().Err(err).Msg("[fragment] failed to set ticket status sent")
	}

	if s.cfg.Fragment.ChatGifts != 0 {
		txLink := "https://tonviewer.com/transaction/" + tx.Hash
		_, err = s.bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: s.cfg.Fragment.ChatGifts},
			ParseMode: telego.ModeHTML,
			LinkPreviewOptions: &telego.LinkPreviewOptions{
				IsDisabled: true,
			},
			Text: fmt.Sprintf(
				"#%d: @%s (id %d)\n%d %s (price %d)\n✅ Successfully <a href=\"%s\">sent</a>",
				msg.TicketID,
				html.EscapeString(msg.User.Username), msg.User.ID,
				msg.Product.Amount, msg.Product.Type, msg.Product.Price,
				txLink,
			),
		})
		if err != nil {
			log.Error().Err(err).Str("tx_hash", tx.Hash).Msg("[fragment] failed to send admin withdraw ok push")
		}
	}

	return 0, nil
}
