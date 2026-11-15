# 0011. Portail maison FastAPI + HTMX

- **Statut** : Accepté (implémentation en phases 6–7)
- **Date** : 2026-11-10

## Contexte
Il faut un portail pour le catalogue, les demandes et l'assistant IA, démontrable en quelques
minutes et maintenable par une personne.

## Décision
Application Python 3.12+ FastAPI, rendu serveur Jinja2 + HTMX, aucun framework JS lourd.
Le portail n'est qu'un **client** de l'API de la plateforme (les CRD générées par kro) et de GitHub : il lit le
cluster (RBAC lecture seule) et n'écrit que des pull requests.

## Alternatives écartées
- Backstage : riche mais lourd à installer, configurer et maintenir pour une démo éphémère.
- SPA React + BFF (approche de la v1) : deux applications et une chaîne de build pour peu de valeur ajoutée.

## Conséquences
- Remplaçable sans toucher à la plateforme : l'API, ce sont les CRD et Git.
