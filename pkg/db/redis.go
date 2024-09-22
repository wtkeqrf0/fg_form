package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wtkeqrf0/tg_form/internal/config"
	"log"
)

func NewRedis(ctx context.Context, cfg config.Redis) *redis.Client {
	cl := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.Db,
	})

	if err := cl.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	context.AfterFunc(ctx, func() {
		if err := cl.Close(); err != nil {
			log.Println(err.Error())
		}
	})
	return cl
}
