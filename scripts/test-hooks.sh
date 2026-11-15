#!/usr/bin/env bash
# Vérifie que le hook .claude/hooks/guard-gcloud.sh bloque/autorise ce qu'il doit.
set -euo pipefail
source "$(dirname "$0")/lib.sh"
hook="$(dirname "$0")/../.claude/hooks/guard-gcloud.sh"
fail=0
check() { # attendu(0|2) commande
  local got=0
  python3 -c 'import json,sys;print(json.dumps({"tool_name":"Bash","tool_input":{"command":sys.argv[1]}}))' "$2" \
    | "$hook" >/dev/null 2>&1 || got=$?
  if [[ "$got" == "$1" ]]; then echo "  ok   ($got) $2"; else echo "  FAIL (attendu $1, obtenu $got) $2"; fail=1; fi
}
log "Hook guard-gcloud"
check 0 'ls -la'
check 0 'make whoami'
check 0 'command -v gcloud'
check 0 'CLOUDSDK_CONFIG=$HOME/.config/gcloud-platform-ai-demo gcloud config list'
check 2 'gcloud auth list'
check 2 'cd infra && gcloud projects list'
check 2 'CLOUDSDK_CONFIG=$HOME/.config/gcloud gcloud config list'
check 2 'echo "$(gcloud config get-value project)"'
check 2 'FOO=1 gcloud compute instances list'
check 2 'sudo gcloud auth login'
check 0 'git commit -m "block gcloud outside the isolated config"'
check 0 'grep -n gcloud scripts/lib.sh'
((fail == 0)) || die "Garde-fou gcloud défaillant"
