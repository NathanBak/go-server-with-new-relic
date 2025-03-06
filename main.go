package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NathanBak/cfgbuild"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/NathanBak/go-server-with-new-relic/internal/server"
	"github.com/NathanBak/go-server-with-new-relic/pkg/storage"
	"github.com/NathanBak/go-server-with-new-relic/pkg/widget"
)

func main() {
	// read .env file and set env vars
	_ = godotenv.Load()

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("go-server-with-new-relic"),
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		newrelic.ConfigCodeLevelMetricsPathPrefixes("go-agent/v3"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// create and initialize config from env vars
	builder := cfgbuild.Builder[*server.Config]{}
	cfg, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	cfg.Storage = &storage.MapStorage[widget.Widget]{}

	s, err := server.New(*cfg, app)
	if err != nil {
		log.Fatal(err)
	}

	// Start server running on separate thread
	go func() {
		err = s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// wait for signal and then shutdown cleanly
	quitchan := make(chan os.Signal, 1)
	signal.Notify(quitchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quitchan
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
