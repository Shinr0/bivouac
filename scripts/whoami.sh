#!/usr/bin/env bash
# Affiche et vérifie l'identité GCP utilisée par le projet.
set -euo pipefail
source "$(dirname "$0")/lib.sh"

require_cmd gcloud curl python3
require_gcloud_auth
log "OK — compte $GCP_ACCOUNT (CLI + ADC), projet $GCP_PROJECT_ID, config $CLOUDSDK_CONFIG"
