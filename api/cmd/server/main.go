// Command server runs the Bongkar Ilusi API and serves the built SPA.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/januarsyah901/ascard/api/internal/app"
	"github.com/januarsyah901/ascard/api/internal/auth"
	"github.com/januarsyah901/ascard/api/internal/config"
	"github.com/januarsyah901/ascard/api/internal/db"
	"github.com/januarsyah901/ascard/api/internal/httpx"
	"github.com/januarsyah901/ascard/api/internal/server"
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "GET /api/health on 127.0.0.1:$PORT and exit 0 if it answers 200 (for Docker HEALTHCHECK)")
	flag.Parse()
	if *healthcheck {
		os.Exit(checkHealth())
	}
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	setupLogger(cfg.LogFormat)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := db.EnsureSchema(ctx, pool); err != nil {
		return err
	}

	handler := server.New(app.Deps{
		Config:   cfg,
		Pool:     pool,
		JWT:      auth.NewJWT(cfg.JWTSecret),
		ClientIP: httpx.ClientIP{Trusted: cfg.TrustedProxies},
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server running", "port", cfg.Port, "env", cfg.AppEnv, "web_dir", cfg.WebDir)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
	}
	return nil
}

// checkHealth is used by the container HEALTHCHECK (distroless has no shell or wget).
func checkHealth() int {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}
	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/api/health")
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "healthcheck: status", resp.StatusCode)
		return 1
	}
	return 0
}

func setupLogger(format string) {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var h slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if format == "json" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(h))
}
