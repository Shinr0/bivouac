package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Shinr0/platform-project/bff/internal/vm"
)

const maxBodyBytes = 1 << 16

// catalog est servi par l'API pour que le frontend n'ait pas à dupliquer la
// liste des services ni les valeurs acceptées par la XRD.
var catalog = []vm.CatalogItem{
	{
		ID:          "virtual-machine",
		Title:       "Virtual Machine",
		Description: "Instance Scaleway provisionnée par Crossplane via la XRD VirtualMachine.",
		Available:   true,
	},
	{
		ID:          "postgres",
		Title:       "PostgreSQL managé",
		Description: "Base de données Scaleway RDB pour les équipes produit.",
		Available:   false,
	},
	{
		ID:          "object-bucket",
		Title:       "Bucket Object Storage",
		Description: "Stockage S3 compatible avec cycle de vie et versioning.",
		Available:   false,
	},
}

func (s *Server) handleCatalog(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"items":   catalog,
		"options": vm.Options(),
	})
}

func (s *Server) handleListVMs(w http.ResponseWriter, r *http.Request) {
	requests, err := s.vms.List(r.Context())
	if err != nil {
		slog.Error("list virtual machines", "error", err)
		writeError(w, http.StatusBadGateway, "impossible de lister les demandes")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": requests})
}

func (s *Server) handleCreateVM(w http.ResponseWriter, r *http.Request) {
	var params vm.Parameters

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&params); err != nil {
		writeError(w, http.StatusBadRequest, "corps de requête invalide: "+err.Error())
		return
	}

	created, err := s.vms.Create(r.Context(), params, requester(r))
	switch {
	case err == nil:
		writeJSON(w, http.StatusCreated, created)

	case errors.Is(err, vm.ErrAlreadyExists):
		writeError(w, http.StatusConflict, "une VM porte déjà ce nom")

	default:
		var invalid *vm.ValidationError
		if errors.As(err, &invalid) {
			writeFieldErrors(w, "paramètres invalides", invalid.Fields)
			return
		}
		slog.Error("create virtual machine", "error", err)
		writeError(w, http.StatusBadGateway, "impossible de créer la demande")
	}
}

func (s *Server) handleGetVM(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	found, err := s.vms.Get(r.Context(), name)
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, found)
	case errors.Is(err, vm.ErrNotFound):
		writeError(w, http.StatusNotFound, "demande introuvable")
	default:
		slog.Error("get virtual machine", "error", err, "name", name)
		writeError(w, http.StatusBadGateway, "impossible de lire la demande")
	}
}

func (s *Server) handleDeleteVM(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	err := s.vms.Delete(r.Context(), name)
	switch {
	case err == nil:
		writeJSON(w, http.StatusNoContent, nil)
	case errors.Is(err, vm.ErrNotFound):
		writeError(w, http.StatusNotFound, "demande introuvable")
	default:
		slog.Error("delete virtual machine", "error", err, "name", name)
		writeError(w, http.StatusBadGateway, "impossible de supprimer la demande")
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleReady vérifie que l'API Kubernetes répond et que la XRD est installée.
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if _, err := s.vms.List(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "API Kubernetes indisponible")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// requester identifie l'appelant. Derrière un proxy d'authentification, l'en-tête
// est renseigné par ce dernier ; sinon la demande est marquée anonyme.
func requester(r *http.Request) string {
	if user := strings.TrimSpace(r.Header.Get("X-Forge-User")); user != "" {
		return user
	}
	return "anonymous"
}
