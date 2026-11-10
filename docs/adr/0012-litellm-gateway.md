# 0012. LiteLLM comme gateway LLM

- **Statut** : Accepté (implémentation en phase 5)
- **Date** : 2026-11-10

## Contexte
Les équipes doivent consommer des LLM via une API unique, avec des clés individuelles, des
budgets et une mesure de consommation, quel que soit le backend (managé ou auto-hébergé).

## Décision
LiteLLM proxy : API compatible OpenAI, clés virtuelles, budgets par clé, routage multi-backends
(Vertex AI, KServe/vLLM). Base Postgres sans PVC (ADR 0007).

## Conséquences
- Syntaxe `model_list` Vertex, API admin des clés et authentification ADC à vérifier en phase 5.
- Le mécanisme de création de clé par claim `LLMAccess` fera l'objet d'un ADR dédié en phase 5.
