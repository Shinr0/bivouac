#!/usr/bin/env bash
# Crée/met à jour le socle persistant (infra/bootstrap). Idempotent.
set -euo pipefail
source "$(dirname "$0")/lib.sh"

require_cmd terraform gcloud
require_env GCP_PROJECT_ID GCP_BILLING_ACCOUNT_ID
require_gcloud_auth

start=$(date +%s)
export TF_VAR_project_id="$GCP_PROJECT_ID"
export TF_VAR_billing_account_id="$GCP_BILLING_ACCOUNT_ID"
export TF_VAR_region="${GCP_REGION:-europe-west4}"
export TF_VAR_budget_amount="${BUDGET_AMOUNT:-50}"
export TF_VAR_budget_currency="${BUDGET_CURRENCY:-EUR}"

dir="$(dirname "$0")/../infra/bootstrap"
log "terraform init (state local)"
terraform -chdir="$dir" init -input=false
log "terraform apply"
terraform -chdir="$dir" apply -input=false

secret=$(terraform -chdir="$dir" output -raw github_token_secret)
if ! gcloud secrets versions list github-token --project "$GCP_PROJECT_ID" \
     --filter=state=enabled --format='value(name)' 2>/dev/null | grep -q .; then
  warn "Le secret $secret n'a aucune version active."
  warn "Ajouter le token : printf '%s' \"\$TOKEN\" | gcloud secrets versions add github-token --data-file=- --project $GCP_PROJECT_ID"
fi

log "Socle prêt en $(( $(date +%s) - start ))s. Bucket de state : $(terraform -chdir="$dir" output -raw tfstate_bucket)"
