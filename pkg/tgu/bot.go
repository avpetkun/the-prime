package tgu

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mymmrac/telego"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Bot struct {
	log zerolog.Logger
	api *telego.Bot
	cfg BotConfig

	myID int64
	auth *AuthParser

	token []byte

	msgLimit *rate.Limiter
	apiLimit *rate.Limiter

	LoadOffset func(ctx context.Context) (offset int, err error)
	SaveOffset func(ctx context.Context, offset int) error

	OnUpdate func(ctx context.Context, upd telego.Update) error

	OnPrivateText func(ctx context.Context, user User, msg *telego.Message) error

	OnSelfJoin   func(ctx context.Context, chat Chat) error
	OnSelfKick   func(ctx context.Context, chat Chat) error
	OnSelfMember func(ctx context.Context, joined bool, chat Chat) error

	OnUserJoin   func(ctx context.Context, user User) error
	OnUserBan    func(ctx context.Context, user User) error
	OnUserMember func(ctx context.Context, joined bool, user User) error

	OnUpdateChat func(ctx context.Context, chat Chat) error
	OnStarsTx    func(ctx context.Context, txID, payload string, amount int) (ok bool, err error)
}

func CreateBot(log zerolog.Logger, cfg BotConfig) (*Bot, error) {
	if cfg.Token == "" {
		return nil, errors.New("bot token is empty")
	}
	if cfg.RateMsgLimit <= 0 {
		cfg.RateMsgLimit = 25
	}
	if cfg.RateApiLimit <= 0 {
		cfg.RateApiLimit = 25
	}
	bot, err := telego.NewBot(cfg.Token, telego.WithDiscardLogger())
	if err != nil {
		return nil, err
	}
	me, err := bot.GetMe()
	if err != nil {
		return nil, err
	}
	tb := &Bot{
		log: log,
		api: bot,
		cfg: cfg,

		myID: me.ID,
		auth: NewAuthParser(cfg.Token),

		token: []byte(cfg.Token),

		msgLimit: rate.NewLimiter(rate.Every(time.Millisecond*time.Duration(1000/cfg.RateMsgLimit)), 1),
		apiLimit: rate.NewLimiter(rate.Every(time.Millisecond*time.Duration(1000/cfg.RateApiLimit)), 1),
	}
	return tb, nil
}

func (b *Bot) AuthParseBase64(webappInitDataBase64 string) (*Auth, error) {
	return b.auth.ParseBase64(webappInitDataBase64)
}

func (b *Bot) AuthParse(webappInitData string) (*Auth, error) {
	return b.auth.Parse(webappInitData)
}

func (b *Bot) GetMe() (User, error) {
	me, err := b.api.GetMe()
	if err != nil {
		return User{}, err
	}
	return userTelego(*me), nil
}

func (b *Bot) CheckChatMember(ctx context.Context, chatID, userID int64) (chatFound, memberExist bool, err error) {
	b.apiLimit.Wait(ctx)
	member, err := b.api.GetChatMember(&telego.GetChatMemberParams{
		ChatID: telego.ChatID{ID: chatID},
		UserID: userID,
	})
	if err != nil {
		errText := err.Error()
		if strings.Contains(errText, "PARTICIPANT_ID_INVALID") {
			// b.log.Warn().Err(err).
			// 	Int64("chat_id", chatID).
			// 	Int64("user_id", userID).
			// 	Msg("[bot] failed to check chat member")
			return false, false, nil
		}
		if strings.Contains(errText, "chat not found") {
			return false, false, nil
		}
		return false, false, err
	}
	return true, member.MemberIsMember(), nil
}

func (b *Bot) SendTextMessage(ctx context.Context, chatID int64, text string) (*telego.Message, error) {
	return b.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: chatID},
		Text:   text,
	})
}

func (b *Bot) SendMessage(ctx context.Context, params *telego.SendMessageParams) (*telego.Message, error) {
	b.msgLimit.Wait(ctx)
	return b.api.SendMessage(params)
}

func (b *Bot) SendPhoto(ctx context.Context, params *telego.SendPhotoParams) (*telego.Message, error) {
	b.msgLimit.Wait(ctx)
	return b.api.SendPhoto(params)
}

func (b *Bot) ForwardMessage(ctx context.Context, params *telego.ForwardMessageParams) (*telego.Message, error) {
	b.msgLimit.Wait(ctx)
	return b.api.ForwardMessage(params)
}

func (b *Bot) CreateStarsInvoice(
	ctx context.Context, starsAmount int,
	payload, title, description, itemLabel string,
) (invoiceURL string, err error) {
	b.apiLimit.Wait(ctx)
	link, err := b.api.CreateInvoiceLink(&telego.CreateInvoiceLinkParams{
		Title:       title,
		Description: description,
		Payload:     payload,
		Currency:    "XTR",
		Prices: []telego.LabeledPrice{
			{Label: itemLabel, Amount: starsAmount},
		},
	})
	if err != nil {
		return "", err
	}
	return *link, nil
}

func (b *Bot) PrepareInviteMessage(
	ctx context.Context, userID int64,
	text, photoURL, buttonText, miniappURL string,
) (*telego.PreparedInlineMessage, error) {
	b.apiLimit.Wait(ctx)
	return b.api.SavePreparedInlineMessage(&telego.SavePreparedInlineMessageParams{
		UserID: userID,

		AllowUserChats:    true,
		AllowBotChats:     true,
		AllowGroupChats:   true,
		AllowChannelChats: true,

		Result: &telego.InlineQueryResultPhoto{
			Type: telego.ResultTypePhoto,

			ID: uuid.NewString(),

			PhotoURL:     photoURL,
			ThumbnailURL: photoURL,

			Caption: text,

			ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: [][]telego.InlineKeyboardButton{{
				{
					Text: buttonText,
					URL:  miniappURL + "?startapp=" + strconv.FormatInt(userID, 10),
				},
			}}},
		},
	})
}
