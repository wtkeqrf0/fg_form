package main

import (
	"context"
	"github.com/wtkeqrf0/tg_form/internal/api"
	"github.com/wtkeqrf0/tg_form/internal/appeal"
	"github.com/wtkeqrf0/tg_form/internal/config"
	"github.com/wtkeqrf0/tg_form/internal/session"
	"github.com/wtkeqrf0/tg_form/pkg/db"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// ---------------- may fail ----------------
	var (
		cfg   = config.SetupConfig()
		psql  = db.NewPostgres(ctx, cfg.Postgres)
		redis = db.NewRedis(ctx, cfg.Redis)
		tgBot = api.New(ctx, cfg.TgToken)
	)

	// --------------- can't fail ---------------
	var (
		sessions = session.New(redis)

		appealR = appeal.NewRepo(psql)
		appealM = appeal.NewMethod(appealR, sessions)
	)

	log.Println("bot started successfully")
	tgBot.Start(appealM)
}
