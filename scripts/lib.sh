#!/usr/bin/env bash
# Fonctions communes aux scripts de la plateforme.
set -euo pipefail

log()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33m[!]\033[0m %s\n' "$*" >&2; }
die()  { printf '\033[1;31m[x]\033[0m %s\n' "$*" >&2; exit 1; }

require_cmd() {
  local missing=()
  for c in "$@"; do command -v "$c" >/dev/null 2>&1 || missing+=("$c"); done
  ((${#missing[@]} == 0)) || die "Outils manquants : ${missing[*]}"
}

require_env() {
  for v in "$@"; do [[ -n "${!v:-}" ]] || die "Variable $v non définie (voir .env.example)"; done
}

# Refuse de continuer si gcloud ou les ADC ne sont pas sur le compte attendu (GCP_ACCOUNT)
# et sur la config isolée du projet. Protège contre l'usage du compte d'un client.
require_gcloud_auth() {
  require_env GCP_ACCOUNT GCP_PROJECT_ID CLOUDSDK_CONFIG GOOGLE_APPLICATION_CREDENTIALS
  [[ "$CLOUDSDK_CONFIG" != "$HOME/.config/gcloud" ]] \
    || die "CLOUDSDK_CONFIG pointe sur la config gcloud globale ; utiliser la config isolée du projet."

  local cli_account
  cli_account=$(gcloud config get-value account 2>/dev/null || true)
  [[ "$cli_account" == "$GCP_ACCOUNT" ]] \
    || die "Compte gcloud actif '${cli_account:-aucun}' ≠ GCP_ACCOUNT. Lancer 'make gcloud-login'."

  [[ -f "$GOOGLE_APPLICATION_CREDENTIALS" ]] \
    || die "ADC absentes ($GOOGLE_APPLICATION_CREDENTIALS). Lancer 'make gcloud-login'."
  local token adc_account
  token=$(gcloud auth application-default print-access-token 2>/dev/null) \
    || die "ADC invalides. Lancer 'make gcloud-login'."
  adc_account=$(curl -fsS "https://oauth2.googleapis.com/tokeninfo?access_token=${token}" \
    | python3 -c 'import json,sys; print(json.load(sys.stdin).get("email",""))') || true
  [[ "$adc_account" == "$GCP_ACCOUNT" ]] \
    || die "ADC sur le compte '${adc_account:-inconnu}' ≠ GCP_ACCOUNT. Lancer 'make gcloud-login'."

  local project
  project=$(gcloud config get-value project 2>/dev/null || true)
  [[ "$project" == "$GCP_PROJECT_ID" ]] \
    || die "Projet gcloud '${project:-aucun}' ≠ GCP_PROJECT_ID. Lancer 'make gcloud-login'."
}
