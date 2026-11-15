# 0008. Argo CD, pattern app-of-apps avec sync waves

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Tout ce qui tourne dans le cluster doit être déclaré dans Git. L'ordre d'installation compte
(CRD avant ressources, Kyverno avant les demandes des équipes).

## Décision
- Argo CD est le seul composant installé par Terraform (provider Helm), avec l'Application racine.
- Le reste est décrit sous `gitops/` : une `Application` par composant, ordonnées par
  `argocd.argoproj.io/sync-wave`. Les composants tiers sont installés par leurs charts Helm.
- Overlays `gcp` et `local` pour les différences d'environnement.

## Alternatives écartées
- Flux : équivalent fonctionnel, mais moins visuel en démonstration.
- Tout installer par Terraform (Helm) : couplage fort, pas de réconciliation continue.

## Conséquences
- `down.sh` supprime l'Application racine en cascade et gère les finalizers.
