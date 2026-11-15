# 0005. GKE Standard zonal, nœuds spot, sans GPU

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Il faut un cluster Kubernetes managé, créé et détruit en quelques minutes, au coût le plus bas
possible, pour des essais de platform engineering qui n'exigent pas de calcul spécialisé.

## Décision
- **GKE Standard, cluster zonal** (europe-west4) : le crédit gratuit GKE couvre les frais de
  gestion d'un cluster zonal ; Standard permet de piloter les node pools (spot, type de machine,
  autoscaling).
- Nœuds **spot** sur tous les pools.
- **Aucun GPU** : pas de pool GPU, et une politique Kyverno refuse toute demande `nvidia.com/gpu`.
- `deletion_protection = false` : le cluster doit pouvoir être détruit sans intervention.

## Alternatives écartées
- GKE Autopilot : facturation par pod pratique, mais moins de contrôle sur les nœuds spot et leur
  dimensionnement, et un coût moins lisible.
- Cluster régional : frais de gestion non couverts par le crédit gratuit, nœuds triplés.
- Pool GPU pour héberger des modèles : quota à demander, coût horaire élevé, hors du sujet du labo.

## Conséquences
- Pas de haute disponibilité du control plane : acceptable pour un labo.
- Les nœuds spot peuvent être préemptés : les workloads doivent le tolérer.
- Montants à revérifier dans le calculateur GCP.
