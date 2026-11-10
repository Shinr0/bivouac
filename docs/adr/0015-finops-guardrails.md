# 0015. Garde-fous FinOps : budget GCP, OpenCost, budgets LiteLLM

- **Statut** : Accepté (budget GCP en phase 1, le reste en phases 3–5)
- **Date** : 2026-11-10

## Contexte
Le coût doit être borné par construction et visible par équipe.

## Décision
- **Budget GCP** sur le projet (socle, phase 1), alertes à 50 / 90 / 100 % de la dépense réelle.
  Un budget alerte, il ne coupe rien : la vraie borne est la courte durée de vie de l'infra.
- **OpenCost** : coût par namespace / label `team`.
- **Budgets LiteLLM** : plafond par clé d'équipe, et plafond global.
- **Kyverno** et **`check-orphans`** : ADR 0007.

## Conséquences
- La devise du budget doit correspondre à celle du compte de facturation.
- Les alertes vont aux administrateurs de facturation (destinataires IAM par défaut).
