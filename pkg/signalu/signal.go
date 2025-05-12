package signalu

import (
	"context"
	"os"
	"os/signal"
)

func WaitExitContext(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, os.Interrupt, os.Kill)
}
