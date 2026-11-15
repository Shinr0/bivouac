# 0018. Reporter l'IA intégrée à la plateforme

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Le plan initial plaçait deux volets IA *dans* la plateforme : un accès LLM par équipe (LiteLLM +
Vertex AI, RGD `LLMAccess`) et un assistant intégré au portail. Ils dépendent de tout le reste
(cluster, GitOps, API de plateforme, portail), ajoutent un prérequis manuel (Model Garden) et des coûts au token,
et déplacent l'attention de la démonstration principale : une IDP éphémère, gouvernée et frugale.

## Décision
- Ces volets deviennent les **phases optionnelles 7 et 8**, après le chemin local.
- Les ADR [0012](0012-litellm-gateway.md), [0013](0013-vertex-ai-claude.md) et
  [0017](0017-vertex-only-no-gpu.md) restent valables pour ces phases.
- L'IA intervient d'abord dans la **manière de construire** la plateforme (ADR 0019).

## Conséquences
- Plus de prérequis Model Garden ni de coût Vertex pour les phases 1 à 6.
- Les politiques Kyverno et le portail sont conçus sans dépendre de LiteLLM.
