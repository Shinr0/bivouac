// Package httpapi expose l'API REST consommée par le portail Forge.
package httpapi

import (
	"net/http"

	"github.com/Shinr0/platform-project/bff/internal/vm"
)

type Server struct {
	vms *vm.Service
}

func NewServer(vms *vm.Service) *Server {
	return &Server{vms: vms}
}

// Handler assemble les routes et les middlewares.
func (s *Server) Handler(allowedOrigins []string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("GET /readyz", s.handleReady)

	mux.HandleFunc("GET /api/v1/catalog", s.handleCatalog)
	mux.HandleFunc("GET /api/v1/vms", s.handleListVMs)
	mux.HandleFunc("POST /api/v1/vms", s.handleCreateVM)
	mux.HandleFunc("GET /api/v1/vms/{name}", s.handleGetVM)
	mux.HandleFunc("DELETE /api/v1/vms/{name}", s.handleDeleteVM)

	return chain(mux, withRecovery, withLogging, withCORS(allowedOrigins))
}
