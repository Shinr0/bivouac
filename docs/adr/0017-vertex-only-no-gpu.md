# 0017. Vertex AI uniquement : ni GPU ni modèle auto-hébergé

- **Statut** : Accepté — remplace [0014](0014-kserve-vllm.md)
- **Date** : 2026-11-15

## Contexte
L'ADR 0014 prévoyait un volet optionnel KServe + vLLM sur un pool GPU L4 spot. Ce volet ajoute :
un quota GPU à demander à la main (souvent nul sur un projet neuf, délai variable), un node pool
GPU, KServe et ses prérequis (cert-manager, éventuellement Knative), et une phase entière de
travail — pour une démonstration secondaire par rapport au message principal (self-service,
politiques, FinOps, assistant IA).

## Décision
- Les LLM sont consommés **uniquement via Vertex AI** (Claude, Model Garden), derrière LiteLLM.
- Aucun node pool GPU, aucun KServe/vLLM, pas de champ `backend` dans `LLMAccess`.
- Kyverno **interdit** toute demande de ressource `nvidia.com/gpu` dans les pods.
- La phase « Volet GPU » est supprimée ; le chemin local kind devient la phase 8.

## Conséquences
- Plus de quota GPU à demander ; `make gpu-on` / `make gpu-off` disparaissent.
- Coût d'infra d'une démo réduit au pool CPU spot ; le coût LLM est au token, plafonné par LiteLLM.
- GKE Standard zonal reste retenu (crédit gratuit, contrôle des nœuds spot — ADR 0005).
- LiteLLM permettrait d'ajouter un autre fournisseur plus tard sans changer l'API des équipes.
