package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wtkeqrf0/tg_form/internal/config"
)

func NewPostgres(ctx context.Context, p config.Postgres) *pgxpool.Pool {
	cfgStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=10",
		p.Host, p.Port, p.User, p.Password, p.DbName, p.SslMode)

	cfg, err := pgxpool.ParseConfig(cfgStr)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}

	if err = pool.Ping(ctx); err != nil {
		panic(err)
	}

	context.AfterFunc(ctx, pool.Close)
	return pool
}
