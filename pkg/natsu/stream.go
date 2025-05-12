package natsu

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/pkg/timeu"
)

type Handler func(ctx context.Context, log zerolog.Logger, msg []byte) (retry time.Duration, err error)

type Stream struct {
	c   *nats.Conn
	js  jetstream.JetStream
	log zerolog.Logger

	mx        sync.Mutex
	handlers  map[string]map[string]bool
	consumers []jetstream.ConsumeContext
}

func Connect(log zerolog.Logger, natsURL string) (*Stream, error) {
	conn, err := nats.Connect(natsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Warn().Err(err).Msg("[nats] disconnected")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info().Str("url", nc.ConnectedUrl()).Msg("[nats] connected")
		}),
	)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}
	stream := &Stream{
		c:   conn,
		js:  js,
		log: log,

		handlers: make(map[string]map[string]bool),
	}
	return stream, nil
}

func (s *Stream) Stop() {
	s.mx.Lock()
	list := s.consumers
	s.consumers = nil
	s.mx.Unlock()
	for _, c := range list {
		c.Stop()
	}
}

func (s *Stream) Publish(ctx context.Context, subject string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.PublishData(ctx, subject, data)
}

func (s *Stream) PublishData(ctx context.Context, subject string, payload []byte) error {
	_, err := s.js.Publish(ctx, subject, payload)
	return err
}

func splitHandlers(subjectHandlers map[string]Handler) (map[string]map[string]Handler, error) {
	streamHandlers := make(map[string]map[string]Handler)
	for subject, h := range subjectHandlers {
		if h == nil {
			return nil, fmt.Errorf("nats stream handler is nil: subject %s", subject)
		}
		i := strings.IndexByte(subject, '.')
		if i == -1 {
			return nil, fmt.Errorf("nats stream subject must contains '.': %s", subject)
		}
		stream := subject[:i]

		if streamHandlers[stream] == nil {
			streamHandlers[stream] = make(map[string]Handler)
		}
		streamHandlers[stream][subject] = h
	}
	return streamHandlers, nil
}

func (s *Stream) rememberHandlers(streamHandlers map[string]map[string]Handler) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	for stream := range streamHandlers {
		if s.handlers[stream] != nil {
			return fmt.Errorf("stream %s already exists", stream)
		}
	}
	for stream, subjects := range streamHandlers {
		s.handlers[stream] = make(map[string]bool)
		for subject := range subjects {
			s.handlers[stream][subject] = true
		}
	}
	return nil
}

func (s *Stream) Subscribe(ctx context.Context, workersPerStream int, subjectHandlers map[string]Handler) error {
	streamHandlers, err := splitHandlers(subjectHandlers)
	if err != nil {
		return err
	}
	if err = s.rememberHandlers(streamHandlers); err != nil {
		return err
	}

	for stream, subjects := range streamHandlers {
		streamAll := stream + ".*"
		str, err := s.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
			Name:      stream,
			Subjects:  []string{streamAll},
			Retention: jetstream.WorkQueuePolicy,
		})
		if err != nil {
			return err
		}
		consumer, err := str.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
			Durable:       stream + "_workers",
			FilterSubject: streamAll,
			AckPolicy:     jetstream.AckExplicitPolicy,
		})
		if err != nil {
			return err
		}
		for range workersPerStream {
			conCtx, err := consumer.Consume(func(msg jetstream.Msg) {
				msgData := msg.Data()
				msgSubject := msg.Subject()
				log := s.log.With().
					Str("subject", msgSubject).
					Str("stream", stream).
					Bytes("data", msgData).
					Logger()

				h, ok := subjects[msgSubject]
				if !ok {
					log.Warn().Msg("[nats] subject not found")
					return
				}
				retry, err := h(ctx, log, msgData)
				if err != nil {
					log.Error().Err(err).Msg("[nats] failed to process message")
				}
				for ctx.Err() == nil {
					if retry != 0 {
						err = msg.NakWithDelay(retry)
					} else if err != nil {
						err = msg.Nak()
					} else {
						err = msg.Ack()
					}
					if err == nil {
						break
					}
					timeu.SleepContext(ctx, time.Second)
				}
			})
			if err != nil {
				return err
			}
			s.mx.Lock()
			s.consumers = append(s.consumers, conCtx)
			s.mx.Unlock()
		}
	}
	return nil
}
