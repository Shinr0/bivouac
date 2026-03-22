#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${CLUSTER_NAME:-platform-eng}"
ARGOCD_NAMESPACE="${ARGOCD_NAMESPACE:-argocd}"
CLAIM_NAMESPACE="${CLAIM_NAMESPACE:-default}"
WAIT_TIMEOUT_SECONDS="${WAIT_TIMEOUT_SECONDS:-600}"
WAIT_INTERVAL_SECONDS="${WAIT_INTERVAL_SECONDS:-10}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { echo -e "${GREEN}[+]${NC} $*"; }
warn() { echo -e "${YELLOW}[!]${NC} $*"; }
err()  { echo -e "${RED}[x]${NC} $*"; exit 1; }

command_exists() { command -v "$1" &>/dev/null; }

check_prerequisites() {
  local missing=()
  for cmd in kind kubectl; do
    command_exists "$cmd" || missing+=("$cmd")
  done

  if (( ${#missing[@]} > 0 )); then
    err "Missing tools: ${missing[*]}. Install them first."
  fi
}

cluster_exists() {
  kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"
}

delete_argocd_application() {
  if kubectl get application crossplane-claims -n "$ARGOCD_NAMESPACE" &>/dev/null; then
    log "Deleting ArgoCD Application 'crossplane-claims'..."
    kubectl delete application crossplane-claims -n "$ARGOCD_NAMESPACE" --wait=true
  else
    warn "ArgoCD Application 'crossplane-claims' not found, skipping."
  fi
}

delete_virtualmachine_claims() {
  if kubectl api-resources --namespaced=true -o name | grep -q '^virtualmachines.compute.platform.local$'; then
    log "Deleting VirtualMachine claims in namespace '${CLAIM_NAMESPACE}'..."
    kubectl delete virtualmachines.compute.platform.local \
      --all -n "$CLAIM_NAMESPACE" --ignore-not-found=true
  else
    warn "CRD 'virtualmachines.compute.platform.local' not found, skipping claims deletion."
  fi
}

count_remaining_gcp_resources() {
  local total=0
  local resource
  local resources=(
    "instances.compute.gcp.upbound.io"
    "addresses.compute.gcp.upbound.io"
    "disks.compute.gcp.upbound.io"
    "attacheddisks.compute.gcp.upbound.io"
  )

  for resource in "${resources[@]}"; do
    if kubectl api-resources -o name | grep -q "^${resource}$"; then
      local count
      count=$(kubectl get "$resource" -A --no-headers 2>/dev/null | wc -l | tr -d ' ')
      total=$((total + count))
    fi
  done

  echo "$total"
}

wait_for_crossplane_cleanup() {
  local waited=0
  local remaining

  log "Waiting for Crossplane managed resources to be deleted..."
  while true; do
    remaining="$(count_remaining_gcp_resources)"
    if [[ "$remaining" == "0" ]]; then
      log "No remaining GCP managed resources detected."
      break
    fi

    if (( waited >= WAIT_TIMEOUT_SECONDS )); then
      warn "Timeout reached (${WAIT_TIMEOUT_SECONDS}s) with ${remaining} resources still present."
      warn "You can inspect with: kubectl get managed -A"
      break
    fi

    warn "Still ${remaining} managed resources present... waiting ${WAIT_INTERVAL_SECONDS}s"
    sleep "$WAIT_INTERVAL_SECONDS"
    waited=$((waited + WAIT_INTERVAL_SECONDS))
  done
}

delete_kind_cluster() {
  if cluster_exists; then
    log "Deleting Kind cluster '${CLUSTER_NAME}'..."
    kind delete cluster --name "$CLUSTER_NAME"
    log "Cluster deleted."
  else
    warn "Kind cluster '${CLUSTER_NAME}' not found, nothing to delete."
  fi
}

main() {
  echo "========================================="
  echo " Platform Engineering Teardown"
  echo "========================================="

  check_prerequisites

  if ! cluster_exists; then
    warn "Cluster '${CLUSTER_NAME}' does not exist. Nothing to stop."
    exit 0
  fi

  # Keep the right kube-context for cluster operations.
  kubectl cluster-info --context "kind-${CLUSTER_NAME}" &>/dev/null \
    || err "Cluster exists but context 'kind-${CLUSTER_NAME}' is unreachable."

  delete_argocd_application
  delete_virtualmachine_claims
  wait_for_crossplane_cleanup
  delete_kind_cluster

  echo ""
  log "Teardown complete."
}

main "$@"
