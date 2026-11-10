// Package config charge la configuration du BFF depuis l'environnement.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	// Addr est l'adresse d'écoute HTTP.
	Addr string
	// Namespace héberge les composites VirtualMachine créés par le portail.
	Namespace string
	// AllowedOrigins active CORS pour le frontend servi depuis une autre origine.
	AllowedOrigins []string
	// KubeconfigPath n'est utilisé qu'en dehors du cluster.
	KubeconfigPath string

	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func FromEnv() (Config, error) {
	cfg := Config{
		Addr:            envOr("FORGE_ADDR", ":8080"),
		Namespace:       envOr("FORGE_NAMESPACE", "forge-workloads"),
		KubeconfigPath:  os.Getenv("KUBECONFIG"),
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    30 * time.Second,
		ShutdownTimeout: 15 * time.Second,
	}

	if raw := strings.TrimSpace(os.Getenv("FORGE_CORS_ORIGINS")); raw != "" {
		for _, origin := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(origin); trimmed != "" {
				cfg.AllowedOrigins = append(cfg.AllowedOrigins, trimmed)
			}
		}
	}

	if cfg.Namespace == "" {
		return Config{}, fmt.Errorf("FORGE_NAMESPACE must not be empty")
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
