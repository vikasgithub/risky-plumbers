package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/vikasgithub/risky-plumbers/internal/config"
	"github.com/vikasgithub/risky-plumbers/internal/healthcheck"
	"github.com/vikasgithub/risky-plumbers/internal/log"
	"github.com/vikasgithub/risky-plumbers/internal/risk"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()
	logger := log.New()

	// load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	healthcheck.RegisterHandlers(r)
	apiRouter := buildApiRouter(logger)
	r.Mount("/api/v1", apiRouter)

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.Server.Port)
	server := &http.Server{
		Addr:    address,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err)
			os.Exit(-1)
		}
	}()
	logger.Infof("server is running at %v", address)

	<-ctx.Done()
	logger.Info("got interruption signal")
	if err := server.Shutdown(context.TODO()); err != nil {
		logger.Errorf("server shutdown returned an err: %v\n", err)
	}

	logger.Info("Server stopped...")
}

func buildApiRouter(logger log.Logger) *chi.Mux {
	r := chi.NewRouter()

	//Add handlers here
	risk.RegisterHandlers(r, risk.NewService(risk.NewRepository(logger), logger))

	return r
}
