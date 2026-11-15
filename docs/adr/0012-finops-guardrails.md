# 0012. Garde-fous FinOps : budget GCP, OpenCost, coûts bornés par construction

- **Statut** : Accepté (budget en phase 1, OpenCost en phase 3)
- **Date** : 2026-11-15

## Contexte
Le coût du labo doit être borné par construction et visible par équipe.

## Décision
- **Budget GCP** sur le projet (socle), alertes à 50 / 90 / 100 % de la dépense réelle. Un budget
  alerte, il ne coupe rien : la vraie borne est la courte durée de vie de l'infrastructure.
- **OpenCost** : coût par namespace / label `team`.
- **Par construction** : nœuds spot, cluster zonal, pas de NAT par défaut, pas de GPU, ressources
  GCP d'équipe quasi gratuites, `check-orphans` en fin de `down` (ADR 0005, 0006, 0007).

## Conséquences
- La devise du budget doit correspondre à celle du compte de facturation.
- Les alertes vont aux administrateurs de facturation (destinataires IAM par défaut).
