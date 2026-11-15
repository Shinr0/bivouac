# 0013. Images sur GHCR, construites par GitHub Actions

- **Statut** : Accepté (implémentation en phase 5)
- **Date** : 2026-11-15

## Contexte
Les images du labo (portail) doivent survivre à `make down` et rester gratuites.

## Décision
GitHub Container Registry, images publiques, construites et étiquetées par SHA de commit par
GitHub Actions. Références épinglées dans `gitops/`.

## Alternatives écartées
- Artifact Registry : facturé et rattaché au projet, à exclure du périmètre détruit.

## Conséquences
- Aucun accès GCP requis pour la CI de build.
