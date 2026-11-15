# 0011. Portail maison FastAPI + HTMX

- **Statut** : Accepté (implémentation en phase 5)
- **Date** : 2026-11-15

## Contexte
Il faut un portail pour le catalogue et les demandes, démontrable en quelques minutes et
maintenable par une personne.

## Décision
Application Python 3.12+ FastAPI, rendu serveur Jinja2 + HTMX, aucun framework JS lourd.
Le portail n'est qu'un **client** de l'API de la plateforme (les CRD générées par kro) et de
GitHub : il lit le cluster (RBAC lecture seule) et n'écrit que des pull requests. Ses formulaires
sont générés depuis le schéma OpenAPI des CRD.

## Alternatives écartées
- Backstage : riche mais lourd à installer, configurer et maintenir pour un labo éphémère.
- SPA + BFF : deux applications et une chaîne de build pour peu de valeur ajoutée.

## Conséquences
- Remplaçable sans toucher à la plateforme : l'API, ce sont les CRD et Git.
