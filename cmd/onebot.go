package cmd

import (
	"github.com/graydovee/xiaoshi/pkg/config"
	"github.com/graydovee/xiaoshi/pkg/core"
	"github.com/graydovee/xiaoshi/pkg/heathz"
	"log/slog"
)

func Run(cfg *config.Config) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// init botCore
	bot, err := core.NewOneBotCore(cfg)
	if err != nil {
		slog.Error("init bot error", "error", err)
		return
	}

	if err := bot.Start(); err != nil {
		slog.Error("start bot error", "error", err)
		return
	}

	if err := heathz.StartHealthz(&cfg.Healthz); err != nil {
		slog.Error("start healthz error", "error", err)
		return
	}

	select {}
}
