# 0013. Vertex AI (Claude via Model Garden) comme backend LLM

- **Statut** : Accepté — reporté en phase optionnelle 7 (ADR 0018)
- **Date** : 2026-11-10

## Contexte
Le backend LLM (unique, cf. ADR 0017) doit fonctionner sans GPU, être facturé à l'usage et s'authentifier sans clé.

## Décision
Claude servi par Vertex AI, appelé par LiteLLM via Workload Identity (`roles/aiplatform.user`).
Les modèles récents peuvent n'être disponibles que sur les endpoints `global`/`eu` : le cluster
reste en europe-west4, seul l'appel Vertex cible l'endpoint disponible.

## Conséquences
- Prérequis manuel : activer les modèles dans Model Garden (README).
- Région/endpoint et identifiants de modèles à vérifier en phase 5.
