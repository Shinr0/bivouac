#!/usr/bin/env bash
# PreToolUse(Bash) : bloque tout appel direct à gcloud qui n'utilise pas la config
# gcloud isolée du projet. Le compte gcloud global de la machine est celui d'un client.
# Passer par `make` (qui exporte CLOUDSDK_CONFIG) ou préfixer CLOUDSDK_CONFIG=...
# Code de sortie 2 = appel bloqué, le message (stderr) est renvoyé à Claude.
set -euo pipefail

cmd=$(python3 -c 'import json,sys; print(json.load(sys.stdin).get("tool_input",{}).get("command",""))')

# gcloud en position de commande : début de ligne, après ; & | ( ` $( ou un préfixe
# (sudo, env, exec, time, xargs, affectations VAR=...). Le mot « gcloud » dans un texte
# (message de commit, grep, echo) ne compte pas.
grep -Eq '(^|[;&|(`]|(sudo|env|exec|time|xargs)[[:space:]]+)[[:space:]]*([A-Za-z_][A-Za-z0-9_]*=[^[:space:]]*[[:space:]]+)*gcloud([[:space:]]|$)' <<<"$cmd" || exit 0
# Simple détection de présence de l'outil.
grep -Eq '(command -v|which|type) +gcloud' <<<"$cmd" && exit 0
# Config isolée explicitement fournie.
grep -Eq 'CLOUDSDK_CONFIG=[^[:space:]]*gcloud-platform-ai-demo' <<<"$cmd" && exit 0

echo "Bloqué par .claude/hooks/guard-gcloud.sh : gcloud appelé hors de la config isolée du projet \
(la config globale est celle d'un client). Utiliser une cible make (make whoami, make bootstrap…) \
ou préfixer CLOUDSDK_CONFIG=\$HOME/.config/gcloud-platform-ai-demo." >&2
exit 2
