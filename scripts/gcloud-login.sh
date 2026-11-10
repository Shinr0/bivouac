#!/usr/bin/env bash
# Connecte le compte GCP du projet dans une config gcloud isolée (CLOUDSDK_CONFIG),
# sans toucher à la config gcloud globale de la machine. Interactif (navigateur).
set -euo pipefail
source "$(dirname "$0")/lib.sh"

require_cmd gcloud
require_env GCP_ACCOUNT GCP_PROJECT_ID CLOUDSDK_CONFIG
[[ "$CLOUDSDK_CONFIG" != "$HOME/.config/gcloud" ]] || die "CLOUDSDK_CONFIG doit être isolé."

mkdir -p "$CLOUDSDK_CONFIG"
log "Config gcloud isolée : $CLOUDSDK_CONFIG"
gcloud auth login "$GCP_ACCOUNT" --update-adc
gcloud config set account "$GCP_ACCOUNT"
gcloud config set project "$GCP_PROJECT_ID"
gcloud config set disable_usage_reporting true
gcloud auth application-default set-quota-project "$GCP_PROJECT_ID"
"$(dirname "$0")/whoami.sh"
