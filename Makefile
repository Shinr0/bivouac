# Plateforme IA éphémère sur GCP — voir PLAN.md § 8.
SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c
.DEFAULT_GOAL := help

# Variables locales (non commitées) : copier .env.example en .env
-include .env
export

# Configuration gcloud isolée pour ce projet : n'utilise jamais le compte gcloud
# global de la machine (souvent celui d'un client). ADC incluses.
CLOUDSDK_CONFIG ?= $(HOME)/.config/gcloud-platform-ai-demo
GOOGLE_APPLICATION_CREDENTIALS ?= $(CLOUDSDK_CONFIG)/application_default_credentials.json
export CLOUDSDK_CONFIG GOOGLE_APPLICATION_CREDENTIALS

TF_DIRS := infra/bootstrap

.PHONY: help gcloud-login whoami bootstrap up demo down check-orphans local-up test \
        test-terraform test-hooks

help: ## Affiche cette aide
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

gcloud-login: ## Connexion du compte perso dans la config gcloud isolée (interactif)
	@scripts/gcloud-login.sh

whoami: ## Vérifie le compte et le projet GCP utilisés par le projet
	@scripts/whoami.sh

bootstrap: ## Socle persistant (une seule fois) : APIs, bucket de state, secret, budget
	@scripts/bootstrap.sh

up: ## Couche éphémère + attente Healthy (phase 2)
	@$(call todo,2)

demo: ## Port-forwards portail / Argo CD / Grafana (phase 3+)
	@$(call todo,3)

down: ## Destruction complète + contrôle d'orphelins (phase 2)
	@$(call todo,2)

check-orphans: ## Liste les ressources facturables résiduelles (phase 2)
	@$(call todo,2)

local-up: ## Chemin kind sans GCP (phase 8)
	@$(call todo,8)

test: test-terraform test-hooks ## Tous les contrôles statiques (s'étoffe à chaque phase)

test-hooks: ## Vérifie les garde-fous Claude Code (hook gcloud)
	@scripts/test-hooks.sh


test-terraform: ## terraform fmt + validate sur chaque couche
	terraform fmt -check -recursive infra
	@for d in $(TF_DIRS); do \
	  echo "==> validate $$d"; \
	  terraform -chdir=$$d init -backend=false -input=false >/dev/null; \
	  terraform -chdir=$$d validate; \
	done

define todo
echo "Cible non implémentée : prévue en phase $(1) (voir PLAN.md § 13)." >&2; exit 2
endef
