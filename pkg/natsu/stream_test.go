package natsu

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/avpetkun/the-prime/pkg/signalu"
)

func TestStream(t *testing.T) {
	t.SkipNow()

	ctx, _ := signalu.WaitExitContext(context.TODO())

	stream, err := Connect(zerolog.New(os.Stdout), "nats://127.0.0.1:4222")
	require.NoError(t, err)
	defer stream.Stop()

	a := 0
	b := 0
	err = stream.Subscribe(ctx, 2, map[string]Handler{
		"test1.a": func(ctx context.Context, log zerolog.Logger, msg []byte) (retry time.Duration, err error) {
			fmt.Printf("test1.a: %s\n", msg)
			a++
			if a < 5 {
				return time.Second, nil
			}
			return
		},
		"test1.b": func(ctx context.Context, log zerolog.Logger, msg []byte) (retry time.Duration, err error) {
			fmt.Printf("test1.b: %s\n", msg)
			b++
			if b < 2 {
				return time.Second, nil
			}
			return
		},
		"test2.c": func(ctx context.Context, log zerolog.Logger, msg []byte) (retry time.Duration, err error) {
			fmt.Printf("test2.c: %s\n", msg)
			return
		},
	})
	require.NoError(t, err)

	err = stream.Subscribe(ctx, 1, map[string]Handler{
		"test2.d": func(ctx context.Context, log zerolog.Logger, msg []byte) (retry time.Duration, err error) {
			return
		},
	})
	require.Error(t, err)

	err = stream.Subscribe(ctx, 1, map[string]Handler{
		"test3.e": func(ctx context.Context, log zerolog.Logger, msg []byte) (retry time.Duration, err error) {
			fmt.Printf("test3.e: %s\n", msg)
			return
		},
	})
	require.NoError(t, err)

	time.Sleep(time.Second)

	err = stream.Publish(ctx, "test1.a", "aaa")
	require.NoError(t, err)
	err = stream.Publish(ctx, "test1.b", "bbb")
	require.NoError(t, err)
	err = stream.Publish(ctx, "test2.c", "ccc1")
	require.NoError(t, err)
	err = stream.Publish(ctx, "test2.c", "ccc2")
	require.NoError(t, err)
	err = stream.Publish(ctx, "test3.e", "eee")
	require.NoError(t, err)

	<-ctx.Done()
}
