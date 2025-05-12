package tgu

import (
	"context"
	"time"

	"github.com/mymmrac/telego"

	"github.com/avpetkun/the-prime/pkg/timeu"
)

func (b *Bot) Listen(ctx context.Context) (err error) {
	startOffset := 0
	if b.LoadOffset != nil {
		startOffset, err = b.LoadOffset(ctx)
		if err != nil {
			return err
		}
	}
	if startOffset == 0 {
		b.log.Warn().Msg("[bot] can't get last offset")
	}
	go b.mainLoop(ctx, startOffset)
	return nil
}

func (b *Bot) mainLoop(ctx context.Context, startOffset int) {
	const errTimeout = time.Second * 3
	params := &telego.GetUpdatesParams{Timeout: 60, Offset: startOffset}

	for ctx.Err() == nil {
		updates, err := b.api.GetUpdates(params)
		if err != nil {
			b.log.Warn().Err(err).Msg("[bot] failed to get telegram updates")
			if timeu.SleepContext(ctx, errTimeout) {
				return
			}
			continue
		}

		for _, u := range updates {
			if u.UpdateID < params.Offset {
				continue
			}
			params.Offset = u.UpdateID + 1
			for {
				err = b.processUpdate(ctx, u)
				if err == nil {
					break
				}
				b.log.Error().Err(err).Any("update", u).Msg("[bot] can't process update")
				if timeu.SleepContext(ctx, errTimeout) {
					return
				}
			}
			for b.SaveOffset != nil {
				err = b.SaveOffset(ctx, params.Offset)
				if err == nil {
					break
				}
				b.log.Error().Err(err).Msg("[bot] can't save last offset")
				if timeu.SleepContext(ctx, errTimeout) {
					return
				}
			}
		}
	}
}

func (b *Bot) processUpdate(ctx context.Context, u telego.Update) (err error) {
	if b.OnUpdate != nil {
		if err = b.OnUpdate(ctx, u); err != nil {
			return err
		}
	}

	if m := u.MyChatMember; m != nil && m.NewChatMember != nil && m.OldChatMember != nil {
		oldIsMember := m.OldChatMember.MemberIsMember()
		newIsMember := m.NewChatMember.MemberIsMember()
		if m.Chat.Type == telego.ChatTypePrivate {
			user := userTelego(m.From)
			if !newIsMember {
				if b.OnUserBan != nil {
					if err = b.OnUserBan(ctx, user); err != nil {
						return err
					}
				}
				if b.OnUserMember != nil {
					if err = b.OnUserMember(ctx, false, user); err != nil {
						return err
					}
				}
			} else if !oldIsMember {
				if b.OnUserJoin != nil {
					if err = b.OnUserJoin(ctx, user); err != nil {
						return err
					}
				}
				if b.OnUserMember != nil {
					if err = b.OnUserMember(ctx, true, user); err != nil {
						return err
					}
				}
			}
		} else {
			chat := chatTelego(m.Chat)
			if !newIsMember {
				if b.OnSelfKick != nil {
					if err = b.OnSelfKick(ctx, chat); err != nil {
						return err
					}
				}
				if b.OnSelfMember != nil {
					if err = b.OnSelfMember(ctx, false, chat); err != nil {
						return err
					}
				}
			} else if !oldIsMember {
				if b.OnSelfJoin != nil {
					if err = b.OnSelfJoin(ctx, chat); err != nil {
						return err
					}
				}
				if b.OnSelfMember != nil {
					if err = b.OnSelfMember(ctx, true, chat); err != nil {
						return err
					}
				}
			}
		}
	}

	if b.OnPrivateText != nil && u.Message != nil &&
		u.Message.Chat.Type == telego.ChatTypePrivate &&
		u.Message.From != nil && u.Message.Text != "" {

		if err = b.OnPrivateText(ctx, userTelego(*u.Message.From), u.Message); err != nil {
			return err
		}
	}

	if b.OnUpdateChat != nil && u.ChannelPost != nil && u.ChannelPost.NewChatTitle != "" {
		if err = b.OnUpdateChat(ctx, chatTelego(u.ChannelPost.Chat)); err != nil {
			return err
		}
	}

	if q := u.PreCheckoutQuery; q != nil {
		ok := b.OnStarsTx != nil
		if ok {
			ok, err = b.OnStarsTx(ctx, q.ID, q.InvoicePayload, q.TotalAmount)
			if err != nil {
				return err
			}
		}
		return b.api.AnswerPreCheckoutQuery(&telego.AnswerPreCheckoutQueryParams{
			PreCheckoutQueryID: q.ID,
			Ok:                 ok,
		})
	}

	return nil
}
