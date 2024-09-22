package session

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/wtkeqrf0/tg_form/pkg/util"
	"time"
)

const (
	day  = time.Hour * 24
	week = day * 7
)

type session struct {
	redis *redis.Client
}

//go:generate ifacemaker -f session.go -o session_if.go -i Session -s session -p session
func New(redis *redis.Client) Session {
	return &session{redis: redis}
}

func (s *session) Set(ctx context.Context, key string, payload *Payload) error {
	var (
		tx  = s.redis.TxPipeline()
		err = tx.HSet(ctx, key, payload).Err()
	)
	if err != nil {
		return err
	}

	if err = tx.Expire(ctx, key, week).Err(); err != nil {
		return err
	}

	_, err = tx.Exec(ctx)
	return err
}

func (s *session) SetValue(ctx context.Context, key string, fieldKey util.PayloadField, value any) error {
	var (
		tx  = s.redis.TxPipeline()
		err = tx.HSet(ctx, key, fieldKey.String(), value).Err()
	)
	if err != nil {
		return err
	}

	if err = tx.Expire(ctx, key, day).Err(); err != nil {
		return err
	}

	_, err = tx.Exec(ctx)
	return err
}

func (s *session) Get(ctx context.Context, key string) (*Payload, error) {
	payload := new(Payload)
	return payload, s.redis.HGetAll(ctx, key).Scan(payload)
}

func (s *session) Delete(ctx context.Context, keys ...string) error {
	return s.redis.Del(ctx, keys...).Err()
}

func (s *session) Keys(ctx context.Context, key string) ([]string, error) {
	return s.redis.HKeys(ctx, key).Result()
}
