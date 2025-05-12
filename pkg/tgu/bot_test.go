package tgu

import (
	"context"
	"os"
	"testing"

	"github.com/mymmrac/telego"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/avpetkun/the-prime/pkg/signalu"
)

func TestListener(t *testing.T) {
	t.SkipNow()

	ctx, _ := signalu.WaitExitContext(context.Background())

	b, err := CreateBot(zerolog.New(os.Stdout), BotConfig{
		Token: os.Getenv("BOT_TOKEN"),
	})
	require.NoError(t, err)

	b.OnUpdate = func(ctx context.Context, upd telego.Update) error {
		return nil
	}

	b.OnUpdateChat = func(ctx context.Context, chat Chat) error {
		log.Printf("on update chat %+v", chat)
		return nil
	}

	b.OnPrivateText = func(ctx context.Context, user User, msg *telego.Message) error {
		log.Print("on private text ", user, " ", msg.Text)
		return nil
	}

	b.OnSelfJoin = func(ctx context.Context, chat Chat) error {
		log.Printf("OnSelfJoin chat %+v", chat)
		return nil
	}
	b.OnSelfKick = func(ctx context.Context, chat Chat) error {
		log.Printf("OnSelfKick chat %+v", chat)
		return nil
	}
	b.OnSelfMember = func(ctx context.Context, joined bool, chat Chat) error {
		log.Printf("OnSelfMember chat %+v joined %t", chat, joined)
		return nil
	}

	b.OnUserJoin = func(ctx context.Context, user User) error {
		log.Printf("OnUserJoin user %+v", user)
		return nil
	}
	b.OnUserBan = func(ctx context.Context, user User) error {
		log.Printf("OnUserBan user %+v", user)
		return nil
	}
	b.OnUserMember = func(ctx context.Context, joined bool, user User) error {
		log.Printf("OnUserMember user %+v joined %t", user, joined)
		return nil
	}

	require.NoError(t, b.Listen(ctx))
	println("listen")

	<-ctx.Done()
}
