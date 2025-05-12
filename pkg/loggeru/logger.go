package loggeru

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetLogger(ctx context.Context) (context.Context, zerolog.Logger) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = false
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	logger := zerolog.New(os.Stdout)
	log.Logger = logger
	return logger.WithContext(ctx), logger
}
