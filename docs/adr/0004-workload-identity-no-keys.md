# 0004. Workload Identity partout, aucune clé de compte de service

- **Statut** : Accepté
- **Date** : 2026-11-10

## Contexte
Une clé JSON de compte de service est un secret longue durée, souvent la cause des fuites
d'identifiants cloud. Le dépôt est public.

## Décision
- Les pods s'authentifient à GCP via **Workload Identity Federation for GKE** (KSA ↔ GSA).
- Le seul secret applicatif (token GitHub) vit dans **Secret Manager** et est projeté dans le
  cluster par External Secrets Operator, lui-même authentifié par Workload Identity.
- L'opérateur utilise ses identifiants `gcloud` (ADC), jamais une clé exportée.
- La CI n'a pas besoin d'accès GCP en phase 1 ; si elle en a besoin, ce sera via Workload
  Identity Federation GitHub OIDC, pas une clé.

## Conséquences
- Aucune clé à faire tourner ni à révoquer.
- La valeur du token GitHub n'apparaît jamais dans le state Terraform.
