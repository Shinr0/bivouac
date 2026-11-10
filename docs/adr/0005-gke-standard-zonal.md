# 0005. GKE Standard zonal, nœuds spot

- **Statut** : Accepté (implémentation en phase 2)
- **Date** : 2026-11-10

## Contexte
Il faut un cluster Kubernetes managé, créé et détruit en quelques minutes, avec un pool GPU
optionnel, au coût le plus bas possible.

## Décision
- **GKE Standard, cluster zonal** (europe-west4) : le crédit gratuit GKE couvre les frais de
  gestion d'un cluster zonal ; Standard permet de piloter finement les node pools (GPU, spot, taints).
- Nœuds **spot** sur tous les pools ; pool GPU L4 à **0** nœud par défaut.
- `deletion_protection = false` : le cluster doit pouvoir être détruit sans intervention.

## Alternatives écartées
- GKE Autopilot : facturation par pod pratique, mais moins de contrôle sur les GPU spot et les
  taints, et un comportement de coût moins lisible en démo.
- Cluster régional : frais de gestion non couverts par le crédit gratuit, triple de nœuds.

## Conséquences
- Pas de haute disponibilité du control plane : acceptable pour une démo.
- Les nœuds spot peuvent être préemptés : les workloads de démo doivent le tolérer.
- Montants à revérifier dans le calculateur GCP (PLAN.md § 11).
