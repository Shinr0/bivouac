# 0016. Images sur GHCR, construites par GitHub Actions

- **Statut** : Accepté (implémentation en phase 6)
- **Date** : 2026-11-10

## Contexte
Les images (portail, job de clés LiteLLM) doivent survivre à `make down` et rester gratuites.

## Décision
GitHub Container Registry, images publiques, construites et étiquetées par SHA de commit par
GitHub Actions. Références épinglées dans `gitops/`.

## Alternatives écartées
- Artifact Registry : facturé et rattaché au projet, à exclure du périmètre détruit.

## Conséquences
- Aucun accès GCP requis pour la CI de build.
