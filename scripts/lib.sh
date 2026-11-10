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

require_gcloud_auth() {
  gcloud auth application-default print-access-token >/dev/null 2>&1 \
    || die "ADC absentes : lancer 'gcloud auth application-default login'"
}
