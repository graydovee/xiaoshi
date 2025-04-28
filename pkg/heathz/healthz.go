package heathz

import (
	"fmt"
	"github.com/graydovee/xiaoshi/pkg/config"
	"log/slog"
	"net"
	"net/http"
)

func StartHealthz(cfg *config.HealthzConfig) error {
	if cfg == nil || !cfg.Enabled {
		slog.Info("healthz is disabled")
		return nil
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port))
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Pattern, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	slog.Info("starting healthz http server on " + cfg.Addr)
	if err := http.Serve(listener, mux); err != nil {
		return err
	}
	return nil
}
