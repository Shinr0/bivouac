# 0004. Workload Identity et OIDC partout, aucune clé stockée

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Une clé JSON de compte de service ou une clé d'API est un secret longue durée, cause classique
de fuite d'identifiants. Le dépôt est public.

## Décision
- Les pods s'authentifient à GCP via **Workload Identity Federation for GKE** (compte Kubernetes ↔
  compte de service Google).
- Le seul secret applicatif (token GitHub) vit dans **Secret Manager** et est projeté dans le
  cluster par External Secrets Operator, lui-même authentifié par Workload Identity.
- L'opérateur utilise ses identifiants `gcloud` (ADC) dans une **configuration gcloud isolée**
  propre au projet, jamais une clé exportée ; les scripts vérifient le compte effectif.
- La CI s'authentifie par **OIDC** (fédération d'identité) quand elle a besoin d'un service
  externe : aucun secret longue durée dans GitHub.

## Conséquences
- Aucune clé à faire tourner ni à révoquer.
- La valeur du token GitHub n'apparaît jamais dans le state Terraform.
