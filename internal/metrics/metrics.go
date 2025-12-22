package metrics

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	Registry = prometheus.NewRegistry()
)

func NewHTTPHandler() http.Handler {
	return promhttp.HandlerFor(Registry, promhttp.HandlerOpts{})
}

func NewPusher(ctx context.Context, pushgw string) error {
	pusher := push.New(pushgw, "ged_trade").Gatherer(Registry)

	// create a ticker that will push the metrics every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := pusher.Push(); err != nil {
					slog.Error("error pushing", "error", err)
				}
			}
		}
	}()
	return nil
}
