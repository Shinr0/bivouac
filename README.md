# Plateforme IA éphémère sur GCP

Une **plateforme interne de développement (IDP) en self-service, avec de l'IA par-dessus**, qui
n'existe que le temps d'une démo : `make up` la construit depuis un projet GCP vide, `make down`
la détruit sans laisser de ressource facturée.

- **LLM-as-a-Service** : une équipe demande un accès LLM par un simple claim Git ; la plateforme
  provisionne la clé, applique quotas et budgets, mesure les coûts.
- **Assistant de plateforme** : un agent en langage naturel aide à formuler une demande et
  **ouvre une pull request** — il ne touche jamais le cluster directement.

> Projet en construction, phase par phase. Spécification complète : [`PLAN.md`](PLAN.md).

## Architecture

```
GitHub (gitops/, claims/) ──sync──► Argo CD ──► Kyverno ──► kro ──► ressources par équipe
                                        │
GKE Standard zonal (spot)               ├─ Portail FastAPI + HTMX (catalogue, demandes, assistant)
                                        ├─ LiteLLM ──► Vertex AI (Claude)
                                        └─ OpenCost, Prometheus/Grafana
```

Détails : [`docs/architecture.md`](docs/architecture.md) · Décisions : [`docs/adr/`](docs/adr/README.md)

| Couche | Outil |
|---|---|
| IaC | Terraform (compatible OpenTofu), deux couches : socle persistant / plateforme éphémère |
| Cluster | GKE Standard zonal, nœuds spot, aucun GPU |
| GitOps | Argo CD, app-of-apps |
| API développeur | kro (ResourceGraphDefinitions) |
| Politiques | Kyverno (CI + admission) |
| Portail | FastAPI + HTMX |
| Gateway LLM | LiteLLM → Vertex AI (Claude), seul backend — pas de modèle auto-hébergé |
| FinOps | Budget GCP, OpenCost, budgets LiteLLM |

## Avancement

| Phase | Contenu | État |
|---|---|---|
| 1 | Squelette, ADR, Makefile, `infra/bootstrap` | ✅ |
| 2 | `infra/platform` (VPC, GKE, IAM, WI) + `up` / `down` / `check-orphans` | ⏳ |
| 3 | Argo CD app-of-apps, Kyverno, ESO, kro, OpenCost | — |
| 4 | RGD `TeamSpace`, politiques, CI Kyverno | — |
| 5 | LiteLLM + Vertex, RGD `LLMAccess` | — |
| 6 | Portail : catalogue, formulaires, PR | — |
| 7 | Assistant IA | — |
| 8 | Chemin local kind, polish | — |

## Prérequis

### Outils

| Outil | Version testée |
|---|---|
| Terraform (ou OpenTofu) | 1.17 (≥ 1.9 requis) |
| Google Cloud SDK (`gcloud`) | 590 |
| GNU Make, Bash | 4.x / 5.x |

`kubectl`, `helm`, la CLI Kyverno et Python seront requis à partir des phases suivantes.

### Actions manuelles (une seule fois)

Ce sont les **seules** opérations hors code (PLAN.md § 10).

1. **Projet GCP avec facturation et identifiants isolés.** Le projet utilise sa propre
   configuration gcloud (`~/.config/gcloud-platform-ai-demo`, ADC comprises) : la configuration
   gcloud globale de la machine n'est jamais lue ni modifiée. Renseigner `GCP_ACCOUNT` et
   `GCP_PROJECT_ID` dans `.env` (étape 3), puis :
   ```bash
   make gcloud-login   # ouvre le navigateur : connexion CLI + ADC sur GCP_ACCOUNT
   make whoami         # vérifie compte (CLI et ADC) et projet
   ```
   Tous les scripts refusent de s'exécuter si le compte actif n'est pas `GCP_ACCOUNT`.
   Le compte doit pouvoir activer des APIs sur le projet et créer un budget sur le compte de
   facturation (`roles/billing.costsManager` ou administrateur de facturation).
2. **Activer les modèles Claude dans Vertex AI Model Garden** et accepter les conditions
   (console GCP › Vertex AI › Model Garden). Non automatisable proprement, par projet.
3. **Configuration locale** :
   ```bash
   cp .env.example .env   # GCP_ACCOUNT, GCP_PROJECT_ID, GCP_BILLING_ACCOUNT_ID
   ```
4. **Socle persistant** :
   ```bash
   make bootstrap
   ```
   Crée : activation des APIs, bucket de state `<PROJECT_ID>-tfstate` (versionné), secret vide
   `github-token`, budget mensuel (50 € par défaut, alertes à 50 / 90 / 100 %). Le state de cette
   couche reste local — voir [ADR 0003](docs/adr/0003-bootstrap-local-state.md).
5. **Token GitHub** : créer un token *fine-grained* limité à ce dépôt (permissions
   *Contents* et *Pull requests* en lecture/écriture), puis l'ajouter au secret sans qu'il
   passe par l'historique du shell :
   ```bash
   export CLOUDSDK_CONFIG=~/.config/gcloud-platform-ai-demo   # config isolée
   read -rs TOKEN && printf '%s' "$TOKEN" | \
     gcloud secrets versions add github-token --data-file=- && unset TOKEN
   ```

## Utilisation

```bash
make help        # liste des cibles
make whoami      # compte et projet GCP utilisés
make test        # contrôles statiques (terraform fmt/validate, puis lint, kyverno, pytest)
make bootstrap   # socle persistant, idempotent
make up          # plateforme éphémère          (phase 2)
make down        # destruction + check-orphans  (phase 2)
```

## Coûts

| Élément | Ordre de grandeur |
|---|---|
| Socle (bucket, secret, budget) | ≈ 0 € / mois |
| Nœud e2-standard-4 spot (europe-west4) | ≈ 0,056 $/h |
| Frais de gestion GKE (cluster zonal) | 0,10 $/h, couverts par le crédit gratuit GKE |
| Vertex AI (Claude) | au token, plafonné par les budgets LiteLLM |
| Démo de 2 h (infra) | < 1 $ |

Prix relevés en novembre 2026, à revérifier dans le calculateur GCP. Garde-fous : budget GCP,
nœuds spot, aucun GPU, budgets LiteLLM, interdiction des PVC/LoadBalancer, contrôle d'orphelins
après chaque `down` ([ADR 0015](docs/adr/0015-finops-guardrails.md)).

## Structure du dépôt

```
├── PLAN.md               spécification et phases
├── CLAUDE.md             conventions pour Claude Code
├── Makefile
├── docs/                 architecture, ADR
├── infra/bootstrap/      socle persistant (Terraform, state local)
├── scripts/              scripts appelés par le Makefile
└── .github/workflows/    CI
```

La structure cible complète est décrite dans PLAN.md § 14. L'ancienne version de ce dépôt
(Crossplane + portail React) est conservée dans la branche `archive/v1-crossplane-forge`.
