# 0014. KServe + vLLM pour les modèles auto-hébergés

- **Statut** : Remplacé par [0017](0017-vertex-only-no-gpu.md) (2026-11-15)
- **Date** : 2026-11-10

## Contexte
Montrer qu'une équipe peut obtenir un modèle open-weights auto-hébergé derrière le même gateway,
sur GPU spot, sans payer de GPU au repos.

## Décision
KServe (control plane de serving CNCF) avec le runtime vLLM, sur un pool L4 spot à 0 nœud par
défaut, activé par `make gpu-on`. Déploiement uniquement si `enable_gpu`.

## Conséquences
- Mode de déploiement (avec/sans Knative), prérequis (cert-manager, Gateway API) et ressource
  (`InferenceService` ou `LLMInferenceService`) à trancher en phase 8, dans un ADR complémentaire.
- Prérequis manuel : quota GPU L4.
