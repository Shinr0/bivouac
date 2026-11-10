// Command server expose l'API Forge qui crée des composites Crossplane
// VirtualMachine dans le cluster. Aucun YAML n'est généré ni commité :
// l'objet Kubernetes est la source de vérité, Crossplane le réconcilie.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shinr0/platform-project/bff/internal/config"
	"github.com/Shinr0/platform-project/bff/internal/httpapi"
	"github.com/Shinr0/platform-project/bff/internal/kube"
	"github.com/Shinr0/platform-project/bff/internal/vm"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}

	client, err := kube.NewDynamicClient(cfg.KubeconfigPath)
	if err != nil {
		return err
	}

	server := httpapi.NewServer(vm.NewService(client, cfg.Namespace))

	httpServer := &http.Server{
		Addr:         cfg.Addr,
		Handler:      server.Handler(cfg.AllowedOrigins),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", cfg.Addr, "namespace", cfg.Namespace)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}
