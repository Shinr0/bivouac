# 0015. IA intégrée à la plateforme : optionnelle, LiteLLM + Vertex AI

- **Statut** : Accepté (phases optionnelles 7–8)
- **Date** : 2026-11-15

## Contexte
Deux usages de l'IA *dans* la plateforme sont envisagés : un accès LLM par équipe avec budget, et
un assistant intégré au portail qui aide à formuler une demande. Ils dépendent de tout le reste
(cluster, GitOps, API de plateforme, portail) et ajoutent un prérequis manuel et des coûts au token.

## Décision
- Ces usages sont des **phases optionnelles (7 et 8)**, après le cœur du labo.
- Gateway **LiteLLM** (API compatible OpenAI, clés virtuelles, budgets par clé) devant **Vertex AI**
  (Claude via Model Garden), authentifié par Workload Identity.
- **Aucun modèle auto-hébergé** : pas de GPU (ADR 0005).
- L'assistant n'a qu'un droit d'écriture : ouvrir une pull request, après confirmation explicite.

## Alternatives écartées
- KServe + vLLM sur GPU spot : quota à demander, coût horaire, hors du sujet du labo.
- Appels directs à Vertex par équipe : pas de budget ni de clé par équipe.

## Conséquences
- Prérequis manuel, pour ces phases seulement : activer les modèles Claude dans Model Garden.
- Région / endpoint Vertex (`global`, `eu`) et syntaxe LiteLLM à vérifier à l'implémentation.
