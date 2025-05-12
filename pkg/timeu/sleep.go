package timeu

import (
	"context"
	"time"
)

func SleepContext(ctx context.Context, d time.Duration) (ctxDone bool) {
	select {
	case <-ctx.Done():
		return true
	case <-time.After(d):
		return false
	}
}
