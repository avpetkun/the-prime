package worker

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/mymmrac/telego"
	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/tgu"
)

type ProductClaimMessage struct {
	Product *common.Product `json:"product"`
	User    tgu.User        `json:"user"`
	ClaimAt time.Time       `json:"claimAt"`
}

func (s *Service) startProductsWorkers(ctx context.Context) error {
	return s.ns.Subscribe(ctx, 1, map[string]natsu.Handler{
		KeyProductClaim: s.productClaimHandler,
	})
}

func (s *Service) productClaimHandler(ctx context.Context, log zerolog.Logger, data []byte) (time.Duration, error) {
	var msg ProductClaimMessage

	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Error().Err(err).Msg("[products] can't unmarshal msg claim")
		return 0, nil
	}

	var ticketID int64

	err = s.db.Tx(ctx, func(tx *dbx.DB) error {
		ticketID, err = tx.SaveProductTicket(ctx, msg.User.ID, msg.ClaimAt, msg.Product)
		if err != nil {
			return err
		}
		if ticketID == 0 {
			return io.EOF
		}
		return tx.WithdrawUserPoints(ctx, msg.User.ID, msg.Product.Price)
	})
	if err != nil && !errors.Is(err, io.EOF) {
		return time.Second, err
	}

	switch msg.Product.Type {
	case common.ProductTgStars, common.ProductTgPremium:
		err = s.ns.Publish(ctx, KeyFragmentSend, FragmentMessage{
			ProductClaimMessage: msg,
			TicketID:            ticketID,
		})
		if err != nil {
			return 0, err
		}
	}

	var text string
	switch msg.User.LanguageCode {
	case "ru":
		text = "Ваш запрос на вывод средств принят и в настоящее время обрабатывается.\nВаше вознаграждение будет зачислено вам в ближайшее время."
	default:
		text = "Your withdrawal request has been accepted and is currently being processed.\nYour reward will be credited to you soon."
	}
	_, err = s.bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: msg.User.ID},
		Text:   text,
	})
	if err != nil {
		log.Error().Err(err).Msg("[products] failed to send user withdraw push")
	}

	return 0, nil
}
