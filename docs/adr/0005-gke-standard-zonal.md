# 0005. GKE Standard zonal, nœuds spot

- **Statut** : Accepté (implémentation en phase 2)
- **Date** : 2026-11-10

## Contexte
Il faut un cluster Kubernetes managé, créé et détruit en quelques minutes, au coût le plus bas possible.

## Décision
- **GKE Standard, cluster zonal** (europe-west4) : le crédit gratuit GKE couvre les frais de
  gestion d'un cluster zonal ; Standard permet de piloter finement les node pools (spot, type de machine, autoscaling).
- Nœuds **spot** sur tous les pools ; aucun GPU ([ADR 0017](0017-vertex-only-no-gpu.md)).
- `deletion_protection = false` : le cluster doit pouvoir être détruit sans intervention.

## Alternatives écartées
- GKE Autopilot : facturation par pod pratique, mais moins de contrôle sur les nœuds spot et leur
  dimensionnement, et un comportement de coût moins lisible en démo.
- Cluster régional : frais de gestion non couverts par le crédit gratuit, triple de nœuds.

## Conséquences
- Pas de haute disponibilité du control plane : acceptable pour une démo.
- Les nœuds spot peuvent être préemptés : les workloads de démo doivent le tolérer.
- Montants à revérifier dans le calculateur GCP (PLAN.md § 11).
